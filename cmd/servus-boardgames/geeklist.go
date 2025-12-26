package main

import (
	"github.com/DictumMortuum/servus-extapi/pkg/bgg"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
)

func GetGeeklist(req *model.Map, res *model.Map) error {

	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	rs, err := bgg.GetGeeklist(id)
	if err != nil {
		return err
	}

	res.Set("options", rs)

	return nil
}
