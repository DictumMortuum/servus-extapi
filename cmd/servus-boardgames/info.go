package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
)

var (
	re_json = regexp.MustCompile(`GEEK.geekitemPreload = (.*);`)
)

func GetBoardgameInfo(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://boardgamegeek.com/boardgame/%d", id)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	res.SetInternal(nil)
	sb := string(body)
	refs := re_json.FindAllStringSubmatch(sb, -1)
	if len(refs) > 0 {
		var rs map[string]any
		err = json.Unmarshal([]byte(refs[0][1]), &rs)
		if err != nil {
			return err
		}

		if val, ok := rs["item"]; ok {
			raw, err := json.Marshal(val)
			if err != nil {
				return err
			}

			var inner map[string]any
			err = json.Unmarshal([]byte(raw), &inner)
			if err != nil {
				return err
			}

			ignore := []string{
				"wiki",
				"description",
				"itemdata",
				"linkedforum_types",
				"summary_video",
				"promoted_ad",
				"special_user",
				"walmart_price",
			}

			for _, item := range ignore {
				delete(inner, item)
			}

			res.SetInternal(inner)
		}
	}

	return nil
}
