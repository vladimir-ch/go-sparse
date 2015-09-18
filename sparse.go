// Copyright 2015 Vladimír Chalupecký. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparse

type MatrixProperties struct {
	Symmetric       bool
	LowerTriangular bool
	UpperTriangular bool
}

type Triplet struct {
	Row, Col int
	Value    float64
}

type rowWise []Triplet

func (r rowWise) Len() int      { return len(r) }
func (r rowWise) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r rowWise) Less(i, j int) bool {
	return r[i].Row < r[j].Row || (r[i].Row == r[j].Row && r[i].Col < r[j].Col)
}
