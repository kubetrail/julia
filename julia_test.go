package julia

import (
	"fmt"
	"testing"
)

func TestNewMatInstantiation(t *testing.T) {
	if _, err := NewMat([]byte{1, 2, 3, 4}, 4); err != nil {
		t.Fatal(err)
	}

	if _, err := NewMat([]byte{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}

	if _, err := NewMat([]byte{1, 2, 3, 4}, 2, 2); err != nil {
		t.Fatal(err)
	}

	if _, err := NewMat([]byte{1, 2, 3, 4}, 2, 3); err == nil {
		t.Fatal("should have failed for dim 2x3")
	}
}

func TestMarshalPrimitiveTypes(t *testing.T) {
	Initialize()
	defer Finalize()

	if _, err := Marshal(true); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal(uint8(0)); err != nil {
		t.Fatal(err)
	}
	if _, err := Marshal(uint16(0)); err != nil {
		t.Fatal(err)
	}
	if _, err := Marshal(uint32(0)); err != nil {
		t.Fatal(err)
	}
	if _, err := Marshal(uint64(0)); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal(int8(0)); err != nil {
		t.Fatal(err)
	}
	if _, err := Marshal(int16(0)); err != nil {
		t.Fatal(err)
	}
	if _, err := Marshal(int32(0)); err != nil {
		t.Fatal(err)
	}
	if _, err := Marshal(int64(0)); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal(float32(0)); err != nil {
		t.Fatal(err)
	}
	if _, err := Marshal(float64(0)); err != nil {
		t.Fatal(err)
	}
}

func TestMarshalSlices(t *testing.T) {
	Initialize()
	defer Finalize()
	if _, err := Marshal([]bool{true, true, false, true}); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal([]uint8{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal([]uint16{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal([]uint32{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal([]uint64{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal([]int8{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal([]int16{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal([]int32{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal([]int64{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal([]float32{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}

	if _, err := Marshal([]float64{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}
}

func TestMarshalMultiDimensional(t *testing.T) {
	Initialize()
	defer Finalize()

	if mat, err := NewMat([]bool{true, true, false, true}, 2, 2); err != nil {
		t.Fatal(err)
	} else {
		if _, err := Marshal(mat); err != nil {
			t.Fatal(err)
		}
	}

	if mat, err := NewMat([]uint8{1, 2, 3, 4}, 2, 2); err != nil {
		t.Fatal(err)
	} else {
		if _, err := Marshal(mat); err != nil {
			t.Fatal(err)
		}
	}

	if mat, err := NewMat([]uint16{1, 2, 3, 4}, 2, 2); err != nil {
		t.Fatal(err)
	} else {
		if _, err := Marshal(mat); err != nil {
			t.Fatal(err)
		}
	}

	if mat, err := NewMat([]uint32{1, 2, 3, 4}, 2, 2); err != nil {
		t.Fatal(err)
	} else {
		if _, err := Marshal(mat); err != nil {
			t.Fatal(err)
		}
	}

	if mat, err := NewMat([]uint64{1, 2, 3, 4}, 2, 2); err != nil {
		t.Fatal(err)
	} else {
		if _, err := Marshal(mat); err != nil {
			t.Fatal(err)
		}
	}

	if mat, err := NewMat([]int8{1, 2, 3, 4}, 2, 2); err != nil {
		t.Fatal(err)
	} else {
		if _, err := Marshal(mat); err != nil {
			t.Fatal(err)
		}
	}

	if mat, err := NewMat([]int16{1, 2, 3, 4}, 2, 2); err != nil {
		t.Fatal(err)
	} else {
		if _, err := Marshal(mat); err != nil {
			t.Fatal(err)
		}
	}

	if mat, err := NewMat([]int32{1, 2, 3, 4}, 2, 2); err != nil {
		t.Fatal(err)
	} else {
		if _, err := Marshal(mat); err != nil {
			t.Fatal(err)
		}
	}

	if mat, err := NewMat([]int64{1, 2, 3, 4}, 2, 2); err != nil {
		t.Fatal(err)
	} else {
		if _, err := Marshal(mat); err != nil {
			t.Fatal(err)
		}
	}

	if mat, err := NewMat([]float32{1, 2, 3, 4}, 2, 2); err != nil {
		t.Fatal(err)
	} else {
		if _, err := Marshal(mat); err != nil {
			t.Fatal(err)
		}
	}

	if mat, err := NewMat([]float64{1, 2, 3, 4}, 2, 2); err != nil {
		t.Fatal(err)
	} else {
		if _, err := Marshal(mat); err != nil {
			t.Fatal(err)
		}
	}
}

func TestEvalFuncPrintlnBase(t *testing.T) {
	Initialize()
	defer Finalize()

	mat, err := NewMat([]float64{1, 2, 3, 4}, 2, 2)
	if err != nil {
		t.Fatal(err)
	}

	data, err := Marshal(mat)
	if err != nil {
		t.Fatal(err)
	}

	data, err = EvalFunc("println", ModuleBase, data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEvalFuncRandn(t *testing.T) {
	Initialize()
	defer Finalize()

	mat, err := NewMat([]int8{2, 2})
	if err != nil {
		t.Fatal(err)
	}

	data, err := Marshal(mat)
	if err != nil {
		t.Fatal(err)
	}

	data, err = EvalFunc("randn", ModuleBase, data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEval(t *testing.T) {
	Initialize()
	defer Finalize()

	if _, err := Eval("f(x::Vector{Int64}) = println(randn(x...))"); err != nil {
		t.Fatal(err)
	}

	mat, err := NewMat([]int64{2, 3})
	if err != nil {
		t.Fatal(err)
	}

	data, err := Marshal(mat)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = EvalFunc("f", ModuleMain, data); err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalOutputOfRandn2x3(t *testing.T) {
	Initialize()
	defer Finalize()

	// testing randn(2,3)
	n, m := 2, 3

	arg1, err := Marshal(int64(n))
	if err != nil {
		t.Fatal(err)
	}

	arg2, err := Marshal(int64(m))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := EvalFunc("randn", ModuleBase, arg1, arg2)
	if err != nil {
		t.Fatal(err)
	}

	out, err := NewMat(make([]float64, n*m), n, m)
	if err != nil {
		t.Fatal(err)
	}

	if err := Unmarshal(resp, out); err != nil {
		t.Fatal(err)
	}

	fmt.Println(out.Elms)
}

func TestTypeofSlice(t *testing.T) {
	Initialize()
	defer Finalize()

	{
		arg, err := Marshal([]bool{true, true, false, true})
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Vector{Int8}" {
			t.Fatal("expected Vector{Int8}, got", argType)
		}
	}

	{
		arg, err := Marshal([]uint8{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Vector{UInt8}" {
			t.Fatal("expected Vector{UInt8}, got", argType)
		}
	}

	{
		arg, err := Marshal([]uint16{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Vector{UInt16}" {
			t.Fatal("expected Vector{UInt16}, got", argType)
		}
	}

	{
		arg, err := Marshal([]uint32{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Vector{UInt32}" {
			t.Fatal("expected Vector{UInt32}, got", argType)
		}
	}

	{
		arg, err := Marshal([]uint64{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Vector{UInt64}" {
			t.Fatal("expected Vector{UInt64}, got", argType)
		}
	}

	{
		arg, err := Marshal([]int8{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Vector{Int8}" {
			t.Fatal("expected Vector{Int8}, got", argType)
		}
	}

	{
		arg, err := Marshal([]int16{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Vector{Int16}" {
			t.Fatal("expected Vector{Int16}, got", argType)
		}
	}

	{
		arg, err := Marshal([]int32{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Vector{Int32}" {
			t.Fatal("expected Vector{Int32}, got", argType)
		}
	}

	{
		arg, err := Marshal([]int64{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Vector{Int64}" {
			t.Fatal("expected Vector{Int64}, got", argType)
		}
	}

	{
		arg, err := Marshal([]float32{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Vector{Float32}" {
			t.Fatal("expected Vector{Float32}, got", argType)
		}
	}

	{
		arg, err := Marshal([]float64{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Vector{Float64}" {
			t.Fatal("expected Vector{Float64}, got", argType)
		}
	}
}

func TestTypeofMatrix(t *testing.T) {
	Initialize()
	defer Finalize()

	{
		x, err := NewMat([]bool{true, true, false, true}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Matrix{Int8}" {
			t.Fatal("expected Matrix{Int8}, got", argType)
		}
	}

	{
		x, err := NewMat([]uint8{1, 2, 3, 4}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Matrix{UInt8}" {
			t.Fatal("expected Matrix{UInt8}, got", argType)
		}
	}

	{
		x, err := NewMat([]uint16{1, 2, 3, 4}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Matrix{UInt16}" {
			t.Fatal("expected Matrix{UInt16}, got", argType)
		}
	}

	{
		x, err := NewMat([]uint32{1, 2, 3, 4}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Matrix{UInt32}" {
			t.Fatal("expected Matrix{UInt32}, got", argType)
		}
	}

	{
		x, err := NewMat([]uint64{1, 2, 3, 4}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Matrix{UInt64}" {
			t.Fatal("expected Matrix{UInt64}, got", argType)
		}
	}

	{
		x, err := NewMat([]int8{1, 2, 3, 4}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Matrix{Int8}" {
			t.Fatal("expected Matrix{Int8}, got", argType)
		}
	}

	{
		x, err := NewMat([]int16{1, 2, 3, 4}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Matrix{Int16}" {
			t.Fatal("expected Matrix{Int16}, got", argType)
		}
	}

	{
		x, err := NewMat([]int32{1, 2, 3, 4}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Matrix{Int32}" {
			t.Fatal("expected Matrix{Int32}, got", argType)
		}
	}

	{
		x, err := NewMat([]int64{1, 2, 3, 4}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Matrix{Int64}" {
			t.Fatal("expected Matrix{Int64}, got", argType)
		}
	}

	{
		x, err := NewMat([]float32{1, 2, 3, 4}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Matrix{Float32}" {
			t.Fatal("expected Matrix{Float32}, got", argType)
		}
	}

	{
		x, err := NewMat([]float64{1, 2, 3, 4}, 2, 2)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Matrix{Float64}" {
			t.Fatal("expected Matrix{Float64}, got", argType)
		}
	}
}

func TestTypeofTensor(t *testing.T) {
	Initialize()
	defer Finalize()

	{
		x, err := NewMat(make([]bool, 2*3*4), 2, 3, 4)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Array{Int8, 3}" {
			t.Fatal("expected Array{Int8, 3}, got", argType)
		}
	}

	{
		x, err := NewMat(make([]uint8, 2*3*4), 2, 3, 4)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Array{UInt8, 3}" {
			t.Fatal("expected Array{UInt8, 3}, got", argType)
		}
	}

	{
		x, err := NewMat(make([]uint16, 2*3*4), 2, 3, 4)
		if err != nil {
			t.Fatal(err)
		}

		arg, err := Marshal(x)
		if err != nil {
			t.Fatal(err)
		}

		argType := arg.Type()

		if argType != "Array{UInt16, 3}" {
			t.Fatal("expected Array{UInt16, 3}, got", argType)
		}
	}
}
