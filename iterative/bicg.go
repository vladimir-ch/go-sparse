// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iterative

import (
	"math"

	"github.com/gonum/matrix/mat64"
)

// BiCG implements the Bi-Conjugate Gradient iterative method with
// preconditioning for solving the linear system Ax = b.
type BiCG struct {
	BreakdownTolerance float64

	resume int
	rho    float64
}

func (bicg *BiCG) Init(ctx *Context) Operation {
	if bicg.BreakdownTolerance == 0 {
		bicg.BreakdownTolerance = 1e-6
	}
	bicg.rho = math.NaN()

	dim := ctx.X.Len()
	if ctx.P == nil || ctx.P.Len() != dim {
		ctx.P = mat64.NewVector(dim, nil)
	}
	if ctx.Ap == nil || ctx.Ap.Len() != dim {
		ctx.Ap = mat64.NewVector(dim, nil)
	}
	if ctx.Q == nil || ctx.Q.Len() != dim {
		ctx.Q = mat64.NewVector(dim, nil)
	}
	if ctx.Aq == nil || ctx.Aq.Len() != dim {
		ctx.Aq = mat64.NewVector(dim, nil)
	}
	if ctx.Z == nil || ctx.Z.Len() != dim {
		ctx.Z = mat64.NewVector(dim, nil)
	}

	bicg.resume = 2
	return SolvePreconditioner
	// Solve M z = r_{i-1}
}

func (bicg *BiCG) Iterate(ctx *Context) Operation {
	switch bicg.resume {
	case 1:
		cg.resume = 2
		return SolvePreconditioner
		// Solve M z = r_{i-1}
	case 2:
		// ρ_i = r_{i-1} · z
		cg.rho = mat64.Dot(ctx.Residual, ctx.Z)
	default:
		panic("unreachable")
	}
}
