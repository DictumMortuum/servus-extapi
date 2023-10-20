package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
)

/*
	integrity checks that all bgplays are correctly mapped to plays.
*/

func integrity(c *cli.Context) error {
	DB, err := db.DatabaseX()
	if err != nil {
		return err
	}
	defer DB.Close()

	rs, err := getMappedBGPlays(DB)
	if err != nil {
		return err
	}

	for _, bgplay := range rs {
		bgstats, err := getBGStats(DB, bgplay)
		if err != nil {
			return err
		}
		bgplay.BGStats = bgstats

		play, err := getPlay(DB, bgplay)
		if err != nil {
			return err
		}

		stats, err := getStats(DB, *play)
		if err != nil {
			return err
		}
		play.Stats = stats

		err = integrityCheck(DB, bgplay, *play)
		if err != nil {
			return err
		}
	}

	return nil
}

func integrityCheck(DB *sqlx.DB, bgplay BGPlay, play Play) error {
	if !comparePlays(bgplay, play) {
		log.Println(bgplay.Id, play.Id)
		log.Println(play.Stats)
		log.Println(bgplay.BGStats)
	}

	if !compareScores(bgplay, play) {
		log.Println(bgplay.Id, play.Id)
		log.Println(play.Stats)
		log.Println(bgplay.BGStats)

		for i := 0; i < len(play.Stats); i++ {
			if bgplay.BGStats[i].PlayerId == play.Stats[i].PlayerId || bgplay.BGStats[i].PlayerId == 79 {
				if bgplay.BGStats[i].Score.Int64 == play.Stats[i].Score.Int64 {
					log.Println("Updating winner")
					err := updateWinner(DB, bgplay.BGStats[i], play.Stats[i])
					if err != nil {
						return err
					}
				}
			}
		}
	}

	flag := false
	if bgplay.LocationId.Valid {
		if play.LocationId.Valid {
			if play.LocationId.Int64 != bgplay.LocationId.Int64 {
				flag = true
			}
		} else {
			flag = true
		}
	}

	if flag {
		log.Println("Updating ", play.Id, " to ", bgplay.LocationId.Int64, " from ", play.LocationId.Int64)
		err := updateLocation(DB, bgplay, play)
		if err != nil {
			return err
		}
	}

	return nil
}
