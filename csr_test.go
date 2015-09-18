// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

import "testing"

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
			c: 3,
			i: []int{0, 1, 3, 2, 3},
			j: []int{1, 0, 0, 2, 1},
			v: []float64{1, 2, 3, 4, 5},
		},
	} {
		dok := NewDOK(test.r, test.c)
		for i := 0; i < len(test.v); i++ {
			dok.InsertEntry(test.i[i], test.j[i], test.v[i])
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
