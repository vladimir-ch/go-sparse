// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

// Vector is a sparse vector represented by a slice of non-zero values and a
// slice denoting their indices.
type Vector struct {
	N       int       // Dimension of the vector.
	Data    []float64 // Non-zero values.
	Indices []int     // Indices of values in Data. Must be zero-based and unique.
}

// NewVector returns a new Vector of dimension n with non-zero elements given
// by data and indices. Both data and indices must have the same length smaller
// than n, otherwise NewVector will panic. Indices must be unique, although no
// checking is done.
func NewVector(n int, data []float64, indices []int) *Vector {
	if len(data) != len(indices) {
		panic("sparse: slice length mismatch")
	}
	if n < len(data) {
		panic("sparse: vector dimension is less than the number of entries")
	}
	return &Vector{
		N:       n,
		Data:    data,
		Indices: indices,
	}
}

// InsertEntry appends the value v with index i to the Vector.
func (v *Vector) InsertEntry(val float64, i int) {
	v.Data = append(v.Data, val)
	v.Indices = append(v.Indices, i)
}

func (v *Vector) reuseAs(n, nnz int) {
	v.N = n
	if cap(v.Data) >= nnz {
		v.Data = v.Data[:nnz]
		v.Indices = v.Indices[:nnz]
	} else {
		v.Data = make([]float64, nnz)
		v.Indices = make([]int, nnz)
	}
}
