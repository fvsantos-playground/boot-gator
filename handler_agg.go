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

	savePosts(s, rss, feed)
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

func savePosts(s *state, rss *RSSFeed, feed database.Feed) error {
	for _, post := range rss.Channel.Item {
		p, perr := time.Parse(time.RFC1123Z, post.PubDate)
		newPost, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			Title: post.Title,
			Url:   post.Link,
			Description: sql.NullString{
				String: post.Description,
				Valid:  len(post.Description) > 0},
			PublishedAt: sql.NullTime{
				Time:  p,
				Valid: perr == nil,
			},
			FeedID:    feed.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			fmt.Printf("Error saving post: %v\n", err)
		} else {
			fmt.Printf("Saved post: %s\n", newPost.Title)
		}
	}
	return nil
}

func printRSS(rss *RSSFeed) {
	fmt.Printf("Title: %s\n", rss.Channel.Title)
	fmt.Printf("Description: %s\n", rss.Channel.Description)
	fmt.Printf("Items: %d\n", len(rss.Channel.Item))
	for _, item := range rss.Channel.Item {
		fmt.Printf("  Title: %s\n", item.Title)
		fmt.Printf("  Link: %s\n", item.Link)
		fmt.Printf("  Description: %s\n", item.Description)
		fmt.Printf("  Published At: %s\n", item.PubDate)
		fmt.Println()
	}
	fmt.Print("\n====================================\n\n")
}
