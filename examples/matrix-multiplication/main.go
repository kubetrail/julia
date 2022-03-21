package main

import (
	"fmt"
	"log"

	"github.com/kubetrail/julia"
)

func main() {
	julia.Initialize()
	defer julia.Finalize()

	// marshal out a data type that can be passed to julia runtime
	// as an argument to a function call
	n, err := julia.Marshal(int64(5))
	if err != nil {
		log.Fatal(err)
	}

	// evaluate randn with arguments created above to create a
	// random matrix
	x, err := julia.EvalFunc("randn", julia.ModuleBase, n, n)
	if err != nil {
		log.Fatal(err)
	}

	// find inverse of the matrix
	y, err := julia.EvalFunc("inv", julia.ModuleBase, x)
	if err != nil {
		log.Fatal(err)
	}

	// then multiply matrix by its inverse and the result should
	// be very close to an identity matrix
	z, err := julia.EvalFunc("*", julia.ModuleBase, x, y)
	if err != nil {
		log.Fatal(err)
	}

	// prepare a matrix to collect the values from julia runtime
	mat, err := julia.NewMat(make([]float64, 25), 5, 5)
	if err != nil {
		log.Fatal(err)
	}

	// unmarshal julia runtime data into matrix
	if err := julia.Unmarshal(z, mat); err != nil {
		log.Fatal(err)
	}

	// print output
	fmt.Println("matrix multiplied by its inverse is an identity matrix:")
	fmt.Println(mat.Elms)
}
