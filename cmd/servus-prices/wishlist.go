package main

// import (
// 	"strings"

// 	"github.com/mrz1836/go-sanitize"
// 	"gorm.io/gorm"
// )

// func addSearch(val string) func(*gorm.DB) *gorm.DB {
// 	terms := strings.Split(sanitize.AlphaNumeric(val, true), " ")

// 	return func(db *gorm.DB) *gorm.DB {
// 		for _, term := range terms {
// 			term = "%" + term + "%"
// 			db = db.Where("name COLLATE utf8mb4_unicode_ci LIKE ?", term)
// 		}

// 		return db
// 	}
// }

// func wishlistFilter(c *gin.Context) {
// 	username := c.Query("username")

// 	m, err := model.ToMap(c, "req")
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	db, err := m.GetGorm()
// 	if err != nil {
// 		c.Error(err)
// 		return
// 	}

// 	if username != "" {
// 		col, err := bgg.Wishlist("Dictum Mortuum")
// 		if err != nil {
// 			c.Error(err)
// 			return
// 		}

// 		for i, item := range col.Items {
// 			tx := db.Session(&gorm.Session{NewDB: true}).Unscoped()

// 			if i == 0 {
// 				db = db.Where(tx.Scopes(addSearch(item.Name)))
// 			} else {
// 				db = db.Or(tx.Scopes(addSearch(item.Name)))
// 			}
// 		}

// 		db = db.Debug()
// 		m.Set("gorm", db)
// 	}
// }
