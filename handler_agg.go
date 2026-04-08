package main

import (
	"context"
	"fmt"
	"log"
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
		fmt.Printf("Found post: %s\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

// func savePosts(s *state, rss *RSSFeed, feed database.Feed) error {
// 	for _, post := range rss.Channel.Item {
// 		p, perr := time.Parse(time.RFC1123Z, post.PubDate)
// 		newPost, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
// 			Title: post.Title,
// 			Url:   post.Link,
// 			Description: sql.NullString{
// 				String: post.Description,
// 				Valid:  len(post.Description) > 0},
// 			PublishedAt: sql.NullTime{
// 				Time:  p,
// 				Valid: perr == nil,
// 			},
// 			FeedID:    feed.ID,
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		})
// 		if err != nil {
// 			fmt.Printf("Error saving post: %v\n", err)
// 		} else {
// 			fmt.Printf("Saved post: %s\n", newPost.Title)
// 		}
// 	}
// 	return nil
// }

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
