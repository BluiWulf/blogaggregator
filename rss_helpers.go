package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/bluiwulf/blogaggregator/internal/database"
)

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	if len(feedURL) == 0 {
		return nil, fmt.Errorf("rss feed URL must be provided")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	feed := RSSFeed{}
	err = xml.Unmarshal(data, &feed)

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func (c *Client) scrapeFeed(s *state) (*RSSFeed, error) {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return nil, err
	}
	
	markFeedParams := database.MarkFeedFetchedParams{
		ID:			nextFeed.ID,
		UpdatedAt:	time.Now(),
	}
	err = s.db.MarkFeedFetched(context.Background(), markFeedParams)
	if err != nil {
		return nil, err
	}

	feed, err := c.fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return nil, err
	}
	fmt.Println()
	fmt.Printf(" * %v *\n", feed.Channel.Title)
	fmt.Println()
	for _, item := range feed.Channel.Item {
		fmt.Printf("%v\n", item.Title)
	}
	fmt.Println()

	return feed, nil
}
