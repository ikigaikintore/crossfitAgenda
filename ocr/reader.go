package ocr

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"github.com/ervitis/crossfitAgenda/domain"
	"github.com/ervitis/crossfitAgenda/ports"
	"os"
)

type (
	fileReader struct {
		path string
	}
)

func (fr fileReader) Read(ctx context.Context) (domain.RawProcessor, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = client.Close()
	}()

	file, err := os.Open(fr.path)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = file.Close()
	}()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		return nil, err
	}

	text, err := client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return nil, err
	}

	return domain.NewRawProcessor(text.Text), nil
}

func (fr fileReader) SetFile(path string) {
	fr.path = path
}

func NewSourceReader(fileName string) ports.SourceReader {
	return &fileReader{path: fileName}
}
