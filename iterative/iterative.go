// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iterative

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gonum/matrix/mat64"
	"github.com/vladimir-ch/sparse"
)

type Operation uint64

const (
	NoOperation Operation = 0
	ComputeAp   Operation = 1 << (iota - 1)
	SolvePreconditioner
	CheckConvergence
)

type Method interface {
	Init(*Context) Operation
	Iterate(*Context) Operation
}

type Stats struct {
	Iterations         int
	MatVecMultiplies   int
	PrecondionerSolves int
	Residual           float64
	StartTime          time.Time
}

type Result struct {
	X       *mat64.Vector
	Stats   Stats
	Runtime time.Duration
}

type Context struct {
	X        *mat64.Vector
	Residual *mat64.Vector
	P        *mat64.Vector
	Ap       *mat64.Vector
	Z        *mat64.Vector
}

type Settings struct {
	Tolerance  float64
	Iterations int
}

func DefaultSettings(dim int) *Settings {
	return &Settings{
		Tolerance:  1e-6,
		Iterations: 10 * dim,
	}
}

func Solve(a sparse.Matrix, b, xInit *mat64.Vector, settings *Settings, method Method) (result Result, err error) {
	stats := Stats{
		StartTime: time.Now(),
	}

	dim, c := a.Dims()
	if dim != c {
		panic("iterative: matrix is not square")
	}
	if xInit != nil && dim != xInit.Len() {
		panic("iterative: mismatched size of the initial guess")
	}
	if b.Len() != dim {
		panic("iterative: mismatched size of the right-hand side vector")
	}

	if xInit == nil {
		xInit = mat64.NewVector(dim, nil)
	}
	if settings == nil {
		settings = DefaultSettings(dim)
	}

	ctx := Context{
		X:        mat64.NewVector(dim, nil),
		Residual: mat64.NewVector(dim, nil),
	}
	// X = xInit
	ctx.X.CopyVec(xInit)
	if mat64.Norm(ctx.X, math.Inf(1)) > 0 {
		// Residual = Ax
		sparse.MulMatVec(ctx.Residual, 1, false, a, ctx.X)
		stats.MatVecMultiplies++
	}
	// Residual = Ax - b
	ctx.Residual.SubVec(ctx.Residual, b)

	if mat64.Norm(ctx.Residual, 2) >= settings.Tolerance {
		err = iterate(method, a, b, settings, &ctx, &stats)
	}

	result = Result{
		X:       ctx.X,
		Stats:   stats,
		Runtime: time.Since(stats.StartTime),
	}
	return result, err
}

func iterate(method Method, a sparse.Matrix, b *mat64.Vector, settings *Settings, ctx *Context, stats *Stats) error {
	bNorm := mat64.Norm(b, 2)
	if bNorm == 0 {
		bNorm = 1
	}

	op := method.Init(ctx)
	for {
		switch op {
		case NoOperation:

		case ComputeAp:
			ctx.Ap.ScaleVec(0, ctx.Ap)
			sparse.MulMatVec(ctx.Ap, 1, false, a, ctx.P)
			stats.MatVecMultiplies++

		case SolvePreconditioner:
			// TODO(vladimir-ch): Add preconditioners.
			// Z = Residual
			ctx.Z.CopyVec(ctx.Residual)
			stats.PrecondionerSolves++

		case CheckConvergence:
			stats.Iterations++
			stats.Residual = mat64.Norm(ctx.Residual, 2) / bNorm
			fmt.Println(stats.Residual)
			if stats.Residual < settings.Tolerance {
				return nil
			}
			if stats.Iterations == settings.Iterations {
				return errors.New("iterative: reached iteration limit")
			}
		}

		op = method.Iterate(ctx)
	}
}
