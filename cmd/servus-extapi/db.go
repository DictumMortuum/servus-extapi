package main

// func Database() (*gorm.DB, *sql.DB, error) {
// 	dsn := config.Cfg.Databases["mariadb"]

// 	sqlDB, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	db, err := gorm.Open(mysql.New(mysql.Config{
// 		Conn:              sqlDB,
// 		DefaultStringSize: 512,
// 	}), &gorm.Config{
// 		PrepareStmt: true,
// 	})
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	db.AutoMigrate(
// 		&model.Boardgame{},
// 		&model.Play{},
// 		&model.Player{},
// 		&model.Stat{},
// 		&model.BoardgamePrice{},
// 		&model.Store{},
// 		&model.Location{},
// 		&model.BGStatsPlayer{},
// 		&model.BGStatsLocation{},
// 		&model.BGStatsGame{},
// 		&model.BGStatsPlay{},
// 		&model.BGStat{},
// 		&model.Price{},
// 		&model.IgnoredPrice{},
// 		&model.CachedPrice{},
// 		&model.IgnoredName{},
// 	)
// 	return db, sqlDB, nil
// }
