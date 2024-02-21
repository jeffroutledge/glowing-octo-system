package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jeffroutledge/glowing-octo-system/internal/database"
)

var ErrDuplicateUrl = errors.New("pq: duplicate key value violates unique constraint \"posts_url_key\"")

type Feed struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	FeedID        uuid.UUID
	UserID        uuid.UUID
	LastFetchedAt time.Time
}

// Rss was generated 2024-02-20 21:27:31 by https://xml-to-go.github.io/ in Ukraine.
type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (cfg *apiConfig) fetchFeeds() {
	for {
		var wg sync.WaitGroup

		feeds, err := cfg.DB.GetNextFeedsToFetch(context.Background(), 10)
		if err != nil {
			log.Default().Println(err.Error())
		}

		for _, f := range feeds {
			wg.Add(1)
			rss := make(chan Rss)
			go getRss(f, rss)
			cfg.addPost(<-rss, f.ID)
			cfg.DB.MarkFeedFetched(context.Background(), f.ID)
			wg.Done()
		}

		wg.Wait()

		time.Sleep(60 * time.Second)
	}
}

func getRss(f database.Feed, rss chan (Rss)) {
	url := f.Url.String
	resp, err := http.Get(url)
	if err != nil {
		log.Default().Println(err.Error())
	}

	decoder := xml.NewDecoder(resp.Body)
	result := Rss{}
	err = decoder.Decode(&result)
	if err != nil {
		log.Default().Println(err.Error())
	}

	rss <- result
}

func (cfg *apiConfig) addPost(rss Rss, feedID uuid.UUID) ([]database.Post, error) {
	posts := make([]database.Post, 0)

	for _, post := range rss.Channel.Item {
		uniqueID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		pubDate, err := time.Parse(time.RFC1123Z, post.PubDate)
		if err != nil {
			log.Default().Println(err.Error())
		}
		newPost := database.CreatePostParams{
			ID:          uniqueID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: post.Title, Valid: true},
			Url:         sql.NullString{String: post.Link, Valid: true},
			Description: sql.NullString{String: post.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: pubDate, Valid: true},
			FeedID:      uuid.NullUUID{UUID: feedID, Valid: true},
		}
		p, err := cfg.DB.CreatePost(context.Background(), newPost)
		if err == ErrDuplicateUrl {
			fmt.Printf("Already have blog post with URL: %s\n", post.Link)
			continue
		} else if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}
