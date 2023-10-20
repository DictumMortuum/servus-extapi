package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/urfave/cli/v2"
)

/*
	load inserts bgplays that are not associated with a play yet in the database.
*/

func load(c *cli.Context) error {
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
		bgstats, err := getBGStats(DB, bgplay)
		if err != nil {
			return err
		}
		bgplay.BGStats = bgstats

		mapped, err := findMappedPlay(DB, bgplay)
		if err != nil {
			return err
		}

		if len(mapped) == 0 {
			log.Println(bgplay)
			play_id, err := insertPlay(DB, bgplay)
			if err != nil {
				return err
			}

			for _, stat := range bgplay.BGStats {
				err := insertStats(DB, play_id, bgplay.BoardgameId, stat)
				if err != nil {
					return err
				}
			}

			err = updateLink(DB, bgplay, Play{
				Id: play_id,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
