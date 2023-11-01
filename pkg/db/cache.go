package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func set(rdb *redis.Client, key string, val any) error {
	p, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = rdb.Set(context.Background(), key, p, time.Second*120).Err()
	if err != nil {
		return err
	}

	return nil
}

func get(rdb *redis.Client, key string, dest any) error {
	p, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(p), dest)
}

func CachedSelect(DB *sqlx.DB, RDB *redis.Client, key string, dest any, query string, args ...any) error {
	err := get(RDB, key, dest)
	if err == redis.Nil {
		err = DB.Select(dest, query, args...)
		if err != nil {
			return err
		}

		err = set(RDB, key, dest)
		if err != nil {
			return err
		}
	}

	return nil
}
