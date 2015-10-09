// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iterative

import (
	"math"

	"github.com/gonum/matrix/mat64"
)

type CG struct {
	first     bool
	resume    int
	rho, rho1 float64
}

func (cg *CG) Init(ctx *Context) Operation {
	cg.first = true
	cg.rho = math.NaN()
	cg.rho1 = math.NaN()

	dim := ctx.X.Len()
	if ctx.P == nil || ctx.P.Len() != dim {
		ctx.P = mat64.NewVector(dim, nil)
	}
	if ctx.Ap == nil || ctx.Ap.Len() != dim {
		ctx.Ap = mat64.NewVector(dim, nil)
	}
	if ctx.Z == nil || ctx.Z.Len() != dim {
		ctx.Z = mat64.NewVector(dim, nil)
	}

	cg.resume = 2
	return SolvePreconditioner
	// Solve M z = r_{i-1}
}

func (cg *CG) Iterate(ctx *Context) Operation {
	switch cg.resume {
	case 1:
		cg.resume = 2
		return SolvePreconditioner
		// Solve M z = r_{i-1}
	case 2:
		// ρ_i = r_{i-1} · z
		cg.rho = mat64.Dot(ctx.Residual, ctx.Z)
		if !cg.first {
			// β = ρ_i / ρ_{i-1}
			beta := cg.rho / cg.rho1
			// z = z + β p_{i-1}
			ctx.Z.AddScaledVec(ctx.Z, beta, ctx.P)
		}
		cg.first = false
		// p_i = z
		ctx.P.CopyVec(ctx.Z)

		cg.resume = 3
		return ComputeAp
		// Compute Ap
	case 3:
		// α = ρ_i / (p_i · Ap_i)
		alpha := cg.rho / mat64.Dot(ctx.P, ctx.Ap)
		// x_i = x_{i-1} + α p_i
		ctx.X.AddScaledVec(ctx.X, alpha, ctx.P)
		// r_i = r_{i-1} - α Ap_i
		ctx.Residual.AddScaledVec(ctx.Residual, -alpha, ctx.Ap)

		cg.rho1 = cg.rho

		cg.resume = 1
		return CheckConvergence
	}
	panic("unreachable")
}
