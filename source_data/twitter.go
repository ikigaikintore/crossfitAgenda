package source_data

import (
	"fmt"
	"github.com/ervitis/crossfitAgenda/ports"
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

func NewTwitterClient() ports.SourceData {
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
