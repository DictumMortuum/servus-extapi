package bgg

import (
	"github.com/DictumMortuum/servus/pkg/boardgames"
	"github.com/DictumMortuum/servus/pkg/boardgames/bgg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type mapping struct {
	Id          int64  `db:"id" json:"id"`
	BoardgameId int64  `db:"boardgame_id" json:"boardgame_id"`
	Name        string `db:"name" json:"name"`
}

func SearchCachedPriceOnBgg(c *gin.Context, db *gorm.DB) (interface{}, error) {
	rs := []mapping{}

	if val, ok := c.Get("mapped_data"); ok {
		price, _ := val.(map[string]any)
		name := boardgames.TransformName(price["name"].(string))
		bgg_results, err := bgg.Search(name)
		if err != nil {
			return nil, err
		}

		for _, result := range bgg_results {
			rs = append(rs, mapping{
				Id:          -1,
				BoardgameId: result.Id,
				Name:        result.Name.Value,
			})
		}
	}

	return rs, nil
}
