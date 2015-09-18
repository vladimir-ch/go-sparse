// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

import "fmt"

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

func dokMulMatVec(alpha float64, transA bool, a *DOK, x []float64, incx int, beta float64, y []float64, incy int) {
	r, _ := a.Dims()
	if beta == 0 {
		for i := 0; i < r; i++ {
			y[i*incy] = 0
		}
	} else {
		for i := 0; i < r; i++ {
			y[i*incy] *= beta
		}
	}
	if transA {
		for k, v := range a.data {
			y[k[1]*incy] += alpha * v * x[k[0]*incx]
		}
	} else {
		for k, v := range a.data {
			y[k[0]*incy] += alpha * v * x[k[1]*incx]
		}
	}
}
