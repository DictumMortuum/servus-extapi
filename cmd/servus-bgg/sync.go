package main

import (
	"fmt"
	"net/http"

	"github.com/DictumMortuum/servus-extapi/pkg/bgg"
	"github.com/urfave/cli/v2"
)

func sync(c *cli.Context) error {
	rs, err := bgg.GetAllBoardgames()
	if err != nil {
		return err
	}

	for _, item := range rs {
		if v, ok := item["id"]; ok {
			err := update(v.(string))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func update(id string) error {
	url := fmt.Sprintf("https://extapi.dictummortuum.com/boardgames/info/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return fmt.Errorf("%s - %s", resp.Status, id)
	}

	return nil
}
