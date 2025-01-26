package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/db"
	"github.com/DictumMortuum/servus-extapi/pkg/scrape"
	"github.com/adjust/rmq/v5"
	"github.com/jmoiron/sqlx"
)

func setUrlToScraped(DB *sqlx.DB, payload map[string]any) error {
	_, err := DB.NamedExec(`
		update
			tscrapeurl
		set
			last_scraped = :scraped,
			last_instock = :instock,
			last_preorder = :preorder,
			last_outofstock = :outofstock,
			last_pages = :pages_visited
		where
			id = :id
	`, payload)
	if err != nil {
		return err
	}

	return nil
}

func consumeFn(task scrape.GenericScrapeRequest) error {
	DB, err := db.DatabaseX()
	if err != nil {
		return err
	}
	defer DB.Close()

	sc, err := getScrape(DB, task.ScrapeUrl.ScrapeId)
	if err != nil {
		return err
	}

	err = scrape.Stale(DB, sc.StoreId)
	if err != nil {
		return err
	}

	payload, _, err := scrape.GenericScrape(*sc, DB, task)
	if err != nil {
		return err
	}

	err = scrape.Cleanup(DB, sc.StoreId)
	if err != nil {
		return err
	}

	log.Println(task.ListOnly, payload)
	err = setUrlToScraped(DB, payload)
	if err != nil {
		return err
	}

	return nil
}

func Consumer(conn rmq.Connection) {
	queue, err := conn.OpenQueue("scrape")
	if err != nil {
		log.Fatal(err)
	}

	err = queue.StartConsuming(10, time.Second)
	if err != nil {
		log.Fatal(err)
	}

	_, err = queue.AddConsumerFunc("scraper", func(d rmq.Delivery) {
		var task scrape.GenericScrapeRequest
		err = json.Unmarshal([]byte(d.Payload()), &task)
		if err != nil {
			d.Reject()
		}

		err = consumeFn(task)
		if err != nil {
			log.Println(err)
			d.Reject()
		}

		d.Ack()
	})
	if err != nil {
		log.Fatal(err)
	}
}
