# julia
lib: Go interface to julia runtime intended for simple i/o between
go and julia runtimes

# disclaimer
> The use of this tool does not guarantee security or usability for any
> particular purpose. Please review the code and use at your own risk.

## installation
This step assumes you have [Go compiler toolchain](https://go.dev/dl/)
installed on your system with version at least Go 1.18.

Download this repo to a folder and cd to it.
```bash
go test
```

You will also need to have `julia` installed at `/usr/local/julia` which
is dynamically linked using `cgo`.

> Please note that this library is experimental and should not be used for
> production settings requiring multithreading and large volumes of data i/o
> to/from julia runtime.
> 
> It is meant for simple julia invocation from Go interface. Use with caution.

## usage
In the example below we create a random matrix in julia and fetch that in go.
Then pass that matrix back to julia runtime and compute its inverse.
```go
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
	fmt.Println("rand mat:", mat.Elms)

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
	fmt.Println("inv mat:", mat.Elms)
}
```

## dockerfile
Build provided examples as a container
```bash
docker build -t go-julia ./
```

Run an example in the container:
```bash
docker run --rm -it go-julia /go-julia/matrix-inversion
rand mat: [1.017768642481323 0.18851208555752758 -0.3488466841230449 0.4406143439802128 -0.9871042912008214 -2.3784718350906355 -1.8453603181415057 -0.09313474878433725 1.3359835582412696 -0.2341551927011799 0.7564688157590539 0.6625939970160911 1.7330849839849714 -1.0073081834684647 0.9176323073548744 -0.19748061907540335 1.0518872982633893 1.3556630206573586 -0.41537427531113674 1.8618622495593484 -0.08692220378171413 0.15844308568914403 1.961931713197817 -0.9674542403020966 0.6027724685356741]
inv mat: [0.017579422927962773 0.33640304238850627 2.0600006052967696 -0.4904229824862231 -1.4617453088048453 0.7276869478620276 -0.6639410803972795 -2.6019577734205965 1.0419842754172768 1.676332344686367 0.6160516767647928 0.20485388151853834 -0.06754593530676123 0.20961519218300062 0.5437915468222431 0.9668365082102905 0.510424584895026 0.3884752124354898 0.557668509045042 -0.5323606068080747 -0.642116529077303 0.3754996808917087 1.8243608863358483 -0.13181720095132604 -1.6168253384478684]
```
