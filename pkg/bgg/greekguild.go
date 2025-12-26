package bgg

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Geeklist struct {
	Items []GeeklistItem `xml:"item"`
}

type GeeklistItem struct {
	ID         int    `xml:"id,attr"`
	ObjectID   int64  `xml:"objectid,attr"`
	ObjectName string `xml:"objectname,attr"`
	Username   string `xml:"username,attr"`
	Number     int    `xml:"number,attr"`
	Thumbs     int    `xml:"thumbs,attr"`
	Body       string `xml:"body"`
}

type Client struct {
	HTTP       *http.Client
	BaseURL    string
	MaxRetries int
	Token      string // optional Bearer token
}

func NewClient(token string) *Client {
	return &Client{
		HTTP:       &http.Client{Timeout: 30 * time.Second},
		BaseURL:    "https://boardgamegeek.com/xmlapi2",
		MaxRetries: 6,
		Token:      token,
	}
}

func (c *Client) GetGeeklistItems(ctx context.Context, listID int, comments bool) ([]GeeklistItem, error) {
	if c.HTTP == nil {
		c.HTTP = &http.Client{Timeout: 30 * time.Second}
	}
	if c.BaseURL == "" {
		c.BaseURL = "https://boardgamegeek.com/xmlapi2"
	}
	if c.MaxRetries <= 0 {
		c.MaxRetries = 6
	}

	u, err := url.Parse(c.BaseURL + "/geeklist/" + strconv.Itoa(listID))
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	q := u.Query()
	if comments {
		q.Set("comments", "1")
	} else {
		q.Set("comments", "0")
	}
	u.RawQuery = q.Encode()

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("new request: %w", err)
		}

		req.Header.Set("User-Agent", "bgg-geeklist-fetcher/1.3 (+https://example.com)")
		if c.Token != "" {
			req.Header.Set("Authorization", "Bearer "+c.Token)
		}

		resp, err := c.HTTP.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			sleepWithJitter(attempt)
			continue
		}

		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			sleepWithJitter(attempt)
			continue
		}

		switch resp.StatusCode {
		case http.StatusOK:
			var gl Geeklist
			if err := xml.Unmarshal(data, &gl); err != nil {
				return nil, fmt.Errorf("unmarshal xml: %w", err)
			}
			return gl.Items, nil // ✅ Correctly returns items here

		case http.StatusAccepted, http.StatusTooManyRequests:
			// BGG is still generating XML or rate limiting; retry
			sleepWithJitter(attempt)
			continue

		default:
			return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(data))
		}
	}

	return nil, fmt.Errorf("failed to fetch geeklist items after %d retries", c.MaxRetries)
}

func sleepWithJitter(attempt int) {
	base := 500 * time.Millisecond
	max := 8 * time.Second
	d := base << attempt
	if d > max {
		d = max
	}
	jitter := time.Duration(rand.Int63n(int64(d / 3)))
	time.Sleep(d/2 + jitter)
}
