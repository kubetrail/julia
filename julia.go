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

func Initialize() {
	/* required: setup the Julia context */
	C.jl_init()
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
func marshal(x interface{}) (*jlValue, error) {
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
		n := uint64(len(v.Dims))
		arrayType := C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_int8_type,
				),
			),
			C.ulong(n),
		)

		array, err := allocArray(arrayType, v.Dims...)
		if err != nil {
			return nil, fmt.Errorf("could not allocate array: %w", err)
		}

		data := array.data
		ptr := unsafe.Pointer(data)

		el := int8(0)
		count, err := dim2NumElms(v.Dims)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			p := (*int8)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
			if v.Elms[i] {
				*p = int8(1)
			} else {
				*p = int8(0)
			}
		}

		return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
	case *Mat[uint8]:
		n := uint64(len(v.Dims))
		arrayType := C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_uint8_type,
				),
			),
			C.ulong(n),
		)

		array, err := allocArray(arrayType, v.Dims...)
		if err != nil {
			return nil, fmt.Errorf("could not allocate array: %w", err)
		}

		data := array.data
		ptr := unsafe.Pointer(data)

		// https://stackoverflow.com/questions/49931051/cgo-how-do-you-use-pointers-in-golang-to-access-data-from-an-array-in-c
		el := byte(0)
		count, err := dim2NumElms(v.Dims)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			// https://stackoverflow.com/a/49961256
			p := (*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
			*p = v.Elms[i]
		}

		return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
	case *Mat[uint16]:
		n := uint64(len(v.Dims))
		arrayType := C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_uint16_type,
				),
			),
			C.ulong(n),
		)

		array, err := allocArray(arrayType, v.Dims...)
		if err != nil {
			return nil, fmt.Errorf("could not allocate array: %w", err)
		}

		data := array.data
		ptr := unsafe.Pointer(data)

		el := uint16(0)
		count, err := dim2NumElms(v.Dims)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			p := (*uint16)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
			*p = v.Elms[i]
		}

		return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
	case *Mat[uint32]:
		n := uint64(len(v.Dims))
		arrayType := C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_uint32_type,
				),
			),
			C.ulong(n),
		)

		array, err := allocArray(arrayType, v.Dims...)
		if err != nil {
			return nil, fmt.Errorf("could not allocate array: %w", err)
		}

		data := array.data
		ptr := unsafe.Pointer(data)

		el := uint32(0)
		count, err := dim2NumElms(v.Dims)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			p := (*uint32)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
			*p = v.Elms[i]
		}

		return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
	case *Mat[uint64]:
		n := uint64(len(v.Dims))
		arrayType := C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_uint64_type,
				),
			),
			C.ulong(n),
		)

		array, err := allocArray(arrayType, v.Dims...)
		if err != nil {
			return nil, fmt.Errorf("could not allocate array: %w", err)
		}

		data := array.data
		ptr := unsafe.Pointer(data)

		el := uint64(0)
		count, err := dim2NumElms(v.Dims)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			p := (*uint64)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
			*p = v.Elms[i]
		}

		return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
	case *Mat[int8]:
		n := uint64(len(v.Dims))
		arrayType := C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_int8_type,
				),
			),
			C.ulong(n),
		)

		array, err := allocArray(arrayType, v.Dims...)
		if err != nil {
			return nil, fmt.Errorf("could not allocate array: %w", err)
		}

		data := array.data
		ptr := unsafe.Pointer(data)

		el := int8(0)
		count, err := dim2NumElms(v.Dims)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			p := (*int8)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
			*p = v.Elms[i]
		}

		return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
	case *Mat[int16]:
		n := uint64(len(v.Dims))
		arrayType := C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_int16_type,
				),
			),
			C.ulong(n),
		)

		array, err := allocArray(arrayType, v.Dims...)
		if err != nil {
			return nil, fmt.Errorf("could not allocate array: %w", err)
		}

		data := array.data
		ptr := unsafe.Pointer(data)

		el := int16(0)
		count, err := dim2NumElms(v.Dims)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			p := (*int16)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
			*p = v.Elms[i]
		}

		return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
	case *Mat[int32]:
		n := uint64(len(v.Dims))
		arrayType := C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_int32_type,
				),
			),
			C.ulong(n),
		)

		array, err := allocArray(arrayType, v.Dims...)
		if err != nil {
			return nil, fmt.Errorf("could not allocate array: %w", err)
		}

		data := array.data
		ptr := unsafe.Pointer(data)

		el := int32(0)
		count, err := dim2NumElms(v.Dims)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			p := (*int32)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
			*p = v.Elms[i]
		}

		return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
	case *Mat[int64]:
		n := uint64(len(v.Dims))
		arrayType := C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_int64_type,
				),
			),
			C.ulong(n),
		)

		array, err := allocArray(arrayType, v.Dims...)
		if err != nil {
			return nil, fmt.Errorf("could not allocate array: %w", err)
		}

		data := array.data
		ptr := unsafe.Pointer(data)

		el := int64(0)
		count, err := dim2NumElms(v.Dims)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			p := (*int64)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
			*p = v.Elms[i]
		}

		return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
	case *Mat[float32]:
		n := uint64(len(v.Dims))
		arrayType := C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_float32_type,
				),
			),
			C.ulong(n),
		)

		array, err := allocArray(arrayType, v.Dims...)
		if err != nil {
			return nil, fmt.Errorf("could not allocate array: %w", err)
		}

		data := array.data
		ptr := unsafe.Pointer(data)

		el := float32(0)
		count, err := dim2NumElms(v.Dims)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			p := (*float32)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
			*p = v.Elms[i]
		}

		return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
	case *Mat[float64]:
		n := uint64(len(v.Dims))
		arrayType := C.jl_apply_array_type(
			(*(C.jl_value_t))(
				unsafe.Pointer(
					C.jl_float64_type,
				),
			),
			C.ulong(n),
		)

		array, err := allocArray(arrayType, v.Dims...)
		if err != nil {
			return nil, fmt.Errorf("could not allocate array: %w", err)
		}

		data := array.data
		ptr := unsafe.Pointer(data)

		el := float64(0)
		count, err := dim2NumElms(v.Dims)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			p := (*float64)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(el)))
			*p = v.Elms[i]
		}

		return &jlValue{value: (*(C.jl_value_t))(unsafe.Pointer(array))}, nil
	default:
		return nil, fmt.Errorf("invalid type, not supported %T", v)
	}
}

// unmarshal unpacks generic jlValue and populates pointer value in x
func unmarshal(data *jlValue, x interface{}) error {
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
		// same logic for *[]bool as done above for *[]byte
		array := (*(C.jl_array_t))(unsafe.Pointer(value))
		data := array.data
		ptr := unsafe.Pointer(data)
		for i := range v.Elms {
			// https://stackoverflow.com/a/49961256
			p := (*bool)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(false)))
			(*v).Elms[i] = *p
		}
	case *Mat[uint8]:
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
			p := (*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(byte(0))))
			(*v).Elms[i] = *p
		}
	case *Mat[uint16]:
		array := (*(C.jl_array_t))(unsafe.Pointer(value))
		data := array.data
		ptr := unsafe.Pointer(data)
		for i := range v.Elms {
			// https://stackoverflow.com/a/49961256
			p := (*uint16)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(uint16(0))))
			(*v).Elms[i] = *p
		}
	case *Mat[uint32]:
		array := (*(C.jl_array_t))(unsafe.Pointer(value))
		data := array.data
		ptr := unsafe.Pointer(data)
		for i := range v.Elms {
			// https://stackoverflow.com/a/49961256
			p := (*uint32)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(uint32(0))))
			(*v).Elms[i] = *p
		}
	case *Mat[uint64]:
		array := (*(C.jl_array_t))(unsafe.Pointer(value))
		data := array.data
		ptr := unsafe.Pointer(data)
		for i := range v.Elms {
			// https://stackoverflow.com/a/49961256
			p := (*uint64)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(uint64(0))))
			(*v).Elms[i] = *p
		}
	case *Mat[int8]:
		array := (*(C.jl_array_t))(unsafe.Pointer(value))
		data := array.data
		ptr := unsafe.Pointer(data)
		for i := range v.Elms {
			// https://stackoverflow.com/a/49961256
			p := (*int8)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(int8(0))))
			(*v).Elms[i] = *p
		}
	case *Mat[int16]:
		array := (*(C.jl_array_t))(unsafe.Pointer(value))
		data := array.data
		ptr := unsafe.Pointer(data)
		for i := range v.Elms {
			// https://stackoverflow.com/a/49961256
			p := (*int16)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(int16(0))))
			(*v).Elms[i] = *p
		}
	case *Mat[int32]:
		array := (*(C.jl_array_t))(unsafe.Pointer(value))
		data := array.data
		ptr := unsafe.Pointer(data)
		for i := range v.Elms {
			// https://stackoverflow.com/a/49961256
			p := (*int32)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(int32(0))))
			(*v).Elms[i] = *p
		}
	case *Mat[int64]:
		array := (*(C.jl_array_t))(unsafe.Pointer(value))
		data := array.data
		ptr := unsafe.Pointer(data)
		for i := range v.Elms {
			// https://stackoverflow.com/a/49961256
			p := (*int64)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(int64(0))))
			(*v).Elms[i] = *p
		}
	case *Mat[float32]:
		// same logic for *[]float32 as done above for *[]byte
		array := (*(C.jl_array_t))(unsafe.Pointer(value))
		data := array.data
		ptr := unsafe.Pointer(data)
		for i := range v.Elms {
			// https://stackoverflow.com/a/49961256
			p := (*float32)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(float32(0))))
			(*v).Elms[i] = *p
		}
	case *Mat[float64]:
		// same logic for *[]float64 as done above for *[]byte
		array := (*(C.jl_array_t))(unsafe.Pointer(value))
		data := array.data
		ptr := unsafe.Pointer(data)
		for i := range v.Elms {
			// https://stackoverflow.com/a/49961256
			p := (*float64)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof(float64(0))))
			(*v).Elms[i] = *p
		}
	default:
		return fmt.Errorf("invalid type, not supported %T", v)
	}

	return nil
}
