package julia

import "fmt"

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
