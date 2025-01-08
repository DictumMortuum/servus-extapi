package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Val struct {
	Metric struct {
		Instance  string `json:"instance"`
		Interface string `json:"interface"`
		IP        string `json:"ip"`
		Nickname  string `json:"nickname"`
		MAC       string `json:"mac"`
	} `json:"metric"`
}

func process(DB *sqlx.DB, result model.Value) error {
	vectorVal := result.(model.Vector)

	if len(vectorVal) == 0 {
		return nil
	}

	err := Reset(DB)
	if err != nil {
		return err
	}

	for _, item := range vectorVal {
		raw, err := item.MarshalJSON()
		if err != nil {
			return err
		}

		var v Val
		err = json.Unmarshal(raw, &v)
		if err != nil {
			return err
		}

		_, err = Insert(DB, v)
		if err != nil {
			return err
		}

		err = Status(DB, v.Metric.MAC, 1)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	client, err := api.NewClient(api.Config{
		Address: "https://prometheus.dictummortuum.com",
	})
	if err != nil {
		log.Fatalf("Error creating client: %v\n", err)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, _, err := v1api.Query(ctx, "deco_client{type='phone'}", time.Now(), v1.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Error querying Prometheus: %v\n", err)
	}

	DB, err := db.DatabaseX()
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	switch result.Type() {
	case model.ValVector:
		{
			err = process(DB, result)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	devices, err := Select(DB)
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		device.Write(os.Stdout)
	}
}
