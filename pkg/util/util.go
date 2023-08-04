package util

import (
	"fmt"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnake(input string) string {
	snake := matchFirstCap.ReplaceAllString(input, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func SnakeToCamel(input string) string {
	camel := ""
	isToUpper := false

	for k, v := range input {
		if k == 0 {
			camel = strings.ToUpper(string(input[0]))
		} else {
			if isToUpper {
				camel += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camel += string(v)
				}
			}
		}
	}

	return camel
}

func RemoveSnake(input string) string {
	out := strings.ReplaceAll(input, "_", "")
	out = strings.ToLower(out)
	fmt.Println(input, out)
	return out
}
