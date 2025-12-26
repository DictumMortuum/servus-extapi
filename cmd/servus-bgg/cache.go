package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/bgg"
	"github.com/DictumMortuum/servus-extapi/pkg/minio"
	"github.com/DictumMortuum/servus-extapi/pkg/util"
	minio_lib "github.com/minio/minio-go/v7"
	"github.com/urfave/cli/v2"
)

// uploadIfMissingOrStale uploads the object if it doesn't exist,
// or if it exists but hasn't been updated for at least one month.
func uploadIfMissingOrStale(ctx context.Context, cli *minio_lib.Client, bucket, objectName string) (didUpload bool, err error) {
	// 1) Check object metadata
	info, statErr := cli.StatObject(ctx, bucket, objectName, minio_lib.StatObjectOptions{})
	shouldUpload := false

	if statErr != nil {
		resp := minio_lib.ToErrorResponse(statErr)
		// Object missing — go ahead and upload
		if resp.Code == "NoSuchKey" || resp.Code == "NotFound" || resp.StatusCode == 404 {
			shouldUpload = true
		} else {
			// Some other error (e.g., permission issue)
			return false, statErr
		}
	} else {
		// Object exists — only upload if it's older than one month
		if info.LastModified.Before(time.Now().AddDate(0, -1, 0)) {
			shouldUpload = true
		}
	}

	if !shouldUpload {
		return false, nil // Nothing to do
	}
	return true, nil
}

func cache(c *cli.Context) error {
	rs, err := bgg.GetAllBoardgameIds()
	if err != nil {
		return err
	}

	bucket := "boardgamegeek"

	for _, item := range rs {
		time.Sleep(200 * time.Millisecond)
		log.Println(item, len(rs))
		ctx := context.Background()
		id := util.Atoi64(item)
		name := fmt.Sprintf("%d.json", id)

		client, err := minio.Init()
		if err != nil {
			return err
		}

		shouldUpload, err := uploadIfMissingOrStale(ctx, client, bucket, name)
		if err != nil {
			return err
		}

		if !shouldUpload {
			log.Println("Skipped upload (object not old enough).")
			continue
		}

		raw, err := bgg.BoardgameInfoNoDB(id)
		if err != nil {
			return err
		}

		log.Println("Uploaded (object was stale).")
		// 2) Upload (overwrite) the object
		reader := bytes.NewReader(raw)
		_, err = client.PutObject(ctx, bucket, name, reader, int64(len(raw)), minio_lib.PutObjectOptions{
			ContentType: "application/json",
		})
		if err != nil {
			return err
		}

	}

	return nil
}
