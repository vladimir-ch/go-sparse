package main

import (
	"bufio"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/vladimir-ch/sparse"
	"github.com/vladimir-ch/sparse/iterative"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("missing file name")
	}
	name := flag.Args()[0]

	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var r io.Reader
	if path.Ext(name) == ".gz" {
		gz, err := gzip.NewReader(f)
		if err != nil {
			log.Fatal(err)
		}
		name = strings.TrimSuffix(name, ".gz")
		r = gz
	} else {
		r = f
	}

	var aDok *sparse.DOK
	switch path.Ext(name) {
	case ".mtx":
		aDok, err = readMatrixMarket(r)
	case ".rsa":
		log.Fatal("reading of Harwell-Boeing format not yet implemented")
	default:
		log.Fatal("unknown file extension")
	}
	if err != nil {
		log.Fatal(err)
	}

	a := sparse.NewCSR(aDok)
	n, _ := a.Dims()
	x := make([]float64, n)
	for i := range x {
		x[i] = 1
	}
	b := make([]float64, n)
	sparse.MulMatVec(1, false, a, x, 1, 0, b, 1)
	for i := range x {
		x[i] = 0
	}
	result, err := iterative.Solve(a, b, x, nil, &iterative.CG{})
	if err != nil {
		log.Fatal(err)
	}
	if len(result.X) > 10 {
		fmt.Println("Solution[:10]:", result.X[:10])
	} else {
		fmt.Println("Solution:", result.X)
	}
}

func readMatrixMarket(r io.Reader) (*sparse.DOK, error) {
	s := bufio.NewScanner(r)
	s.Scan()
	line := s.Text()
	if line != "%%MatrixMarket matrix coordinate real symmetric" {
		return nil, errors.New("matrix not symmetric")
	}

	for s.Scan() {
		line = s.Text()
		if !strings.HasPrefix(line, "%") {
			break
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	fields := strings.Fields(line)
	rows, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}
	cols, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, err
	}
	if rows != cols {
		return nil, errors.New("matrix is not square")
	}
	nnz, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, err
	}

	a := sparse.NewDOK(rows, cols)
	var count int
	for s.Scan() {
		line = s.Text()
		fields := strings.Fields(line)

		i, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		j, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		v, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, err
		}

		a.Set(i-1, j-1, v)
		if i != j {
			a.Set(j-1, i-1, v)
		}
		count++
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	if count != nnz {
		return nil, errors.New("mismatched number of non-zeros")
	}

	return a, nil
}
