package util

import (
	"fmt"
)

func ToArray[T any](col any, t T) ([]T, error) {
	switch v := col.(type) {
	case []T:
		return col.([]T), nil
	default:
		return nil, fmt.Errorf("%+v not a valid type array", v)
	}
}
