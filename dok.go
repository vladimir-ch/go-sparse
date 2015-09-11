// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

type DOK struct {
	rows, cols int

	data map[[2]int]float64
}

func NewDOK(r, c int) *DOK {
	return &DOK{
		rows: r,
		cols: c,
		data: make(map[[2]int]float64),
	}
}

func (m *DOK) Dims() (r, c int) {
	return m.rows, m.cols
}

func (m *DOK) At(r, c int) float64 {
	return m.data[[2]int{r, c}]
}

func (m *DOK) Set(r, c int, v float64) {
	m.data[[2]int{r, c}] = v
}

func (m *DOK) Add(r, c int, v float64) {
	m.data[[2]int{r, c}] += v
}

func (m *DOK) Triplets() []Triplet {
	var t []Triplet
	for k, v := range m.data {
		t = append(t, Triplet{k[0], k[1], v})
	}
	return t
}

func dokMulMatVec(alpha float64, a *DOK, x []float64, incx int, y []float64, incy int) {
	for k, v := range a.data {
		y[k[0]*incy] += alpha * v * x[k[1]*incx]
	}
}