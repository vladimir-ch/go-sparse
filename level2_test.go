// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

import (
	"reflect"
	"testing"
)

func TestMulMatVec(t *testing.T) {
	for id, test := range []struct {
		r, c int
		i, j []int
		v    []float64

		x     []float64
		incx  int
		incy  int
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
			incx:  1,
			incy:  1,
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
			incx:  1,
			incy:  1,
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
			incx:  1,
			incy:  1,
			alpha: 1,
			trans: false,

			want: []float64{2, -10, 9},
		},
		{
			r: 3,
			c: 3,
			i: []int{0, 1, 1, 2},      // 0 1  0
			j: []int{1, 0, 2, 2},      // 2 0 -4
			v: []float64{1, 2, -4, 3}, // 0 0  3

			x:     []float64{1, 2, 3},
			incx:  1,
			incy:  1,
			alpha: 1,
			trans: true,

			want: []float64{4, 1, 1},
		},
		{
			r: 3,
			c: 3,
			i: []int{0, 1, 1, 2},      // 0 1  0
			j: []int{1, 0, 2, 2},      // 2 0 -4
			v: []float64{1, 2, -4, 3}, // 0 0  3

			x:     []float64{1, 2, 3},
			incx:  1,
			incy:  1,
			alpha: 2,
			trans: true,

			want: []float64{8, 2, 2},
		},
	} {
		dok := NewDOK(test.r, test.c)
		for i := 0; i < len(test.v); i++ {
			dok.InsertEntry(test.i[i], test.j[i], test.v[i])
		}
		y := make([]float64, test.r*test.incy)

		MulMatVec(test.alpha, test.trans, dok, test.x, test.incx, 1, y, test.incy)
		if !reflect.DeepEqual(y, test.want) {
			t.Errorf("test %d: unexpected result for DOK", id+1)
		}

		for i := range y {
			y[i] = 0
		}
		csr := NewCSR(dok)
		MulMatVec(test.alpha, test.trans, csr, test.x, test.incx, 1, y, test.incy)
		if !reflect.DeepEqual(y, test.want) {
			t.Errorf("test %d: unexpected result for CSR", id+1)
		}
	}
}
