package nas

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func YoutubeDL(payload map[string]any) error {
	output, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST", "https://nas.dictummortuum.com/hooks/youtube-dl", bytes.NewBuffer(output))
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
