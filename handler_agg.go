package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fvsantos-playground/boot-gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	time_between_reqs := time.Duration(time.Second * 45)

	if len(cmd.Args) > 0 {
		var err error
		time_between_reqs, err = time.ParseDuration(cmd.Args[0])
		fmt.Printf("... %v", err)
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	rss, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	printRSS(rss)

	s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: feed.ID,
	})

	return nil
}

func printRSS(rss *RSSFeed) {
	fmt.Printf("Title: %s\n", rss.Channel.Title)
	fmt.Printf("Description: %s\n", rss.Channel.Description)
	fmt.Printf("Items: %d\n", len(rss.Channel.Item))
	for _, item := range rss.Channel.Item {
		fmt.Printf("Item: %s\n", item.Title)
	}
	fmt.Print("\n====================================\n\n")
}
