package main

import (
	"context"
	"log"

	"github.com/ervitis/crossfitAgenda/calendar"
	"github.com/ervitis/crossfitAgenda/credentials"
	"github.com/ervitis/crossfitAgenda/ocr"
	"github.com/ervitis/crossfitAgenda/source_data"
)

func main() {
	resourceManager := source_data.NewResourceManager(
		source_data.WithSourceDataClient(source_data.NewTwitterClient()),
	)

	name, err := resourceManager.DownloadPicture()
	if err != nil {
		log.Printf("error happened: %s\n", err)
	}

	ocrClient := ocr.NewSourceReader(name)
	processor, err := ocrClient.Read()
	if err != nil {
		log.Printf("error in ocr client: %s\n", err)
	}

	monthWod := processor.Convert()

	credManager := credentials.New()
	_ = credManager.SetConfigWithScopes(calendar.Scope, calendar.EventsScope)
	calService, _ := calendar.New(context.Background(), credManager)
	events, err := calService.GetCrossfitEvents()
	if err != nil {
		log.Printf("error getting events: %s\n", err)
	}

	if err := calService.UpdateEvents(events, monthWod); err != nil {
		log.Fatal(err)
	}
}
