package julia

/*
// Start with the basic example from https://docs.julialang.org/en/v1/manual/embedding/
//
// Obviously the paths below may need to be modified to match your julia install location and version number.
//
#cgo CFLAGS: -fPIC -DJULIA_INIT_DIR="/usr/local/julia/lib" -I/usr/local/julia/include/julia -I.
#cgo LDFLAGS: -L/usr/local/julia/lib/julia  -L/usr/local/julia/lib -Wl,-rpath,/usr/local/julia/lib -ljulia
#include <julia.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type ModuleType int

const (
	ModuleBase ModuleType = iota
	ModuleMain
)

const (
	jlValueTypeOf = "__jlValueTypeOf"
)

func Initialize() {
	/* required: setup the Julia context */
	C.jl_init()

	// declare a few functions for use in this library
	_, _ = Eval(fmt.Sprintf("%s(x) = Vector{UInt8}(string(typeof(x)))", jlValueTypeOf))
}

func Finalize() {
	/* strongly recommended: notify Julia that the
	   program is about to terminate. this allows
	   Julia time to cleanup pending write requests
	   and run all finalizers
	*/
	C.jl_atexit_hook(0)
}

// jlValue represents generic data type to pass to/from Julia runtime.
// Users are not expected to instantiate this struct and an instance
// of it is typically accessed via Marshal/Unmarshal functions.
type jlValue struct {
	value *C.jl_value_t
}

// Type evaluates to julia representation of typeof
func (g *jlValue) Type() string {
	resp, _ := EvalFunc(jlValueTypeOf, ModuleMain, g)
	n := Len(resp)

	out, _ := NewMat(make([]uint8, n))
	_ = Unmarshal(resp, out)

	return string(out.GetElms())
}

func Len(g *jlValue) int {
	length, _ := EvalFunc("length", ModuleBase, g)

	var n int64
	_ = Unmarshal(length, &n)

	return int(n)
}

func Marshal[T PrimitiveTypes | PrimitiveSliceTypes | MatTypes](x T) (*jlValue, error) {
	return marshal(x)
}

func Unmarshal[T PrimitivePointerTypes | MatTypes](data *jlValue, x T) error {
	return unmarshal(data, x)
}

// Eval evaluates input as if it were julia code
func Eval(input string) (*jlValue, error) {
	return &jlValue{value: C.jl_eval_string(C.CString(input))}, nil
}

// EvalFunc evaluates a function literal, represented by name and module it is defined in,
// and passes any optional arguments to it, returning any output from julia runtime
func EvalFunc(name string, moduleType ModuleType, args ...*jlValue) (*jlValue, error) {
	var f *C.jl_function_t
	switch moduleType {
	case ModuleBase:
		f = C.jl_get_function(C.jl_base_module, C.CString(name))
	case ModuleMain:
		f = C.jl_get_function(C.jl_main_module, C.CString(name))
	}

	inputs := make([]*C.jl_value_t, len(args))
	for i, arg := range args {
		inputs[i] = arg.value
	}

	if len(args) > 0 {
		return &jlValue{value: C.jl_call(f, &(inputs[0]), C.int(len(inputs)))}, nil
	} else {
		return &jlValue{value: C.jl_call0(f)}, nil
	}
}

// https://discourse.julialang.org/t/problems-scaling-jl-alloc-array-2d-c-api/63341
func allocArray(arrayType *C.jl_value_t, dims ...int) (*C.jl_array_t, error) {
	var array *C.jl_array_t
	switch n := len(dims); n {
	case 0:
		return nil, fmt.Errorf("pl input at least one dimension")
	case 1:
		array = C.jl_alloc_array_1d(
			arrayType,
			C.ulong(dims[0]),
		)
	case 2:
		array = C.jl_alloc_array_2d(
			arrayType,
			C.ulong(uint64(dims[0])),
			C.ulong(uint64(dims[1])),
		)
	case 3:
		array = C.jl_alloc_array_3d(
			arrayType,
			C.ulong(dims[0]),
			C.ulong(dims[1]),
			C.ulong(dims[2]),
		)
	default:
		jdims := make([]C.ulong, n)
		for i := range dims {
			jdims[i] = C.ulong(uint64(dims[i]))
		}

		jdimPtr := (*(C.jl_value_t))(unsafe.Pointer(&jdims[0]))

		array = C.jl_new_array(arrayType, jdimPtr)
	}

	return array, nil
}

// marshal packs supported input to a generic jlValue type to pass to
// julia runtime
func marshal(x any) (*jlValue, error) {
	switch v := x.(type) {
	case bool:
		if v {
			return &jlValue{value: C.jl_box_bool(C.schar(int8(1)))}, nil
		} else {
			return &jlValue{value: C.jl_box_bool(C.schar(int8(0)))}, nil
		}
	case uint8:
		return &jlValue{value: C.jl_box_uint8(C.uchar(v))}, nil
	case uint16:
		return &jlValue{value: C.jl_box_uint16(C.ushort(v))}, nil
	case uint32:
		return &jlValue{value: C.jl_box_uint32(C.uint(v))}, nil
	case uint64:
		return &jlValue{value: C.jl_box_uint64(C.ulong(v))}, nil
	case int8:
		return &jlValue{value: C.jl_box_int8(C.schar(v))}, nil
	case int16:
		return &jlValue{value: C.jl_box_int16(C.short(v))}, nil
	case int32:
		return &jlValue{value: C.jl_box_int32(C.int(v))}, nil
	case int64:
		return &jlValue{value: C.jl_box_int64(C.long(v))}, nil
	case float32:
		return &jlValue{value: C.jl_box_float32(C.float(v))}, nil
	case float64:
		return &jlValue{value: C.jl_box_float64(C.double(v))}, nil
	case []bool:
		m, err := NewMat(v, len(v))
		if err != nil {
			return nil, err
		}
		return Marshal(m)
	case []uint8:
		m, err := NewMat(v, len(v))
		if err != nil {
			return nil, err
		}
		return Marshal(m)
	case []uint16:
		m, err := NewMat(v, len(v))
		if err != nil {
			return nil, err
		}
		return Marshal(m)
	case []uint32:
		m, err := NewMat(v, len(v))
		if err != nil {
			return nil, err
		}
		return Marshal(m)
	case []uint64:
		m, err := NewMat(v, len(v))
		if err != nil {
			return nil, err
		}
		return Marshal(m)
	case []float32:
		m, err := NewMat(v, len(v))
		if err != nil {
			return nil, err
		}
		return Marshal(m)
	case []float64:
		m, err := NewMat(v, len(v))
		if err != nil {
			return nil, err
		}
		return Marshal(m)
	case []int8:
		m, err := NewMat(v, len(v))
		if err != nil {
			return nil, err
		}
		return Marshal(m)
	case []int16:
		m, err := NewMat(v, len(v))
		if err != nil {
			return nil, err
		}
		return Marshal(m)
	case []int32:
		m, err := NewMat(v, len(v))
		if err != nil {
			return nil, err
		}
		return Marshal(m)
	case []int64:
		m, err := NewMat(v, len(v))
		if err != nil {
			return nil, err
		}
		return Marshal(m)
	case *Mat[bool]:
		return marshalMat[bool, *bool](v)
	case *Mat[uint8]:
		return marshalMat[uint8, *uint8](v)
	case *Mat[uint16]:
		return marshalMat[uint16, *uint16](v)
	case *Mat[uint32]:
		return marshalMat[uint32, *uint32](v)
	case *Mat[uint64]:
		return marshalMat[uint64, *uint64](v)
	case *Mat[int8]:
		return marshalMat[int8, *int8](v)
	case *Mat[int16]:
		return marshalMat[int16, *int16](v)
	case *Mat[int32]:
		return marshalMat[int32, *int32](v)
	case *Mat[int64]:
		return marshalMat[int64, *int64](v)
	case *Mat[float32]:
		return marshalMat[float32, *float32](v)
	case *Mat[float64]:
		return marshalMat[float64, *float64](v)
	default:
		return nil, fmt.Errorf("invalid type, not supported %T", v)
	}
}

// unmarshal unpacks generic jlValue and populates pointer value in x
func unmarshal(data *jlValue, x any) error {
	value := data.value
	switch v := x.(type) {
	case *bool:
		if C.jl_unbox_bool(value) == 1 {
			*v = true
		} else {
			*v = false
		}
	case *uint8:
		*v = uint8(C.jl_unbox_uint8(value))
	case *uint16:
		*v = uint16(C.jl_unbox_uint16(value))
	case *uint32:
		*v = uint32(C.jl_unbox_uint32(value))
	case *uint64:
		*v = uint64(C.jl_unbox_uint64(value))
	case *int8:
		*v = int8(C.jl_unbox_int8(value))
	case *int16:
		*v = int16(C.jl_unbox_int16(value))
	case *int32:
		*v = int32(C.jl_unbox_int32(value))
	case *int64:
		*v = int64(C.jl_unbox_int64(value))
	case *float32:
		*v = float32(C.jl_unbox_float32(value))
	case *float64:
		*v = float64(C.jl_unbox_float64(value))
	case *Mat[bool]:
		unmarshalMat[bool, *bool](data, v)
	case *Mat[uint8]:
		unmarshalMat[uint8, *uint8](data, v)
	case *Mat[uint16]:
		unmarshalMat[uint16, *uint16](data, v)
	case *Mat[uint32]:
		unmarshalMat[uint32, *uint32](data, v)
	case *Mat[uint64]:
		unmarshalMat[uint64, *uint64](data, v)
	case *Mat[int8]:
		unmarshalMat[int8, *int8](data, v)
	case *Mat[int16]:
		unmarshalMat[int16, *int16](data, v)
	case *Mat[int32]:
		unmarshalMat[int32, *int32](data, v)
	case *Mat[int64]:
		unmarshalMat[int64, *int64](data, v)
	case *Mat[float32]:
		unmarshalMat[float32, *float32](data, v)
	case *Mat[float64]:
		unmarshalMat[float64, *float64](data, v)
	default:
		return fmt.Errorf("invalid type, not supported %T", v)
	}

	return nil
}

// marshalMat is a generic serialization of input matrix to julia value.
// since type casting to pointer of T is required, it seems it is
// required to parametrize the pointer of T!
func marshalMat[T PrimitiveTypes, PtrT *T](v *Mat[T]) (*jlValue, error) {
	n := uint64(len(v.Dims))
	var el T

	arrayType, _ := getArrayType(n, el)
	array, err := allocArray(arrayType, v.Dims...)
	if err != nil {
		return nil, fmt.Errorf("could not allocate array: %w", err)
	}

	data := array.data
	ptr := unsafe.Pointer(data)

	for i := range v.Elms {
		p := (PtrT)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
		*p = v.Elms[i]
	}

	return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
}

// unmarshalMat is a generic way to unmarshal julia value into matrix type
// type-parametrized by primitive types. interestingly, we need to
// type-parametrize this function using both T and its pointer.
func unmarshalMat[T PrimitiveTypes, PtrT *T](jlValue *jlValue, v *Mat[T]) {
	var el T
	value := jlValue.value

	// cast value as unsafe pointer first, which makes it
	// equivalent to void* in C, then cast it to
	// pointer of jl_array_t
	array := (*(C.jl_array_t))(unsafe.Pointer(value))

	// access the data field
	data := array.data

	// cast it as unsafe pointer in order to perform
	// pointer arithmetics
	ptr := unsafe.Pointer(data)

	// better be sure that returned data is of that specific length
	for i := range v.Elms {
		// https://stackoverflow.com/a/49961256
		p := (PtrT)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
		(*v).Elms[i] = *p
	}
}

// getArrayType takes element el as empty interface type because we can't do
// type switch on generics!
func getArrayType(n uint64, el any) (*C.jl_value_t, error) {
	switch el.(type) {
	case bool:
		return C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_int8_type,
				),
			),
			C.ulong(n),
		), nil
	case uint8:
		return C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_uint8_type,
				),
			),
			C.ulong(n),
		), nil
	case uint16:
		return C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_uint16_type,
				),
			),
			C.ulong(n),
		), nil
	case uint32:
		return C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_uint32_type,
				),
			),
			C.ulong(n),
		), nil
	case uint64:
		return C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_uint64_type,
				),
			),
			C.ulong(n),
		), nil
	case int8:
		return C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_int8_type,
				),
			),
			C.ulong(n),
		), nil
	case int16:
		return C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_int16_type,
				),
			),
			C.ulong(n),
		), nil
	case int32:
		return C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_int32_type,
				),
			),
			C.ulong(n),
		), nil
	case int64:
		return C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_int64_type,
				),
			),
			C.ulong(n),
		), nil
	case float32:
		return C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_float32_type,
				),
			),
			C.ulong(n),
		), nil
	case float64:
		return C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_float64_type,
				),
			),
			C.ulong(n),
		), nil
	default:
		return nil, fmt.Errorf("invalid type, not supported %T", el)
	}
}

// dim2NumElms returns total number of elements inferred by dimension sizes
func dim2NumElms(dims []int) (int, error) {
	var numElements int
	if len(dims) == 0 {
		return 0, fmt.Errorf("invalid dims, length needs to be greater than 0")
	}

	for i, dim := range dims {
		if dim <= 0 {
			return 0, fmt.Errorf("invalid dims, needs to be greater than 0")
		}
		if i == 0 {
			numElements = dim
			continue
		}

		numElements *= dim
	}

	return numElements, nil
}
