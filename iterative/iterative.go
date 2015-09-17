// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iterative

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gonum/floats"
	"github.com/vladimir-ch/sparse"
)

type RequestType uint64

const (
	NoRequest RequestType = 0
	ComputeAp RequestType = 1 << (iota - 1)
	SolvePreconditioner
	CheckConvergence
)

type Method interface {
	Init(*Context) RequestType
	Iterate(*Context) RequestType
}

type Stats struct {
	Iterations         int
	MatVecMultiplies   int
	PrecondionerSolves int
	Residual           float64
	StartTime          time.Time
}

type Result struct {
	X       []float64
	Stats   Stats
	Runtime time.Duration
}

type Context struct {
	X        []float64
	Residual []float64
	P        []float64
	Ap       []float64
	Z        []float64
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

func Solve(a sparse.Matrix, b, xInit []float64, settings *Settings, method Method) (result Result, err error) {
	stats := Stats{
		StartTime: time.Now(),
	}

	dim := len(xInit)
	if dim == 0 {
		panic("iterative: invalid dimension")
	}

	r, c := a.Dims()
	if r != c {
		panic("iterative: matrix is not square")
	}
	if c != dim {
		panic("iterative: mismatched size of the matrix")
	}
	if len(b) != dim {
		panic("iterative: mismatched size of the right-hand side vector")
	}

	if settings == nil {
		settings = DefaultSettings(dim)
	}

	ctx := Context{
		X:        make([]float64, dim),
		Residual: make([]float64, dim),
	}
	copy(ctx.X, xInit)
	copy(ctx.Residual, b)
	if floats.Norm(ctx.X, math.Inf(1)) > 0 {
		sparse.MulMatVec(-1, false, a, ctx.X, 1, 1, ctx.Residual, 1)
		stats.MatVecMultiplies++
	}

	if floats.Norm(ctx.Residual, 2) >= settings.Tolerance {
		err = iterate(method, a, b, settings, &ctx, &stats)
	}

	result = Result{
		X:       ctx.X,
		Stats:   stats,
		Runtime: time.Since(stats.StartTime),
	}
	return result, err
}

func iterate(method Method, a sparse.Matrix, b []float64, settings *Settings, ctx *Context, stats *Stats) error {
	dim := len(ctx.X)
	bNorm := floats.Norm(b, 2)
	if bNorm == 0 {
		bNorm = 1
	}

	request := method.Init(ctx)
	for {
		switch request {
		case NoRequest:

		case ComputeAp:
			ctx.Ap = resize(ctx.Ap, dim)
			sparse.MulMatVec(1, false, a, ctx.P, 1, 0, ctx.Ap, 1)
			stats.MatVecMultiplies++

		case SolvePreconditioner:
			ctx.Z = resize(ctx.Z, dim)
			copy(ctx.Z, ctx.Residual)
			stats.PrecondionerSolves++

		case CheckConvergence:
			stats.Iterations++
			stats.Residual = floats.Norm(ctx.Residual, 2) / bNorm
			fmt.Println(stats.Residual)
			if stats.Residual < settings.Tolerance {
				return nil
			}
			if stats.Iterations == settings.Iterations {
				return errors.New("iterative: reached iteration limit")
			}
		}

		request = method.Iterate(ctx)
	}
}

// resize resizes x to the length dim, reusing its memory if possible.
func resize(x []float64, dim int) []float64 {
	if cap(x) >= dim {
		return x[:dim]
	}
	return make([]float64, dim)
}
