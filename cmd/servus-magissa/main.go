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

	episodes := scrapeEpisodes("https://www.antenna.gr/magissa")

	if len(episodes) == 0 {
		return
	}

	// version

	for _, episode := range episodes {
		part := "https://www.antenna.gr" + episode

		rs, _ := nas.Exists(DB, part)
		if rs == nil {
			_, err := nas.Insert(DB, part)
			if err != nil {
				log.Fatal(err)
			}

			payload := map[string]any{
				"url":   part,
				"path":  "/volume1/plex/greek series/Magissa/",
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
