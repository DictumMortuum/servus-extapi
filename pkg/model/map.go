package model

import (
	"fmt"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
)

func ToMap(c *gin.Context, key string) (*Map, error) {
	if val, ok := c.Get(key); ok {
		switch v := val.(type) {
		case *Map:
			return v, nil
		default:
			return nil, fmt.Errorf("req already exists in the context")
		}
	}

	v := Map{}
	c.Set(key, &v)
	return &v, nil
}

type Map struct {
	Internal map[string]any
}

func (m *Map) Close() error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	db.Close()

	return nil
}

func (m *Map) Set(key string, val any) {
	if m.Internal == nil {
		m.Internal = make(map[string]any)
	}

	m.Internal[key] = val
}

func (m *Map) GetString(key string) (string, error) {
	return cast.ToStringE(m.Internal[key])
}

func (m *Map) GetInt64(key string) (int64, error) {
	return cast.ToInt64E(m.Internal[key])
}

func (m *Map) GetBool(key string) (bool, error) {
	return cast.ToBoolE(m.Internal[key])
}

func (m *Map) Get(key string) (any, bool) {
	val, ok := m.Internal[key]
	return val, ok
}

func (m *Map) GetDB() (*sqlx.DB, error) {
	if val, ok := m.Internal["db"]; ok {
		conn, ok := val.(*sqlx.DB)
		if !ok {
			return nil, fmt.Errorf("error with retrieving database pointer")
		} else {
			return conn, nil
		}
	} else {
		db, err := sqlx.Open("mysql", config.Cfg.Databases["mariadb"])
		if err != nil {
			return nil, err
		}

		db.MapperFunc(util.ToSnake)
		m.Internal["db"] = db
		return db, nil
	}
}
