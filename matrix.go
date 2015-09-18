// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

// Matrix is a sparse matrix.
type Matrix interface {
	// Dims returns the dimensions of the matrix.
	Dims() (r, c int)

	// At returns the value of the matrix entry at (r, c). It will panic if r
	// or c are out of bounds for the matrix.
	At(r, c int) float64
}

// MutableMatrix is a matrix that can modify its non-zero entries without
// changing its sparsity structure.
type MutableMatrix interface {
	Matrix

	// SetSparse sets the entry at (r, c) to v. It will panic if r or c are out
	// of bounds for the matrix.
	// Whether SetSparse panics if the entry does not already exist in the
	// matrix is implementation-specific.
	SetSparse(r, c int, v float64)
}

// MatrixBuilder can build a sparse matrix by modifying its sparsity structure.
type MatrixBuilder interface {
	Begin()
	InsertEntry(r, c int, v float64)
	InsertEntries()
	InsertClique()
	End()
}
