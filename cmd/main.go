package main

import (
	"fmt"
	"github.com/ervitis/crossfitAgenda/ocr"
	"github.com/ervitis/crossfitAgenda/source_data"
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
}
