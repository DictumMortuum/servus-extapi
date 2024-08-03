package model

import (
	"fmt"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	v.Internal = make(map[string]any)
	v.Headers = make(map[string]string)
	c.Set(key, &v)
	return &v, nil
}

type Map struct {
	Internal map[string]any
	Headers  map[string]string
	Paginate func(*gorm.DB) *gorm.DB
	Sort     func(*gorm.DB) *gorm.DB
	Filter   func(*gorm.DB) *gorm.DB
}

func (m *Map) Close() error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	db.Close()

	return nil
}

func (m *Map) Unset(key string) {
	delete(m.Internal, key)
}

func (m *Map) Set(key string, val any) {
	if m.Internal == nil {
		m.Internal = make(map[string]any)
	}

	m.Internal[key] = val
}

func (m *Map) SetInternal(val map[string]any) {
	m.Internal = val
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

func (m *Map) GetByte(key string) ([]byte, error) {
	if val, ok := m.Internal[key]; ok {
		switch s := val.(type) {
		case []byte:
			return s, nil
		default:
			return nil, fmt.Errorf("not a valid type")
		}
	} else {
		return nil, fmt.Errorf("could not find key")
	}
}

func (m *Map) Get(key string) (any, bool) {
	val, ok := m.Internal[key]
	return val, ok
}

func (m *Map) GetModel() (Routable, error) {
	if val, ok := m.Internal["apimodel"]; ok {
		conn, ok := val.(Routable)
		if !ok {
			return nil, fmt.Errorf("error with retrieving routable")
		} else {
			return conn, nil
		}
	}

	return nil, nil
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

func (m *Map) GetGorm() (*gorm.DB, error) {
	if val, ok := m.Internal["gorm"]; ok {
		conn, ok := val.(*gorm.DB)
		if !ok {
			return nil, fmt.Errorf("error with retrieving gorm pointer")
		} else {
			return conn, nil
		}
	} else {
		DB, err := m.GetDB()
		if err != nil {
			return nil, err
		}

		db, err := gorm.Open(mysql.New(mysql.Config{
			Conn:              DB,
			DefaultStringSize: 512,
		}), &gorm.Config{
			// DisableAutomaticPing: true,
			// PrepareStmt: false,
		})
		if err != nil {
			return nil, err
		}

		m.Internal["gorm"] = db
		return db, nil
	}
}

func (m *Map) GetRedis() (*redis.Client, error) {
	if val, ok := m.Internal["redis"]; ok {
		conn, ok := val.(*redis.Client)
		if !ok {
			return nil, fmt.Errorf("error with retrieving redis pointer")
		} else {
			return conn, nil
		}
	} else {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		m.Internal["redis"] = rdb
		return rdb, nil
	}
}
