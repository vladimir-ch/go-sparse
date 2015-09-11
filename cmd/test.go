package main

import (
	"fmt"

	"github.com/vladimir-ch/sparse"
)

func main() {
	m := sparse.NewDOK(5, 5)
	m.Set(0, 0, 1)
	m.Set(0, 1, -1)
	m.Set(0, 3, -3)

	m.Set(1, 0, -2)
	m.Set(1, 1, 5)

	m.Set(2, 2, 4)
	m.Set(2, 3, 6)
	m.Set(2, 4, 4)

	m.Set(3, 0, -4)
	m.Set(3, 2, 2)
	m.Set(3, 3, 7)

	m.Set(4, 1, 8)
	m.Set(4, 4, -5)

	c := sparse.NewCSR(m)

	y := make([]float64, 5)
	x := []float64{1, 2, 3, 4, 5}

	sparse.MulMatVec(1, c, x, 1, y, 1)

	fmt.Println(y)
}
