package bgg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/jmoiron/sqlx"
)

var (
	re_json = regexp.MustCompile(`GEEK.geekitemPreload = (.*);`)
)

func BoardgameInfo(db *sqlx.DB, id int64) (map[string]any, error) {
	url := fmt.Sprintf("https://boardgamegeek.com/boardgame/%d", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	sb := string(body)
	refs := re_json.FindAllStringSubmatch(sb, -1)
	if len(refs) > 0 {
		var rs map[string]any
		err = json.Unmarshal([]byte(refs[0][1]), &rs)
		if err != nil {
			return nil, err
		}

		if val, ok := rs["item"]; ok {
			raw, err := json.Marshal(val)
			if err != nil {
				return nil, err
			}

			var inner map[string]any
			err = json.Unmarshal([]byte(raw), &inner)
			if err != nil {
				return nil, err
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

			return inner, nil
		}
	}

	return nil, nil
}
