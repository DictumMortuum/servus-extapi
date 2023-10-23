package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/nas"
)

func main() {
	episodes := scrapeEpisodes("https://www.antenna.gr/magissa")

	if len(episodes) == 0 {
		return
	}

	for _, episode := range episodes {
		part := "https://www.antenna.gr" + episode

		payload := map[string]any{
			"url":   part,
			"path":  "/volume1/plex/greek series/Magissa/",
			"owner": "dimitris@dictummortuum.com",
			"group": "dimitris@dictummortuum.com",
		}

		err := nas.YoutubeDL(payload)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(payload)
	}
}
