package screenshot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/DictumMortuum/servus-extapi/pkg/minio"
	minio_lib "github.com/minio/minio-go/v7"
)

func _req(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, config.Cfg.ApiFlash.Host, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("url", url)
	q.Add("access_key", config.Cfg.ApiFlash.Key)
	q.Add("format", "jpeg")
	q.Add("no_ads", "true")
	q.Add("no_tracking", "true")
	q.Add("no_cookie_banners", "true")
	q.Add("latitude", "37.9838")
	q.Add("longitude", "23.7275")
	q.Add("wait_until", "page_loaded")
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func Do(url string, id int64) error {
	ctx := context.Background()

	req, err := _req(url)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	client, err := minio.Init()
	if err != nil {
		return err
	}

	_, err = client.PutObject(ctx, "screenshots", fmt.Sprintf("wish-%d.jpg", id), res.Body, res.ContentLength, minio_lib.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return err
	}

	return nil
}

func Get(id int64) (*minio_lib.Object, error) {
	ctx := context.Background()

	client, err := minio.Init()
	if err != nil {
		return nil, err
	}

	return client.GetObject(ctx, "screenshots", fmt.Sprintf("wish-%d.jpg", id), minio_lib.GetObjectOptions{})
}
