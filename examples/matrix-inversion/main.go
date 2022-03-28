package main

import (
	"fmt"
	"log"

	"github.com/kubetrail/julia"
)

func main() {
	julia.Initialize()
	defer julia.Finalize()

	// testing inv(5,5)
	n := 5

	// create an argument to pass to julia function in the form
	// of a scalar value of int64 data type
	arg, err := julia.Marshal(int64(n))
	if err != nil {
		log.Fatal(err)
	}

	// get response from julia after evaluating a function called randn
	// in base module and passing two arguments to it of same value
	// created in the previous step
	resp, err := julia.EvalFunc("randn", julia.ModuleBase, arg, arg)
	if err != nil {
		log.Fatal(err)
	}

	// initialize a new matrix to populate it with response received above.
	// create a float64 matrix since output from randn is a float64 data type
	mat, err := julia.NewMat(make([]float64, n*n), n, n)
	if err != nil {
		log.Fatal(err)
	}

	// unmarshal response into matrix. it is important to unmarshal into
	// types that match exactly those returned by julia, otherwise you will
	// get segfaults
	if err := julia.Unmarshal(resp, mat); err != nil {
		log.Fatal(err)
	}

	// print matrix elements
	fmt.Println("rand mat:", mat.elms)

	// now pass this matrix back to julia to compute its inverse.
	// marshaling is required to pass any data to julia runtime
	data, err := julia.Marshal(mat)
	if err != nil {
		log.Fatal(err)
	}

	// obtain response from julia runtime
	resp, err = julia.EvalFunc("inv", julia.ModuleBase, data)
	if err != nil {
		log.Fatal(err)
	}

	// unmrshal response back into the matrix since it fits perfect
	// for the data type and size
	if err := julia.Unmarshal(resp, mat); err != nil {
		log.Fatal(err)
	}

	// print elements of inverted matrix
	fmt.Println("inv mat:", mat.GetElms())
}
