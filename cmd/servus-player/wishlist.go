package main

import (
	"encoding/json"
	"errors"

	"github.com/DictumMortuum/servus-extapi/pkg/bgg"
	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/redis/go-redis/v9"
)

func GetPlayerWishlist(req *model.Map, res *model.Map) error {
	id, err := req.GetString("id")
	if err != nil {
		return err
	}

	if id == "" {
		return errors.New("username is not valid")
	}

	RDB, err := req.GetRedis()
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	key := "GetPlayerWishlist" + id
	var wishlist *bgg.WishlistRs
	err = db.Get(RDB, key, &wishlist)
	if err == redis.Nil || len(wishlist.Items) == 0 {
		rs, err := bgg.Wishlist(id)
		if err != nil {
			return err
		}

		err = db.Set(RDB, key, rs)
		if err != nil {
			return err
		}

		wishlist = rs
	}

	if len(wishlist.Items) > 0 {
		raw, err := json.Marshal(wishlist.Items)
		if err != nil {
			return err
		}

		_, err = DB.NamedExec(`
		update tboardgameplayers set collection = :collection where bgg_username = :id
		`, map[string]any{
			"id":         id,
			"collection": string(raw),
		})
		if err != nil {
			return err
		}

		res.Set("synced", true)
	} else {
		res.Set("synced", false)
	}

	res.Set("data", wishlist.Items)

	return nil
}
