// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

func Dot(x []float64, index []int, y []float64, incy int) float64 {
	if len(x) != len(index) {
		panic("sparse: slice length mismatch")
	}

	var r float64
	for i, xi := range x {
		r += xi * y[index[i]*incy]
	}
	return r
}

func Axpy(a float64, x []float64, index []int, y []float64, incy int) {
	if len(x) != len(index) {
		panic("sparse: slice length mismatch")
	}

	for i, xi := range x {
		y[index[i]*incy] += a * xi
	}
}

func Gather(y []float64, incy int, x []float64, index []int) {
	if len(x) != len(index) {
		panic("sparse: slice length mismatch")
	}

	for i := range x {
		x[i] = y[index[i]*incy]
	}
}

func GatherZero(y []float64, incy int, x []float64, index []int) {
	if len(x) != len(index) {
		panic("sparse: slice length mismatch")
	}

	for i := range x {
		x[i] = y[index[i]*incy]
		y[index[i]*incy] = 0
	}
}

func Scatter(x, y []float64, incy int, index []int) {
	if len(x) != len(index) {
		panic("sparse: slice length mismatch")
	}

	for i, xi := range x {
		y[index[i]*incy] = xi
	}
}
