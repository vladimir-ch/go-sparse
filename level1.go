// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

import "github.com/gonum/matrix/mat64"

// Dot computes the dot product of the sparse vector x with the dense vector y.
// The vectors must have the same dimension.
func Dot(x *Vector, y *mat64.Vector) (dot float64) {
	if x.N != y.Len() {
		panic("sparse: vector dimension mismatch")
	}

	raw := y.RawVector()
	for i, index := range x.Indices {
		dot += x.Data[i] * raw.Data[index*raw.Inc]
	}
	return
}

// Axpy scales the sparse vector x by alpha and adds the result to the dense
// vector y, i.e., it computes
//
//  y[index[i]*incy] += alpha*x[i]
//
// If alpha is zero, y is not modified.
func Axpy(alpha float64, x []float64, index []int, y []float64, incy int) {
	if len(x) != len(index) {
		panic("sparse: slice length mismatch")
	}

	if alpha == 0 {
		return
	}
	for i, idx := range index {
		y[idx*incy] += alpha * x[i]
	}
}

// Gather gathers entries given by indices of the dense vector y into the sparse
// vector x. Indices must not be nil.
func Gather(x *Vector, y *mat64.Vector, indices []int) {
	if indices == nil {
		panic("sparse: slice is nil")
	}

	x.reuseAs(y.Len(), len(indices))
	copy(x.Indices, indices)
	raw := y.RawVector()
	for i, index := range x.Indices {
		x.Data[i] = raw.Data[index*raw.Inc]
	}
}

// Gather gathers entries given by indices of the dense vector y into the sparse
// vector x and sets the corresponding values of y to zero.
func GatherZero(x *Vector, y *mat64.Vector, indices []int) {
	if indices == nil {
		panic("sparse: slice is nil")
	}

	x.reuseAs(y.Len(), len(indices))
	copy(x.Indices, indices)
	raw := y.RawVector()
	for i, index := range x.Indices {
		x.Data[i] = raw.Data[index*raw.Inc]
		raw.Data[index*raw.Inc] = 0
	}
}

// Scatter copies the values of x into the corresponding locations in the dense
// vector y, i.e., it assigns
//
//  y[index[i]*incy] = x[i]
//
func Scatter(x []float64, index []int, y []float64, incy int) {
	if len(x) != len(index) {
		panic("sparse: slice length mismatch")
	}

	for i, idx := range index {
		y[idx*incy] = x[i]
	}
}
