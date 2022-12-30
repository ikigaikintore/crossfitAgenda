package source_data

import (
	"log"

	"github.com/ervitis/crossfitAgenda/ports"
)

type (
	source struct {
		sourceClient ports.SourceData
	}

	SourceOption func(*source)

	dumbClientSourceData struct {
		log *log.Logger
	}
)

func defaultClient() ports.SourceData {
	return &dumbClientSourceData{log: log.Default()}
}

func (d dumbClientSourceData) DownloadPicture() (string, error) {
	d.log.Println("nothing to do here")
	return "", nil
}

func defaultSourceOption() *source {
	return &source{sourceClient: defaultClient()}
}

func (s source) DownloadPicture() (string, error) {
	return s.sourceClient.DownloadPicture()
}

func WithSourceDataClient(data ports.SourceData) SourceOption {
	return func(s *source) {
		s.sourceClient = data
	}
}

func NewResourceManager(opts ...SourceOption) ports.ResourceManager {
	sd := defaultSourceOption()

	for _, opt := range opts {
		opt(sd)
	}

	return sd
}
