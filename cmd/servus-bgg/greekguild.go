package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DictumMortuum/servus-extapi/pkg/bgg"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/urfave/cli/v2"
)

func getRank(cfg map[string]any) string {
	rankinfo_tmp := cfg["rankinfo"].([]any)

	if len(rankinfo_tmp) == 0 {
		return "N/A"
	}

	for _, rank := range rankinfo_tmp {
		rankinfo := rank.(map[string]any)
		name := rankinfo["shortprettyname"].(string)

		if name == "Overall Rank" {
			return rankinfo["rank"].(string)
		}
	}

	return "N/A"
}

func getDesigners(cfg map[string]any) string {
	links := cfg["links"].(map[string]any)
	designers_tmp := links["boardgamedesigner"].([]any)

	if len(designers_tmp) == 0 {
		return "N/A"
	}

	designers := []string{}
	for _, designer := range designers_tmp {
		designerinfo := designer.(map[string]any)
		name := designerinfo["name"].(string)
		designers = append(designers, name)
	}

	return strings.Join(designers, "|")
}

func getBoardgameCategories(cfg map[string]any) string {
	links := cfg["links"].(map[string]any)
	tmp := links["boardgamecategory"].([]any)

	if len(tmp) == 0 {
		return "N/A"
	}

	categories := []string{}
	for _, category := range tmp {
		info := category.(map[string]any)
		name := info["name"].(string)
		categories = append(categories, name)
	}

	return strings.Join(categories, "|")
}

func guild(c *cli.Context) error {
	ctx := context.Background()
	token := "b5cc5df6-26db-4dd6-91ac-4e15da314116"

	cli := bgg.NewClient(token)
	items, err := cli.GetGeeklistItems(ctx, 365879, false)
	if err != nil {
		return err
	}

	// id="12252456"
	// objecttype="thing"
	// subtype="boardgame"
	// objectid="70512"
	// objectname="Luna"
	// username="Vagos"
	// postdate="Wed, 29 Oct 2025 14:55:49 +0000"
	// editdate="Wed, 29 Oct 2025 14:55:49 +0000"
	// thumbs="2"

	fmt.Printf("bgg_id,name,url,username,thumbs,min_player,max_player,weight,yearpublished,subdomain,rank,designers,categories\n")
	for _, it := range items {
		raw, err := getBoardgameInfo("boardgamegeek", it.ObjectID)
		if err != nil {
			return err
		}

		var cfg map[string]any
		err = json.Unmarshal(raw, &cfg)
		if err != nil {
			return err
		}

		stats := cfg["stats"].(map[string]any)
		links := cfg["links"].(map[string]any)
		subdomain_tmp := links["boardgamesubdomain"].([]any)

		var subdomain string
		if len(subdomain_tmp) > 0 {
			sub := subdomain_tmp[0].(map[string]any)
			subdomain = sub["name"].(string)
		} else {
			subdomain = "N/A"
		}

		fmt.Printf("%d,%s,%s,%s,%d,%d,%d,%f,%s,%s,%s,%s,%s\n",
			it.ObjectID,
			strings.Trim(it.ObjectName, ","),
			fmt.Sprintf("https://boardgamegeek.com/boardgame/%d", it.ObjectID),
			it.Username,
			it.Thumbs,
			util.Atoi(cfg["minplayers"].(string)),
			util.Atoi(cfg["maxplayers"].(string)),
			util.Atof(stats["avgweight"].(string)),
			cfg["yearpublished"],
			subdomain,
			getRank(cfg),
			getDesigners(cfg),
			getBoardgameCategories(cfg),
		)
	}

	return nil
}
