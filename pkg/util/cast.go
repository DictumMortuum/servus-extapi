package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ToArray[T any](col any, t T) ([]T, error) {
	switch v := col.(type) {
	case []T:
		return col.([]T), nil
	default:
		return nil, fmt.Errorf("%+v not a valid type array", v)
	}
}

func Atoi(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func Atoi64(s string) int64 {
	i, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	return i
}

func Atof(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
