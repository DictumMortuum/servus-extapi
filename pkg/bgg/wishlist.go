package bgg

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type Status struct {
	Own              string `xml:"own,attr" json:"own"`
	PrevOwned        string `xml:"prevowned,attr" json:"prevowned"`
	ForTrade         string `xml:"fortrade,attr" json:"fortrade"`
	Want             string `xml:"want,attr" json:"want"`
	WantToBuy        string `xml:"wanttobuy,attr" json:"wanttobuy"`
	Wishlist         string `xml:"wishlist,attr" json:"wishlist"`
	WishlistPriority string `xml:"wishlistpriority,attr" json:"wishlistpriority"`
	Preordered       string `xml:"preordered,attr" json:"preordered"`
}

type WishlistItem struct {
	// XMLName  xml.Name `xml:"item" json:"-"`
	ObjectId int64  `xml:"objectid,attr" json:"id"`
	Name     string `xml:"name" json:"name"`
	Status   Status `xml:"status"`
	Stats    Info   `json:"stats"`
}

type WishlistRs struct {
	// XMLName xml.Name       `xml:"items" json:"-"`
	Items []WishlistItem `xml:"item" json:"items"`
}

func Wishlist(g *gorm.DB, db *sqlx.DB, name string) ([]WishlistItem, error) {
	link := fmt.Sprintf("https://api.geekdo.com/xmlapi2/collection?username=%s", url.QueryEscape(name))
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	conn := &http.Client{}

	resp, err := conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tmp := WishlistRs{}

	err = xml.Unmarshal(body, &tmp)
	if err != nil {
		return nil, err
	}

	rs := []WishlistItem{}
	for _, item := range tmp.Items {
		info, err := getInfoForBoardgame(db, item.ObjectId)
		if err == sql.ErrNoRows {
			info2, err := BoardgameInfo(db, item.ObjectId)
			if err != nil {
				return nil, err
			}

			log.Println(item.ObjectId)

			raw, err := json.Marshal(info2)
			if err != nil {
				return nil, err
			}

			boardgame := model.Boardgame{}
			_, err = boardgame.Create(g, raw)
			if err != nil {
				return nil, err
			}

			info, err = getInfoForBoardgame(db, item.ObjectId)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}

		if info.Subtype != nil {
			if *info.Subtype == "\"boardgameexpansion\"" {
				continue
			}
		} else {
			continue
		}

		item.Stats = *info
		rs = append(rs, item)
	}

	return rs, nil
}

type Info struct {
	Id             int64                  `json:"id,omitempty"`
	Name           string                 `json:"name,omitempty"`
	Rank           models.JsonNullInt64   `json:"rank,omitempty"`
	Square200      models.JsonNullString  `json:"url,omitempty"`
	Weight         models.JsonNullFloat64 `json:"weight,omitempty"`
	Average        model.JsonNullString   `json:"average,omitempty"`
	MinPlayTime    model.JsonNullString   `db:"minplaytime" json:"minplaytime"`
	MaxPlayTime    model.JsonNullString   `db:"maxplaytime" json:"maxplaytime"`
	MinPlayers     model.JsonNullString   `db:"minplayers" json:"minplayers"`
	MaxPlayers     model.JsonNullString   `db:"maxplayers" json:"maxplayers"`
	BestMinPlayers model.JsonNullString   `db:"bestminplayers" json:"bestminplayers"`
	BestMaxPlayers model.JsonNullString   `db:"bestmaxplayers" json:"bestmaxplayers"`
	Subtype        *string                `db:"subtype" json:"subtype"`
	Cooperative    bool                   `json:"cooperative"`
}

func getInfoForBoardgame(db *sqlx.DB, id int64) (*Info, error) {
	var rs Info
	err := db.Get(&rs, `
		select
			g.id,
			g.name,
			g.rank,
			g.square200,
			json_extract(g.bgg_data, '$.polls.boardgameweight.averageweight') weight,
			json_extract(g.bgg_data, '$.stats.average') average,
			json_extract(g.bgg_data, '$.minplaytime') minplaytime,
			json_extract(g.bgg_data, '$.maxplaytime') maxplaytime,
			json_extract(g.bgg_data, '$.minplayers') minplayers,
			json_extract(g.bgg_data, '$.maxplayers') maxplayers,
			json_extract(g.bgg_data, '$.polls.userplayers.best[0].min') bestminplayers,
			json_extract(g.bgg_data, '$.polls.userplayers.best[0].max') bestmaxplayers,
			json_extract(g.bgg_data, '$.subtype') subtype,
			cooperative
		from
			tboardgames g
		where
			g.id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	return &rs, nil
}
