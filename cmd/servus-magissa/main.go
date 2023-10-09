package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func main() {
	episodes := scrapeEpisodes("https://www.antenna.gr/magissa")

	if len(episodes) == 0 {
		return
	}

	part := "https://www.antenna.gr" + episodes[0]
	video := scrapeVideo(part)
	name := strings.Split(part, "/")

	output, err := json.Marshal(map[string]any{
		"name":  name[len(name)-1],
		"url":   video,
		"path":  "/volume1/plex/greek series/Magissa/",
		"owner": "dimitris@dictummortuum.com",
		"group": "dimitris@dictummortuum.com",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
}
