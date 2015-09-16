// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

// MulMatVec multiplies the dense vector x by a sparse matrix A (or its
// transpose) and adds the result to the dense vector y, i.e., it computes
//
//  y = alpha * op(A) * x + beta * y,
//
// where op(A) is either A or Aᵀ.
func MulMatVec(alpha float64, transA bool, a Matrix, x []float64, incx int, beta float64, y []float64, incy int) {
	switch a := a.(type) {
	case *CSR:
		csrMulMatVec(alpha, transA, a, x, incx, beta, y, incy)
	case *DOK:
		dokMulMatVec(alpha, transA, a, x, incx, beta, y, incy)
	default:
		panic("unsupported matrix type")
	}
}
