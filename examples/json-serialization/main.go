package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kubetrail/julia"
)

const (
	// run julia code to load required packages
	// and declare a few functions.
	// these functions will be invoked from go code
	jCode = `
# load JSON2 package for data serialization
using JSON2;

# unmarshal takes byte buffer and unmarshals it into a vector of strings 
unmarshal(x::Vector{UInt8}) = JSON2.read(String(x), Vector{String});

# marshal takes any x and serializes it to a byte buffer
marshal(x) = Vector{UInt8}(JSON2.write(x));
`
)

func main() {
	julia.Initialize()
	defer julia.Finalize()

	// load julia code to marshal and unmarshal list of strings
	if _, err := julia.Eval(jCode); err != nil {
		log.Fatal(err)
	}

	// input is a list of strings to be passed to julia
	// as a json serialized byte buffer
	input := []string{"abcd", "123456", "%^&%##*"}

	// convert input arguments to byte buffer
	jb, err := json.Marshal(input)
	if err != nil {
		log.Fatal(err)
	}

	// convert byte buffer to julia argument
	arg, err := julia.Marshal(jb)
	if err != nil {
		log.Fatal(err)
	}

	// pass argument to unarshal function
	// unmarshal runs in julia and delivers a Vector{String}
	resp, err := julia.EvalFunc("unmarshal", julia.ModuleMain, arg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("this is an example to show how to pass list of strings to julia using json serialization")
	fmt.Println(input, "passed to julia as list of strings")
	fmt.Println("expected data type in julia: Vector{String}")
	fmt.Println("received data type in julia:", resp.Type())
}
