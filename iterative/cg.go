// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iterative

import (
	"math"

	"github.com/gonum/floats"
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

	ctx.P = resize(ctx.P, len(ctx.X))

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
		cg.rho = floats.Dot(ctx.Residual, ctx.Z)
		if !cg.first {
			// β = ρ_i / ρ_{i-1}
			beta := cg.rho / cg.rho1
			// z = z + β p_{i-1}
			floats.AddScaled(ctx.Z, beta, ctx.P)
		}
		cg.first = false
		// p_i = z
		copy(ctx.P, ctx.Z)

		cg.resume = 3
		return ComputeAp
		// Compute Ap
	case 3:
		// α = ρ_i / (p_i · Ap_i)
		alpha := cg.rho / floats.Dot(ctx.P, ctx.Ap)
		// x_i = x_{i-1} + α p_i
		floats.AddScaled(ctx.X, alpha, ctx.P)
		// r_i = r_{i-1} - α Ap_i
		floats.AddScaled(ctx.Residual, -alpha, ctx.Ap)

		cg.rho1 = cg.rho

		cg.resume = 1
		return CheckConvergence
	}
	panic("unreachable")
}
