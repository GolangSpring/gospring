package helper

import (
	"errors"
	"fmt"
)

// Performs a destructuring on slice/array values using generics
func SliceUnpack[T any](src []T, dests ...*T) error {
	if len(dests) == 0 {
		return errors.New("unpack: destination can't be empty")
	}

	if len(src) < len(dests) {
		return errors.New("unpack: not enough source values")
	}

	for idx, dest := range dests {
		if dest == nil {
			return fmt.Errorf("unpack: destination at index %d is nil", idx)
		}
		*dests[idx] = src[idx]
	}
	return nil
}
