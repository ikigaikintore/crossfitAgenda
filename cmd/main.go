package main

import (
	"github.com/ervitis/crossfitAgenda/source_data"
	"log"
)

func main() {
	resourceManager := source_data.NewResourceManager(
		source_data.WithSourceDataClient(source_data.NewTwitterClient()),
	)

	if err := resourceManager.DownloadPicture(); err != nil {
		log.Printf("error happened: %s\n", err)
	}
}
