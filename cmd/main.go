package main

import (
	"context"
	"fmt"
	"github.com/ervitis/crossfitAgenda/credentials"
	"github.com/ervitis/crossfitAgenda/crossfit_events"
	"github.com/ervitis/crossfitAgenda/ocr"
	"github.com/ervitis/crossfitAgenda/source_data"
	"google.golang.org/api/calendar/v3"
	"log"
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

	fmt.Println(processor.Convert())

	credManager := credentials.New()
	_ = credManager.SetConfigWithScopes(calendar.CalendarScope, calendar.CalendarEventsScope)
	calService, _ := crossfit_events.New(context.Background(), credManager)
	events, err := calService.GetCrossfitEvents()
	if err != nil {
		log.Printf("error getting events: %s\n", err)
	}

	fmt.Println(events)
}
