// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

import (
	"fmt"

	"github.com/gonum/matrix/mat64"
)

type Index [2]int

type DOK struct {
	rows, cols int
	data       map[Index]float64
	props      MatrixProperties
}

func NewDOK(r, c int) *DOK {
	return &DOK{
		rows: r,
		cols: c,
		data: make(map[Index]float64),
	}
}

func (m *DOK) Dims() (r, c int) {
	return m.rows, m.cols
}

func (m *DOK) At(r, c int) float64 {
	if r >= m.rows || r < 0 {
		panic("sparse: row index out of range")
	}
	if c >= m.cols || c < 0 {
		panic("sparse: column index out of range")
	}

	return m.data[Index{r, c}]
}

func (m *DOK) Properties() MatrixProperties {
	return m.props
}

func (m *DOK) SetSparse(r, c int, v float64) {
	if r >= m.rows || r < 0 {
		panic("sparse: row index out of range")
	}
	if c >= m.cols || c < 0 {
		panic("sparse: column index out of range")
	}

	if _, exists := m.data[Index{r, c}]; !exists {
		panic(fmt.Sprintf("sparse: entry at (%d,%d) does not exist", r, c))
	}
	m.data[Index{r, c}] = v
}

func (m *DOK) InsertEntry(r, c int, v float64) {
	if r >= m.rows || r < 0 {
		panic("sparse: row index out of range")
	}
	if c >= m.cols || c < 0 {
		panic("sparse: column index out of range")
	}

	m.data[Index{r, c}] = v
}

func (m *DOK) Triplets() []Triplet {
	var t []Triplet
	for k, v := range m.data {
		t = append(t, Triplet{k[0], k[1], v})
	}
	return t
}

func dokMulMatVec(y *mat64.Vector, alpha float64, transA bool, a *DOK, x *mat64.Vector) {
	r, c := a.Dims()
	if transA {
		if r != x.Len() || c != y.Len() {
			panic("sparse: dimension mismatch")
		}
	} else {
		if r != y.Len() || c != x.Len() {
			panic("sparse: dimension mismatch")
		}
	}

	if alpha == 0 {
		return
	}

	xRaw := x.RawVector()
	yRaw := y.RawVector()
	if transA {
		for ij, aij := range a.data {
			yRaw.Data[ij[1]*yRaw.Inc] += alpha * aij * xRaw.Data[ij[0]*xRaw.Inc]
		}
	} else {
		for ij, aij := range a.data {
			yRaw.Data[ij[0]*yRaw.Inc] += alpha * aij * xRaw.Data[ij[1]*xRaw.Inc]
		}
	}
}
