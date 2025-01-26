package main

import (
	"encoding/json"

	"github.com/DictumMortuum/servus-extapi/pkg/model"
	"github.com/DictumMortuum/servus-extapi/pkg/queries"
	"github.com/DictumMortuum/servus-extapi/pkg/scrape"
	"github.com/jmoiron/sqlx"
)

func getScrape(DB *sqlx.DB, id int64) (*model.Scrape, error) {
	sc := model.Scrape{}
	err := DB.Get(&sc, `
		select
			*
		from
			tscrape
		where
			id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	return &sc, nil
}

func getScrapeUrl(DB *sqlx.DB, id int64) (*model.ScrapeUrl, error) {
	u := model.ScrapeUrl{}
	err := DB.Get(&u, `
		select
			*
		from
			tscrapeurl
		where
			id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func getScrapeUrls(DB *sqlx.DB, id int64) ([]model.ScrapeUrl, error) {
	u := []model.ScrapeUrl{}
	err := DB.Select(&u, `
		select
			*
		from
			tscrapeurl
		where
			scrape_id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func setUrlToPending(DB *sqlx.DB, id int64) error {
	_, err := DB.Exec(`
		update
			tscrapeurl
		set
			last_scraped = NULL,
			last_instock = NULL,
			last_preorder = NULL,
			last_outofstock = NULL,
			last_pages = NULL
		where
			id = ?
	`, id)
	if err != nil {
		return err
	}

	return nil
}

func scrapeSingle(req *model.Map, r *scrape.GenericScrapeRequest) error {
	conn, err := req.GetRmq()
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	scrape, err := conn.OpenQueue("scrape")
	if err != nil {
		return err
	}

	raw, err := json.Marshal(r)
	if err != nil {
		return err
	}

	err = scrape.Publish(string(raw))
	if err != nil {
		return err
	}

	err = setUrlToPending(DB, r.ScrapeUrl.Id)
	if err != nil {
		return err
	}

	return nil
}

type scrapeBody struct {
	ListOnly bool `json:"list_only"`
}

func ScrapeUrl(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	body, err := req.GetByte("body")
	if err != nil {
		return err
	}

	var payload scrapeBody
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return err
	}

	u, err := getScrapeUrl(DB, id)
	if err != nil {
		return err
	}

	cfg, err := queries.GetConfig(DB, "SCRAPE_CACHE")
	if err != nil {
		return err
	}

	r := scrape.GenericScrapeRequest{
		ScrapeUrl: *u,
		Cache:     cfg.Value,
		ListOnly:  payload.ListOnly,
	}

	err = scrapeSingle(req, &r)
	if err != nil {
		return err
	}

	res.SetInternal(map[string]any{
		"req": r,
	})

	return nil
}

func ScrapeStore(req *model.Map, res *model.Map) error {
	id, err := req.GetInt64("id")
	if err != nil {
		return err
	}

	DB, err := req.GetDB()
	if err != nil {
		return err
	}

	body, err := req.GetByte("body")
	if err != nil {
		return err
	}

	var payload scrapeBody
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return err
	}

	sc, err := getScrape(DB, id)
	if err != nil {
		return err
	}

	u, err := getScrapeUrls(DB, sc.Id)
	if err != nil {
		return err
	}

	cfg, err := queries.GetConfig(DB, "SCRAPE_CACHE")
	if err != nil {
		return err
	}

	for _, url := range u {
		r := scrape.GenericScrapeRequest{
			ScrapeUrl: url,
			Cache:     cfg.Value,
			ListOnly:  payload.ListOnly,
		}

		err = scrapeSingle(req, &r)
		if err != nil {
			return err
		}
	}

	return nil
}
