package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fvsantos-playground/boot-gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	timeBetweenRequests := time.Duration(time.Second * 45)

	if len(cmd.Args) > 0 {
		var err error
		timeBetweenRequests, err = time.ParseDuration(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return err
	}

	log.Printf("Found a feed to fetch: %s", feed.Url)
	scrapeFeed(s.db, feed)

	// rss, err := fetchFeed(context.Background(), feed.Url)
	// if err != nil {
	// 	log.Println("Couldn't fetch feed", err)
	// 	return err
	// }

	// savePosts(s, rss, feed)
	// printRSS(rss)

	return nil
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		savePosts(db, feed, item)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

func savePosts(db *database.Queries, feed database.Feed, post RSSItem) {
	publishedAt := sql.NullTime{}
	if t, err := time.Parse(time.RFC1123Z, post.PubDate); err == nil {
		publishedAt = sql.NullTime{
			Time:  t,
			Valid: true,
		}
	}

	_, err := db.CreatePost(context.Background(), database.CreatePostParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		Title:     post.Title,
		Description: sql.NullString{
			String: post.Description,
			Valid:  true,
		},
		Url:         post.Link,
		PublishedAt: publishedAt,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return
		}
		log.Printf("Couldn't create post: %v", err)
		return
	}
}

// func printRSS(rss *RSSFeed) {
// 	fmt.Printf("Title: %s\n", rss.Channel.Title)
// 	fmt.Printf("Description: %s\n", rss.Channel.Description)
// 	fmt.Printf("Items: %d\n", len(rss.Channel.Item))
// 	for _, item := range rss.Channel.Item {
// 		fmt.Printf("  Title: %s\n", item.Title)
// 		fmt.Printf("  Link: %s\n", item.Link)
// 		fmt.Printf("  Description: %s\n", item.Description)
// 		fmt.Printf("  Published At: %s\n", item.PubDate)
// 		fmt.Println()
// 	}
// 	fmt.Print("\n====================================\n\n")
// }
