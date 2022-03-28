package julia

import "fmt"

// PrimitiveTypes are type constraint on julia input
type PrimitiveTypes interface {
	~bool |
		~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

// PrimitiveSliceTypes are type constraints on julia inputs
type PrimitiveSliceTypes interface {
	~[]bool |
		~[]uint8 | ~[]uint16 | ~[]uint32 | ~[]uint64 |
		~[]int8 | ~[]int16 | ~[]int32 | ~[]int64 |
		~[]float32 | ~[]float64
}

// PrimitivePointerTypes are type constraints on julia output
type PrimitivePointerTypes interface {
	~*bool |
		~*uint8 | ~*uint16 | ~*uint32 | ~*uint64 |
		~*int8 | ~*int16 | ~*int32 | ~*int64 |
		~*float32 | ~*float64
}

// MatTypes represents constraints on parametrized Mat type
// to be uses as both julia inputs and outputs
type MatTypes interface {
	*Mat[bool] |
		*Mat[uint8] | *Mat[uint16] | *Mat[uint32] | *Mat[uint64] |
		*Mat[int8] | *Mat[int16] | *Mat[int32] | *Mat[int64] |
		*Mat[float32] | *Mat[float64]
}

// Mat represents the matrix for supported data types
// parameterized by primitive types
type Mat[T PrimitiveTypes] struct {
	elms []T   `json:"elms,omitempty"`
	dims []int `json:"dims,omitempty"`
}

func (g *Mat[T]) GetElms() []T {
	return g.elms
}

func (g *Mat[T]) GetDims() []int {
	return g.dims
}

// NewMat creates a new instance of matrix and validates if the length of
// elements is satisfied by the dimensions
func NewMat[T PrimitiveTypes](values []T, dims ...int) (*Mat[T], error) {
	if len(dims) == 0 {
		dims = []int{len(values)}
	}

	m := &Mat[T]{
		dims: dims,
		elms: values,
	}

	if len(m.dims) == 0 {
		return nil, fmt.Errorf("invalid dimensions")
	}

	numElements, err := dim2NumElms(m.dims)
	if err != nil {
		return nil, err
	}

	if numElements != len(values) {
		return nil, fmt.Errorf("dims and len elms mismatch")
	}

	return m, nil
}
