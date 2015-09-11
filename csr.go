// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

import "sort"

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
	for j := m.rowIndex[r]; j < m.rowIndex[r+1]; j++ {
		if m.columns[j] == c {
			return m.values[j]
		}
	}
	return 0
}

func csrMulMatVec(alpha float64, a *CSR, x []float64, incx int, y []float64, incy int) {
	r, _ := a.Dims()
	for i := 0; i < r; i++ {
		sum := y[i*incy]
		for j := a.rowIndex[i]; j < a.rowIndex[i+1]; j++ {
			sum += alpha * a.values[j] * x[a.columns[j]*incx]
		}
		y[i*incy] = sum
	}
}
