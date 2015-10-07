// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

import (
	"reflect"
	"testing"

	"github.com/gonum/matrix/mat64"
)

func TestMulMatVec(t *testing.T) {
	for id, test := range []struct {
		r, c int
		i, j []int
		v    []float64

		x     []float64
		alpha float64
		trans bool

		want []float64
	}{
		{
			r: 2,
			c: 2,
			i: []int{0, 1}, // 0 1
			j: []int{1, 0}, // 2 0
			v: []float64{1, 2},

			x:     []float64{0, 0},
			alpha: 1,
			trans: false,

			want: []float64{0, 0},
		},
		{
			r: 2,
			c: 2,
			i: []int{0, 1}, // 0 1
			j: []int{1, 0}, // 2 0
			v: []float64{1, 2},

			x:     []float64{1, 1},
			alpha: 1,
			trans: false,

			want: []float64{1, 2},
		},
		{
			r: 3,
			c: 3,
			i: []int{0, 1, 1, 2},      // 0 1  0
			j: []int{1, 0, 2, 2},      // 2 0 -4
			v: []float64{1, 2, -4, 3}, // 0 0  3

			x:     []float64{1, 2, 3},
			alpha: 1,
			trans: false,

			want: []float64{2, -10, 9},
		},
		{
			r: 3,
			c: 3,
			i: []int{0, 1, 1, 2},      // 0 1  0  ^T
			j: []int{1, 0, 2, 2},      // 2 0 -4
			v: []float64{1, 2, -4, 3}, // 0 0  3

			x:     []float64{1, 2, 3},
			alpha: 1,
			trans: true,

			want: []float64{4, 1, 1},
		},
		{
			r: 3,
			c: 3,
			i: []int{0, 1, 1, 2},      // 0 1  0  ^T
			j: []int{1, 0, 2, 2},      // 2 0 -4
			v: []float64{1, 2, -4, 3}, // 0 0  3

			x:     []float64{1, 2, 3},
			alpha: 2,
			trans: true,

			want: []float64{8, 2, 2},
		},
		{
			r: 3,
			c: 5,
			i: []int{0, 0, 1, 1, 2, 2},       // 0 1  0 1  0
			j: []int{1, 3, 0, 2, 2, 4},       // 2 0 -4 0  0
			v: []float64{1, 1, 2, -4, 3, -5}, // 0 0  3 0 -5

			x:     []float64{1, 2, 3, 4, 5},
			alpha: 2,
			trans: false,

			want: []float64{12, -20, -32},
		},
		{
			r: 3,
			c: 5,
			i: []int{0, 0, 1, 1, 2, 2},       // 0 1  0 1  0  ^T
			j: []int{1, 3, 0, 2, 2, 4},       // 2 0 -4 0  0
			v: []float64{1, 1, 2, -4, 3, -5}, // 0 0  3 0 -5

			x:     []float64{1, 2, 3},
			alpha: 2,
			trans: true,

			want: []float64{8, 2, 2, 2, -30},
		},
	} {
		var r, c int
		if test.trans {
			r, c = test.c, test.r
		} else {
			r, c = test.r, test.c
		}
		x := mat64.NewVector(c, test.x)
		y := mat64.NewVector(r, nil)

		dok := NewDOK(test.r, test.c)
		for i := 0; i < len(test.v); i++ {
			dok.InsertEntry(test.i[i], test.j[i], test.v[i])
		}

		MulMatVec(y, test.alpha, test.trans, dok, x)
		if !reflect.DeepEqual(y.RawVector().Data, test.want) {
			t.Errorf("test %d: unexpected result for DOK, want = %v, got = %v", id+1, test.want, y.RawVector().Data)
		}

		for i := 0; i < y.Len(); i++ {
			y.SetVec(i, 0)
		}
		csr := NewCSR(dok)

		MulMatVec(y, test.alpha, test.trans, csr, x)
		if !reflect.DeepEqual(y.RawVector().Data, test.want) {
			t.Errorf("test %d: unexpected result for CSR, want = %v, got = %v", id+1, test.want, y.RawVector().Data)
		}
	}
}
