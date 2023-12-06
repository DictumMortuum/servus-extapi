package db

import (
	"fmt"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
)

func YearConstraint(req *model.Map, start string) string {
	q := ""

	yearFlag, err := req.GetBool("year_flag")
	if err == nil && yearFlag {
		year, err := req.GetInt64("year")
		if err == nil && year != 0 {
			q = fmt.Sprintf("%s date >= '%d-01-01' and date < '%d-01-01'", start, year, year+1)
		}
	}

	return q
}

func YearHavingConstraint(req *model.Map) string {
	q := ""

	yearFlag, err := req.GetBool("year_flag")
	if err == nil && yearFlag {
		year, err := req.GetInt64("year")
		if err == nil && year != 0 {
			q = fmt.Sprintf("MAX(date) < '%d-01-01'", year)
		}
	}

	return q
}

func Limit(req *model.Map, n int64) string {
	if n == -1 {
		return ""
	} else {
		return fmt.Sprintf("limit %d", n)
	}
}
