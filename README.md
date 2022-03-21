# julia
This is a `go` interface to `julia` runtime intended for single threaded
and lightweight use for simple code execution.

Please start with [embedding](https://docs.julialang.org/en/v1/manual/embedding/)
docs, which is the basic design on which this library is built.

# disclaimer
> The use of this tool does not guarantee security or usability for any
> particular purpose. Please review the code and use at your own risk.
> 
> Please also see known limitations at the end of this doc.

## installation
This step assumes you have [Go compiler toolchain](https://go.dev/dl/)
installed on your system with version at least Go 1.18. The library
makes use of `go` generics.

You will also need to have `julia` installed at `/usr/local/julia` which
is dynamically linked using `cgo`.

> Please note that this library is experimental and should not be used for
> production settings requiring multithreading and large volumes of data i/o
> to/from julia runtime.
> 
> It is meant for simple julia invocation from Go interface. Use with caution.

Download this repo to a folder and cd to it.
```bash
go test
```

## usage
In the example below we create a random matrix in `julia` and fetch that in `go`.
Then pass that matrix back to `julia` runtime and compute its inverse.
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

Here is another example of matrix multiplication:
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
```
## dockerfile
A `dockerfile` can be used to build these examples
```bash
docker build -t go-julia ./
```

Run an example in the container:
```bash
docker run --rm -it go-julia /go-julia/matrix-inversion
rand mat: [1.017768642481323 0.18851208555752758 -0.3488466841230449 0.4406143439802128 -0.9871042912008214 -2.3784718350906355 -1.8453603181415057 -0.09313474878433725 1.3359835582412696 -0.2341551927011799 0.7564688157590539 0.6625939970160911 1.7330849839849714 -1.0073081834684647 0.9176323073548744 -0.19748061907540335 1.0518872982633893 1.3556630206573586 -0.41537427531113674 1.8618622495593484 -0.08692220378171413 0.15844308568914403 1.961931713197817 -0.9674542403020966 0.6027724685356741]
inv mat: [0.017579422927962773 0.33640304238850627 2.0600006052967696 -0.4904229824862231 -1.4617453088048453 0.7276869478620276 -0.6639410803972795 -2.6019577734205965 1.0419842754172768 1.676332344686367 0.6160516767647928 0.20485388151853834 -0.06754593530676123 0.20961519218300062 0.5437915468222431 0.9668365082102905 0.510424584895026 0.3884752124354898 0.557668509045042 -0.5323606068080747 -0.642116529077303 0.3754996808917087 1.8243608863358483 -0.13181720095132604 -1.6168253384478684]
```

```bash
docker run --rm -it go-julia /go-julia/matrix-multiplication
matrix multiplied by its inverse is an identity matrix:
[0.9999999999999994 2.393094211668032e-16 -1.0471181754244775e-15 7.021075814654354e-16 3.076390399372395e-16 8.757211127093011e-17 0.9999999999999989 2.7294181293716953e-16 3.5320707519199396e-16 -3.2696702378511e-16 1.7145691212372143e-16 5.635048558235537e-17 1 -1.706945437414341e-16 -2.3262142907258317e-17 9.85627416443657e-17 -3.884989330059967e-16 2.2894671093460267e-16 0.9999999999999999 -3.7378253110592907e-16 -1.1418060527744768e-16 9.42406854610802e-16 -2.2780738535360076e-16 -8.110240689058509e-16 1.0000000000000002]
```

## architecture
The library is not intended to cover all possible `julia` execution scenarios and meant for 
simple use cases where `go` front end is required on top of `julia` backend runtime.

A type-parametrized `matrix` type is defined:
```go
// Mat represents the matrix for supported data types
// parameterized by primitive types
type Mat[T PrimitiveTypes] struct {
	Elms []T   `json:"elms,omitempty"`
	Dims []int `json:"dims,omitempty"`
}
```

This is the main data structure to send and receive values between `go` and `julia` runtimes.

Furthermore, `Marshal` and `Unmarshal` functions are defined to work with `Mat` data
structure to pack/unpack data into a `julia` native generic data type.

## known issues
Foreign function interface to `julia` via its `C-API` should be used with
caution and preferably run in a single threaded mode. Considering `go` allows
concurrency very easily, extra care needs to be taken to ensure that `julia`
does not run as part of multiple goroutines.

* https://discourse.julialang.org/t/problems-scaling-jl-alloc-array-2d-c-api/63341
