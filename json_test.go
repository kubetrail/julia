package julia

import (
	"encoding/json"
	"testing"
)

const (
	preload = `
using JSON2
parse(x) = JSON2.read(String(x), Vector{String})
serialize(x) = Vector{UInt8}(JSON2.write(x))
`
)

// TestJsonSerialization sends a list of strings to julia
// as json serialized byte buffer which is then read back
// as Vector{String} in julia. Needless to say this method
// should work for arbitrary data types including structs,
// maps, slices etc.
func TestJsonSerializationSendToJulia(t *testing.T) {
	Initialize()
	defer Finalize()

	if _, err := Eval(preload); err != nil {
		t.Fatal(err)
	}

	listOfStrings := []string{"abcd", "12345678"}
	jb, _ := json.Marshal(listOfStrings)

	arg, err := Marshal(jb)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := EvalFunc("parse", ModuleMain, arg)
	if err != nil {
		t.Fatal(err)
	}

	respType := resp.Type()
	if respType != "Vector{String}" {
		t.Fatal("expected Vector{String}, got", respType)
	}
}

func TestJsonSerializationReceiveFromJulia(t *testing.T) {
	Initialize()
	defer Finalize()

	if _, err := Eval(preload); err != nil {
		t.Fatal(err)
	}

	resp, err := Eval("serialize([\"abcd\", \"12345679\"])")
	if err != nil {
		t.Fatal(err)
	}

	n := Len(resp)
	mat, err := NewMat(make([]byte, n))
	if err != nil {
		t.Fatal(err)
	}

	if err := Unmarshal(resp, mat); err != nil {
		t.Fatal(err)
	}

	var list []string
	if err := json.Unmarshal(mat.GetElms(), &list); err != nil {
		t.Fatal(err)
	}

	if list[0] != "abcd" ||
		list[1] != "12345679" {
		t.Fatal("did not receive expected values")
	}
}
