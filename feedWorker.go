package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jeffroutledge/glowing-octo-system/internal/database"
)

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
			go getFeed(f)
			cfg.DB.MarkFeedFetched(context.Background(), f.ID)
			wg.Done()
		}

		wg.Wait()

		time.Sleep(60 * time.Second)
	}
}

func getFeed(f database.Feed) Rss {
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

	fmt.Println(result.Channel.Item[0].Title)

	return result
}
