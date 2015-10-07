// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

import "github.com/gonum/matrix/mat64"

// MulMatVec multiplies the dense vector x by a sparse matrix A (or its
// transpose) and adds the result to the dense vector y, i.e., it computes
//
//  y += alpha * op(A) * x,
//
// where op(A) is either A or Aᵀ.
func MulMatVec(y *mat64.Vector, alpha float64, transA bool, a Matrix, x *mat64.Vector) {
	switch a := a.(type) {
	case *CSR:
		csrMulMatVec(y, alpha, transA, a, x)
	case *DOK:
		dokMulMatVec(y, alpha, transA, a, x)
	default:
		panic("unsupported matrix type")
	}
}
