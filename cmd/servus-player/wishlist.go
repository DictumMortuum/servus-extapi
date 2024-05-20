package main

import (
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

	RDB, err := req.GetRedis()
	if err != nil {
		return err
	}

	key := "GetPlayerWishlist" + id
	var wishlist *bgg.WishlistRs
	err = db.Get(RDB, key, &wishlist)
	if err == redis.Nil {
		rs, err := bgg.Wishlist(id)
		if err != nil {
			return err
		}

		err = db.Set(RDB, key, rs)
		if err != nil {
			return err
		}
	}

	res.Set("wishlist", wishlist)

	return nil
}
