package julia

import (
	"testing"
)

// data struct is type-parametrized
type data[T PrimitiveTypes] struct {
	values []T
	name   string
}

// Len is a method on data parametrized by T
func (g *data[T]) Len() int {
	return len(g.values)
}

func (g *data[T]) IsInt32() bool {
	var el T
	switch any(el).(type) {
	case int32:
		return true
	default:
		return false
	}
}

// IsSameType checks sameness of type parameter inputs by
// using a type switch on a new instance of a variable.
// type switches cannot be used on variable types constrained
// by type parameters, hence it needs to be wrapped in an
// interface such as any
func IsSameType[T, W PrimitiveTypes]() bool {
	var el T
	switch any(el).(type) {
	case W:
		return true
	default:
		return false
	}
}

func newData[T PrimitiveTypes](x []T, name string) *data[T] {
	return &data[T]{
		values: x,
		name:   name,
	}
}

func TestLen(t *testing.T) {
	d := newData([]int32{1, 2, 3}, "data")
	if d.Len() != 3 {
		t.Fatal("expected length to be 3 using method")
	}
}

func TestIsSameType(t *testing.T) {
	if IsSameType[int32, float64]() {
		t.Fatal("expected types to be different")
	}

	if !IsSameType[int32, int32]() {
		t.Fatal("expected types to be the same")
	}
}
