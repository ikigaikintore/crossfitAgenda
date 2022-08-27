package source_data

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

type (
	twitter struct {
		client *twitterscraper.Scraper
		uri    string
	}
)

func NewTwitterClient() SourceData {
	return &twitter{
		client: twitterscraper.New(),
		uri:    "Haleo_DKY",
	}
}

func (tw twitter) DownloadPicture() (string, error) {
	tweets, _, err := tw.client.FetchTweets(tw.uri, 1, "")
	if err != nil {
		return "", fmt.Errorf("downloading profile: %w", err)
	}

	if len(tweets) == 0 || len(tweets[0].Photos) == 0 {
		log.Println("no pictures to download")
		return "", nil
	}

	client := http.Client{Timeout: 20 * time.Second}

	dateTweet := tweets[0].TimeParsed
	now := time.Now()

	if dateTweet.Month() != now.Month() {
		return "", fmt.Errorf("month of the tweet is not actual, wait a little")
	}

	resp, err := client.Get(tweets[0].Photos[0])
	if err != nil {
		return "", fmt.Errorf("downloading photo: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	tempFile, err := os.CreateTemp("", "crossfitMonth.jpg")
	if err != nil {
		return "", fmt.Errorf("creating file: %w", err)
	}

	// send an event the file was copied successful
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("copy body file: %w", err)
	}

	defer func() {
		_ = tempFile.Close()
	}()

	return tempFile.Name(), nil
}
