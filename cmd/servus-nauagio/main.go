package main

import (
	"log"

	"github.com/DictumMortuum/servus-extapi/pkg/nas"
)

func main() {
	episodes := scrapeEpisodes("https://www.megatv.com/ekpompes/1110551/to-nayagio/")

	if len(episodes) == 0 {
		return
	}

	for _, episode := range episodes[0:1] {
		payload := map[string]any{
			"url":   episode,
			"path":  "/volume1/plex/greek series/Nauagio/",
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
