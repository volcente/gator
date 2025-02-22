package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

const feedURL = "https://www.wagslane.dev/index.xml"

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	feedRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to construct request: %w", err)
	}

	feedRequest.Header.Set("User-Agent", "gator")
	res, err := client.Do(feedRequest)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request: %w", err)
	}
	defer res.Body.Close()

	rssFeed := RSSFeed{}
	if err = xml.NewDecoder(res.Body).Decode(&rssFeed); err != nil {
		return nil, fmt.Errorf("Failed to parse response: %w", err)
	}

	decodeUnescapedHTML(&rssFeed)
	return &rssFeed, nil
}

func decodeUnescapedHTML(rawFeed *RSSFeed) {
	rawFeed.Channel.Title = html.UnescapeString(rawFeed.Channel.Title)
	rawFeed.Channel.Description = html.UnescapeString(rawFeed.Channel.Description)
	// REMEMBER range loop creates a copy for each iterated item, modifying item doesn't change the source object! always commit changes via obj.item[idx] = modifiedItem
	for idx, item := range rawFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rawFeed.Channel.Item[idx] = item
	}
}
