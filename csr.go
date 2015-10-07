// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

import (
	"sort"

	"github.com/gonum/matrix/mat64"
)

type CSR struct {
	rows, cols int

	values   []float64
	columns  []int
	rowIndex []int
}

func NewCSR(dok *DOK) *CSR {
	triplets := dok.Triplets()
	nnz := len(triplets)

	// Triplets from DOK are unique, but not sorted. Alternatively, we could
	// have something like SortIndices() method to turn the matrix into the
	// canonical form.
	sort.Sort(rowWise(triplets))

	rows, cols := dok.Dims()
	values := make([]float64, nnz)
	columns := make([]int, nnz)
	rowIndex := make([]int, rows+1)

	// Count the number of entries in each row.
	for i := range triplets {
		rowIndex[triplets[i].Row]++
	}

	// Cumulative sum of entries per row.
	for i, sum := 0, 0; i < rows; i++ {
		tmp := rowIndex[i]
		rowIndex[i] = sum
		sum += tmp
	}
	rowIndex[rows] = nnz

	offset := make([]int, rows) // Instead of allocating we could modify the rowIndex slice.
	for _, t := range triplets {
		dest := rowIndex[t.Row] + offset[t.Row]
		columns[dest] = t.Col
		values[dest] = t.Value
		offset[t.Row]++
	}

	return &CSR{
		rows:     rows,
		cols:     cols,
		values:   values,
		columns:  columns,
		rowIndex: rowIndex,
	}
}

func (m *CSR) Dims() (r, c int) {
	return m.rows, m.cols
}

func (m *CSR) At(r, c int) float64 {
	if r >= m.rows || r < 0 {
		panic("sparse: row index out of range")
	}
	if c >= m.cols || c < 0 {
		panic("sparse: column index out of range")
	}

	for j := m.rowIndex[r]; j < m.rowIndex[r+1]; j++ {
		if m.columns[j] == c {
			return m.values[j]
		}
	}
	return 0
}

func csrMulMatVec(y *mat64.Vector, alpha float64, transA bool, a *CSR, x *mat64.Vector) {
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

	yRaw := y.RawVector()
	if transA {
		row := Vector{N: y.Len()}
		for i := 0; i < r; i++ {
			start := a.rowIndex[i]
			end := a.rowIndex[i+1]
			row.Data = a.values[start:end]
			row.Indices = a.columns[start:end]
			Axpy(y, alpha*x.At(i, 0), &row)
		}
	} else {
		row := Vector{N: x.Len()}
		for i := 0; i < r; i++ {
			start := a.rowIndex[i]
			end := a.rowIndex[i+1]
			row.Data = a.values[start:end]
			row.Indices = a.columns[start:end]
			yRaw.Data[i*yRaw.Inc] += alpha * Dot(&row, x)
		}
	}
}
