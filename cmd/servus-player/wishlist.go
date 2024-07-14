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

	G, err := req.GetGorm()
	if err != nil {
		return err
	}

	key := "GetPlayerWishlist" + id
	var wishlist []bgg.WishlistItem
	err = db.Get(RDB, key, &wishlist)
	if err == redis.Nil || len(wishlist) == 0 {
		rs, err := bgg.Wishlist(G, DB, id)
		if err != nil {
			return err
		}

		err = db.Set(RDB, key, rs)
		if err != nil {
			return err
		}

		wishlist = rs
	}

	if len(wishlist) > 0 {
		raw, err := json.Marshal(wishlist)
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

	res.Set("data", wishlist)

	return nil
}

func GetFinderUserWishlist(req *model.Map, res *model.Map) error {
	id, err := req.GetString("id")
	if err != nil {
		return err
	}

	force, err := req.GetBool("force")
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

	G, err := req.GetGorm()
	if err != nil {
		return err
	}

	key := "GetFinderUserWishlist" + id
	var wishlist []bgg.WishlistItem
	err = db.Get(RDB, key, &wishlist)
	if err == redis.Nil || len(wishlist) == 0 {
		rs, err := bgg.Wishlist(G, DB, id)
		if err != nil {
			return err
		}

		err = db.Set(RDB, key, rs)
		if err != nil {
			return err
		}

		wishlist = rs
	}

	if len(wishlist) > 0 || force {
		ratings, err := bgg.GetRatingsFromBgg(id)
		if err != nil {
			return err
		}

		rs := []bgg.WishlistItem{}
		for _, item := range wishlist {
			for _, rating := range ratings {
				if rating.Id == int(item.ObjectId) {
					item.UserRating = rating.Rating
					break
				}
			}

			rs = append(rs, item)
		}

		raw, err := json.Marshal(rs)
		if err != nil {
			return err
		}

		_, err = DB.NamedExec(`
		update tfinderusers set collection = :collection where bgg_username = :id
		`, map[string]any{
			"id":         id,
			"collection": string(raw),
		})
		if err != nil {
			return err
		}

		res.Set("synced", true)
		res.Set("data", rs)
	} else {
		res.Set("synced", false)
		res.Set("data", wishlist)
	}

	return nil
}
