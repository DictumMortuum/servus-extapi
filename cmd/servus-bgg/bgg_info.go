package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/DictumMortuum/servus-extapi/pkg/bgg"
	"github.com/DictumMortuum/servus-extapi/pkg/minio"
	minio_lib "github.com/minio/minio-go/v7"
)

func getBoardgameInfo(bucket string, id int64) ([]byte, error) {
	time.Sleep(200 * time.Millisecond)

	ctx := context.Background()
	name := fmt.Sprintf("%d.json", id)

	client, err := minio.Init()
	if err != nil {
		return nil, err
	}

	shouldUpload, err := uploadIfMissingOrStale(ctx, client, bucket, name)
	if err != nil {
		return nil, err
	}

	if !shouldUpload {
		obj, err := client.GetObject(ctx, bucket, name, minio_lib.GetObjectOptions{})
		if err != nil {
			return nil, err
		}
		defer obj.Close()

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, obj); err != nil {
			return nil, err
		}

		data := buf.Bytes() // []byte containing the JSON
		return data, nil
	}

	raw, err := bgg.BoardgameInfoNoDB(id)
	if err != nil {
		return nil, err
	}

	// 2) Upload (overwrite) the object
	reader := bytes.NewReader(raw)
	_, err = client.PutObject(ctx, bucket, name, reader, int64(len(raw)), minio_lib.PutObjectOptions{
		ContentType: "application/json",
	})
	if err != nil {
		return nil, err
	}

	return raw, nil
}
