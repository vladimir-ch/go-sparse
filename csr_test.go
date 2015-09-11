// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

import (
	"reflect"
	"testing"
)

func TestCSR(t *testing.T) {
	for _, test := range []struct {
		r, c int
		i, j []int
		v    []float64
	}{
		{
			r: 2,
			c: 2,
			i: []int{0, 1},
			j: []int{1, 0},
			v: []float64{1, 2},
		},
		{
			r: 4,
			c: 2,
			i: []int{0, 1, 3, 2, 3},
			j: []int{1, 0, 0, 2, 1},
			v: []float64{1, 2, 3, 4, 5},
		},
	} {
		dok := NewDOK(test.r, test.c)
		for i := 0; i < len(test.v); i++ {
			dok.Set(test.i[i], test.j[i], test.v[i])
		}

		csr := NewCSR(dok)

		for i := 0; i < len(test.v); i++ {
			v := csr.At(test.i[i], test.j[i])
			if v != test.v[i] {
				t.Errorf("entries not equal at (%d,%d): want %v, got %v\n", test.i[i], test.j[i], test.v[i], v)
			}
		}
	}
}

func TestCSRMulMatVec(t *testing.T) {
	for id, test := range []struct {
		r, c int
		i, j []int
		v    []float64

		x     []float64
		incx  int
		incy  int
		alpha float64

		want []float64
	}{
		{
			r:     2,
			c:     2,
			i:     []int{0, 1},
			j:     []int{1, 0},
			v:     []float64{1, 2},
			x:     []float64{1, 1},
			incx:  1,
			incy:  1,
			alpha: 1,
			want:  []float64{1, 2},
		},
		{
			r:     2,
			c:     2,
			i:     []int{0, 1},
			j:     []int{1, 0},
			v:     []float64{1, 2},
			x:     []float64{0, 0},
			incx:  1,
			incy:  1,
			alpha: 1,
			want:  []float64{0, 0},
		},
		{
			r:     3,
			c:     3,
			i:     []int{0, 1, 2},
			j:     []int{1, 0, 2},
			v:     []float64{1, 2, 3},
			x:     []float64{1, 2, 3},
			incx:  1,
			incy:  1,
			alpha: 1,
			want:  []float64{2, 2, 9},
		},
	} {
		dok := NewDOK(test.r, test.c)
		for i := 0; i < len(test.v); i++ {
			dok.Set(test.i[i], test.j[i], test.v[i])
		}

		csr := NewCSR(dok)

		y := make([]float64, test.r*test.incy)
		csrMulMatVec(test.alpha, csr, test.x, test.incx, y, test.incy)

		if !reflect.DeepEqual(y, test.want) {
			t.Errorf("test %d: unexpected result", id)
		}
	}
}
