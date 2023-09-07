package util

import "sort"

func Unique(col []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range col {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	sort.Strings(list)
	return list
}
