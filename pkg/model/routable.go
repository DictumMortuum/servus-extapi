package model

import (
	"gorm.io/gorm"
)

type Routable interface {
	List(*gorm.DB, ...func(*gorm.DB) *gorm.DB) (any, error)
	Get(*gorm.DB, int64) (any, error)
	Update(*gorm.DB, int64, []byte) (any, error)
	Create(*gorm.DB, []byte) (any, error)
	Delete(*gorm.DB, int64) (any, error)
	DefaultFilter(*gorm.DB) *gorm.DB
}
