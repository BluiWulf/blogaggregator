package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bluiwulf/blogaggregator/internal/database"
	"github.com/google/uuid"
)

func commandFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("%v command expects one argument, the feed URL", cmd.name)
	}
	feedURL := cmd.args[0]
	feed, err := s.db.GetFeedWithUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("feed at '%v' is not available: %v", feedURL, err)
	}

	followParams := database.CreateFeedFollowParams{
		ID: 		uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		UserID:		user.ID,
		FeedID:		feed.ID,
	}
	follow, err := s.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return fmt.Errorf("failed to follow RSS feed '%v': %v", feedURL, err)
	}

	fmt.Println()
	fmt.Println("Feed follow was successful")
	fmt.Println()
	fmt.Printf("Feed Name: %v\n", follow.FeedName)
	fmt.Printf("User Name: %v\n", follow.UserName)
	fmt.Println()

	return nil
}

func commandFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("%v command expects no argument", cmd.name)
	}
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get list of feed follows for user '%v': %v", user.Name, err)
	}

	fmt.Println()
	if len(follows) == 0 {
		fmt.Println("User is currently not following any RSS feeds")
	} else {
		for _, follow := range follows {
			fmt.Printf("%v\n", follow.FeedName)
		}
	}
	fmt.Println()

	return nil
}

func commandUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("%v command expects one argument, the feed URL", cmd.name)
	}
	feedURL := cmd.args[0]
	feed, err := s.db.GetFeedWithUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("feed at '%v' is not available: %v", feedURL, err)
	}

	deleteFollowParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), deleteFollowParams)
	if err != nil {
		return fmt.Errorf("failed to unfollow RSS feed '%v': %v", feedURL, err)
	}

	fmt.Println()
	fmt.Println("Feed unfollow was successful")
	fmt.Println()

	return nil
}
