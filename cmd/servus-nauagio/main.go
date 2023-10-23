package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/nas"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	DB, err := db.DatabaseX()
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	episodes := scrapeEpisodes("https://www.megatv.com/ekpompes/1110551/to-nayagio/")

	if len(episodes) == 0 {
		return
	}

	for _, episode := range episodes {
		rs, _ := nas.Exists(DB, episode)
		if rs == nil {
			_, err := nas.Insert(DB, episode)
			if err != nil {
				log.Fatal(err)
			}

			payload := map[string]any{
				"url":   episode,
				"path":  "/volume1/plex/greek series/Nauagio/",
				"owner": "dimitris@dictummortuum.com",
				"group": "dimitris@dictummortuum.com",
			}

			err = nas.YoutubeDL(payload)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(payload)
		}
	}
}
