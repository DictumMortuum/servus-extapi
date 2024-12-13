package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type CollectionItem struct {
	Id            int64          `json:"gameId"`
	Name          string         `json:"name"`
	Image         string         `json:"image"`
	Thumbnail     string         `json:"thumbnail"`
	MinPlayers    int            `json:"minPlayers"`
	MaxPlayers    int            `json:"maxPlayers"`
	PlayingTime   int            `json:"playingTime"`
	IsExpansion   bool           `json:"isExpansion"`
	YearPublished int            `json:"yearPublished"`
	BggRating     float64        `json:"bggRating"`
	AverageRating float64        `json:"averageRating"`
	Rank          int            `json:"rank"`
	NumPlays      int            `json:"numPlays"`
	Rating        float64        `json:"rating"`
	Owned         bool           `json:"owned"`
	PreOrdered    bool           `json:"preOrdered"`
	ForTrade      bool           `json:"forTrade"`
	PreviousOwned bool           `json:"previousOwned"`
	Want          bool           `json:"want"`
	WantToPlay    bool           `json:"wantToPlay"`
	WantToBuy     bool           `json:"wantToBuy"`
	Wishlist      bool           `json:"wishlist"`
	UserComment   string         `json:"userComment"`
	Info          *BoardgameInfo `json:"info"`
}

func getCollection(username string) ([]CollectionItem, error) {
	r, err := http.NewRequest("GET", fmt.Sprintf("http://bgg-json.azurewebsites.net/collection/%s", username), nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rs []CollectionItem
	err = json.Unmarshal(body, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

type BoardgameInfo struct {
	Id             int64                  `json:"id,omitempty"`
	Mechanics      models.JsonArray       `json:"mechanics,omitempty"`
	Designers      models.JsonArray       `json:"designers,omitempty"`
	Categories     models.JsonArray       `json:"categories,omitempty"`
	Subdomains     models.JsonArray       `json:"subdomains,omitempty"`
	Families       models.JsonArray       `json:"families,omitempty"`
	Weight         models.JsonNullFloat64 `json:"weight,omitempty"`
	BestMinPlayers models.JsonNullInt64   `json:"best_min_players,omitempty"`
	BestMaxPlayers models.JsonNullInt64   `json:"best_max_players,omitempty"`
}

func getBoardgameInfo(DB *sqlx.DB, RDB *redis.Client, id int64) (*BoardgameInfo, error) {
	var rs BoardgameInfo
	key := fmt.Sprintf("boardgameInfo%d", id)

	err := db.Get(RDB, key, &rs)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	err = DB.QueryRowx(`
		select
			g.id,
			json_extract(g.bgg_data, '$.links.boardgamemechanic') mechanics,
			json_extract(g.bgg_data, '$.links.boardgamedesigner') designers,
			json_extract(g.bgg_data, '$.links.boardgamecategory') categories,
			json_extract(g.bgg_data, '$.links.boardgamefamily') families,
			json_extract(g.bgg_data, '$.links.boardgamesubdomain') subdomains,
			json_extract(g.bgg_data, '$.polls.boardgameweight.averageweight') weight,
			json_extract(g.bgg_data, '$.polls.userplayers.best[0].min') best_min_players,
			json_extract(g.bgg_data, '$.polls.userplayers.best[0].max') best_max_players
		from
			tboardgames g
		where
			g.id = ?
	`, id).StructScan(&rs)
	if err != nil {
		return nil, err
	}

	err = db.SetT(RDB, key, rs, 5*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}

func GetCollectionInfo(req *model.Map, res *model.Map) error {
	body, err := req.GetByte("body")
	if err != nil {
		return err
	}

	type Body struct {
		Ids []int64 `json:"ids"`
	}

	var b Body
	err = json.Unmarshal(body, &b)
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	RDB, err := req.GetRedis()
	if err != nil {
		return err
	}

	retval := []*BoardgameInfo{}
	for _, id := range b.Ids {
		info, err := getBoardgameInfo(DB, RDB, id)
		log.Println(id, info)
		if err != nil {
			continue
		}

		retval = append(retval, info)
	}

	res.Set("data", retval)

	return nil
}

func GetPlayerCollection(req *model.Map, res *model.Map) error {
	id, err := req.GetString("id")
	if err != nil {
		return err
	}

	if id == "" {
		return errors.New("username is not valid")
	}

	rs, err := getCollection(id)
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	RDB, err := req.GetRedis()
	if err != nil {
		return err
	}

	retval := []CollectionItem{}
	for _, item := range rs {
		info, err := getBoardgameInfo(DB, RDB, item.Id)
		if err != nil {
			item.Info = nil
		}

		item.Info = info
		retval = append(retval, item)
	}

	res.Set("data", retval)

	return nil
}
