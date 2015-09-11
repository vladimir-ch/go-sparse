// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

type Vector struct {
	n      int
	values []float64
	index  []int
}

func NewVector(n int, values []float64, indices []int) *Vector {
	if len(values) != len(indices) {
		panic("slice length mismatch")
	}
	if n < len(values) {
		panic("vector dimension is less than the number of entries")
	}
	return &Vector{
		n:      n,
		values: values,
		index:  indices,
	}
}

func (v *Vector) Dim() int {
	return v.n
}

func (v *Vector) NonZeros() int {
	return len(v.values)
}

func (v *Vector) InsertEntry(value float64, index int) {
	v.values = append(v.values, value)
	v.index = append(v.index, index)
}
