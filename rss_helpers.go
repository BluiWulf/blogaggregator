package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bluiwulf/blogaggregator/internal/database"
	"github.com/google/uuid"
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
		var pubAt sql.NullTime
		if parsedT, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			pubAt = sql.NullTime{
				Time: parsedT,
				Valid: true,
			}
		} else if parsedT, err := time.Parse(time.RFC1123, item.PubDate); err == nil {
			pubAt = sql.NullTime{
				Time: parsedT,
				Valid: true,
			}
		}

		postParams := database.CreatePostParams{
			ID:				uuid.New(),
			CreatedAt:		time.Now(),
			UpdatedAt:		time.Now(),
			Title:			item.Title,
			Url:			item.Link,
			Description:	item.Description,
			PublishedAt:	pubAt.Time,
			FeedID:			nextFeed.ID,
		}
		_, err := s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			if !strings.Contains(err.Error(), "no_dupl_posts") {
				return nil, err
			}
		}
	}
	fmt.Println()

	return feed, nil
}
