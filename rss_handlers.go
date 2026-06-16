package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bluiwulf/blogaggregator/internal/database"
	"github.com/google/uuid"
)

type Client struct {
	httpClient http.Client
}

type RSSFeed struct {
	Channel		RSSChannel	`xml:"channel"`
}

type RSSChannel struct {
	Title		string		`xml:"title"`
	Link		string		`xml:"link"`
	Description	string		`xml:"description"`
	Item		[]RSSItem	`xml:"item"`
}

type RSSItem struct {
	Title		string		`xml:"title"`
	Link		string		`xml:"link"`
	Description	string		`xml:"description"`
	PubDate		string		`xml:"pubDate"`
}

func commandAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("%v command expects one argument, the time between requests", cmd.name)
	}
	reqInterval, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("the time between requests provided is invalid: %v", err)
	}
	ticker := time.NewTicker(reqInterval)
	rssClient := NewClient(5 * time.Second)

	fmt.Println()
	fmt.Printf("Collecting feeds every %v", reqInterval)
	fmt.Println()
	for ; ; <- ticker.C {
		fmt.Println("Scraping next feed...")
		_, err = rssClient.scrapeFeed(s)
		if err != nil {
			return fmt.Errorf("failed to scrape rss feed: %v", err)
		}
		fmt.Println()
	}

	return nil
}

func commandAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("%v command expects two arguments, the feed name and URL", cmd.name)
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]
	feed, err := s.db.GetFeedWithUrl(context.Background(), feedURL)
	if err == nil {
		if feed.UserID == user.ID {
			return fmt.Errorf("'%v' already created the RSS feed '%v'", user.Name, feed.Name)
		}
	}

	feedParams := database.CreateFeedParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		feedName,
		Url:		feedURL,
		UserID:		user.ID,
	}
	feed, err = s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("failed to create RSS feed '%v': %v", feedName, err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		UserID:		user.ID,
		FeedID:		feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("failed to follow RSS feed '%v': %v", feedURL, err)
	}

	fmt.Println()
	fmt.Println("Feed has been created")
	fmt.Println()
	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("Url: %v\n", feed.Url)
	fmt.Printf("UserID: %v\n", feed.UserID)
	fmt.Println()

	return nil
}

func commandFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("%v command expects no argument", cmd.name)
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get list of feeds: %v", err)
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserWithId(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user '%v' information: %v", feed.UserID, err)
		}
		fmt.Println()
		fmt.Printf("Name: %v\n", feed.Name)
		fmt.Printf("Url: %v\n", feed.Url)
		fmt.Printf("User: %v\n", user.Name)
	}
	fmt.Println()

	return nil
}

func commandBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("%v command expects one optional argument max, the post limit", cmd.name)
	}
	limit := 2
	if len(cmd.args) == 1 {
		optLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("limit passed is not valid: %v", err)
		}
		limit = optLimit
	}

	userPostsParams:= database.GetPostsForUserParams{
		UserID:	user.ID,
		Limit:	int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), userPostsParams)
	if err != nil {
		return fmt.Errorf("failed to get list of user posts: %v", err)
	}
	for _, post := range posts {
		fmt.Println("---")
		fmt.Printf("Title: %v\n", post.Title)
		fmt.Printf("Link: %v\n", post.Url)
		fmt.Printf("Published: %v\n", post.PublishedAt)
		fmt.Println()
		fmt.Printf("Description: %v\n", post.Description)
		fmt.Println("---")
	}

	return nil
}
