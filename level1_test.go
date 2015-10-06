// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

import (
	"math"
	"reflect"
	"testing"

	"github.com/gonum/matrix/mat64"
)

func TestDot(t *testing.T) {
	for _, test := range []struct {
		n       int
		x, y    []float64
		indices []int

		want float64
	}{
		{
			n:       5,
			x:       []float64{1, 2, 3},
			indices: []int{0, 2, 4},
			y:       []float64{1, math.NaN(), 3, math.NaN(), 5},

			want: 22,
		},
	} {
		x := NewVector(test.n, test.x, test.indices)
		y := mat64.NewVector(len(test.y), test.y)
		got := Dot(x, y)
		if got != test.want {
			t.Errorf("want = %v, got %v\n", test.want, got)
		}
	}
}

func TestAxpy(t *testing.T) {
	for _, test := range []struct {
		a     float64
		x, y  []float64
		index []int
		incy  int

		want []float64
	}{
		{
			a:     1,
			x:     []float64{1, 2, 3},
			index: []int{0, 2, 3},
			y:     []float64{0, 0, 0, 0},
			incy:  1,

			want: []float64{1, 0, 2, 3},
		},
		{
			a:     2,
			x:     []float64{1, 2, 3},
			index: []int{0, 2, 3},
			y:     []float64{0, 0, 0, 0},
			incy:  1,

			want: []float64{2, 0, 4, 6},
		},
		{
			a:     2,
			x:     []float64{1, 2, 3},
			index: []int{0, 2, 3},
			y:     []float64{0, 0, 0, 0, 0, 0, 0, 0},
			incy:  2,

			want: []float64{2, 0, 0, 0, 4, 0, 6, 0},
		},
		{
			a:     2,
			x:     []float64{1, 2, 3},
			index: []int{0, 2, 3},
			y:     []float64{1, 2, 3, 4, 5, 6, 7, 8},
			incy:  2,

			want: []float64{3, 2, 3, 4, 9, 6, 13, 8},
		},
	} {
		Axpy(test.a, test.x, test.index, test.y, test.incy)
		if !reflect.DeepEqual(test.y, test.want) {
			t.Errorf("want = %v, got %v\n", test.want, test.y)
		}
	}
}

func TestGather(t *testing.T) {
	for i, test := range []struct {
		y       []float64
		indices []int

		want []float64
	}{
		{
			y:       []float64{1, 2, 3, 4},
			indices: []int{0, 2, 3},

			want: []float64{1, 3, 4},
		},
		{
			indices: []int{0, 2, 3, 6},
			y:       []float64{1, 2, 3, 4, 5, 6, 7, 8},

			want: []float64{1, 3, 4, 7},
		},
	} {
		y := mat64.NewVector(len(test.y), test.y)
		var x Vector
		Gather(&x, y, test.indices)

		if x.N != y.Len() {
			t.Errorf("%d: wrong dimension, want = %v, got = %v ", i, y.Len(), x.N)
		}
		if !reflect.DeepEqual(x.Data, test.want) {
			t.Errorf("%d: data not equal, want = %v, got %v\n", i, test.want, x.Data)
		}
	}
}

func TestGatherZero(t *testing.T) {
	for j, test := range []struct {
		x, y  []float64
		index []int
		incy  int

		want []float64
	}{
		{
			x:     []float64{math.NaN(), math.NaN(), math.NaN()},
			index: []int{0, 2, 3},
			y:     []float64{1, 2, 3, 4},
			incy:  1,

			want: []float64{1, 3, 4},
		},
		{
			x:     []float64{math.NaN(), math.NaN(), math.NaN()},
			index: []int{0, 2, 3},
			y:     []float64{1, 2, 3, 4, 5, 6, 7, 8},
			incy:  2,

			want: []float64{1, 5, 7},
		},
	} {
		GatherZero(test.y, test.incy, test.x, test.index)
		if !reflect.DeepEqual(test.x, test.want) {
			t.Errorf("want = %v, got %v\n", test.want, test.x)
		}
		for _, idx := range test.index {
			if test.y[idx*test.incy] != 0 {
				t.Errorf("test %d: %d-th element not set to zero", j, idx*test.incy)
			}
		}
	}
}

func TestScatter(t *testing.T) {
	for _, test := range []struct {
		x, y  []float64
		index []int
		incy  int

		want []float64
	}{
		{
			x:     []float64{1, 2, 3},
			index: []int{0, 2, 3},
			y:     []float64{math.NaN(), 0, math.NaN(), math.NaN()},
			incy:  1,

			want: []float64{1, 0, 2, 3},
		},
		{
			x:     []float64{1, 2, 3},
			index: []int{0, 2, 3},
			y:     []float64{math.NaN(), 0, 0, 0, math.NaN(), 0, math.NaN(), 0},
			incy:  2,

			want: []float64{1, 0, 0, 0, 2, 0, 3, 0},
		},
	} {
		Scatter(test.x, test.index, test.y, test.incy)
		if !reflect.DeepEqual(test.y, test.want) {
			t.Errorf("want = %v, got %v\n", test.want, test.y)
		}
	}
}
