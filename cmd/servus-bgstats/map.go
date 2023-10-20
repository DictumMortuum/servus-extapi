package main

import (
	"log"
	"sort"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/urfave/cli/v2"
)

/*
	mapMissing grabs all bgplays and tries to map them to existing plays in the database.
*/

func mapMissing(c *cli.Context) error {
	DB, err := db.DatabaseX()
	if err != nil {
		return err
	}
	defer DB.Close()

	rs, err := getUnmappedBGPlays(DB)
	if err != nil {
		return err
	}

	for _, bgplay := range rs {
		similar, err := getSimilarBGPlays(DB, bgplay)
		if err != nil {
			return err
		}

		if len(similar) > 1 {
			// log.Println("Found multiple plays for ", bgplay.Id, " : ", similar)
			continue
		}

		bgstats, err := getBGStats(DB, bgplay)
		if err != nil {
			return err
		}
		bgplay.BGStats = bgstats

		mapped, err := findMappedPlay(DB, bgplay)
		if err != nil {
			return err
		}

		for _, play := range mapped {
			stats, err := getStats(DB, play)
			if err != nil {
				return err
			}
			play.Stats = stats

			// if comparePlays(bgplay, play) && compareScores(bgplay, play) {
			if comparePlays(bgplay, play) {
				if !bgplay.PlayId.Valid {
					log.Println("Compare plays", play.Id, " <- ", bgplay.Id)
					log.Println("play", play.Stats)
					log.Println("bgplay", bgplay.BGStats)

					// err := updateLink(DB, bgplay, play)
					// if err != nil {
					// 	return err
					// }
				}
			}
		}
	}

	return nil
}

func comparePlays(bgplay BGPlay, play Play) bool {
	sort.Slice(play.Stats, func(i int, j int) bool {
		if play.Stats[i].Score.Int64 != play.Stats[j].Score.Int64 {
			return play.Stats[i].Score.Int64 > play.Stats[j].Score.Int64
		} else {
			return play.Stats[i].PlayerId > play.Stats[j].PlayerId
		}
	})

	sort.Slice(bgplay.BGStats, func(i int, j int) bool {
		if bgplay.BGStats[i].Score.Int64 != bgplay.BGStats[j].Score.Int64 {
			return bgplay.BGStats[i].Score.Int64 > bgplay.BGStats[j].Score.Int64
		} else {
			return bgplay.BGStats[i].PlayerId > bgplay.BGStats[j].PlayerId
		}
	})

	for i := 0; i < len(play.Stats); i++ {
		if bgplay.BGStats[i].PlayerId != play.Stats[i].PlayerId && bgplay.BGStats[i].PlayerId != 79 {
			return false
		}
	}

	return true
}

func compareScores(bgplay BGPlay, play Play) bool {
	flags := []bool{}

	for i := 0; i < len(play.Stats); i++ {
		flags = append(flags, false)
	}

	for i := 0; i < len(play.Stats); i++ {
		for j := 0; j < len(bgplay.BGStats); j++ {
			if play.Stats[i].Score.Int64 == bgplay.BGStats[j].Score.Int64 && !flags[i] {
				if play.Stats[i].Winner == bgplay.BGStats[j].Winner {
					flags[i] = true
					break
				}
			}
		}
	}

	for _, flag := range flags {
		if !flag {
			return false
		}
	}

	return true
}
