package source_data

import "log"

type (
	ResourceManager interface {
		DownloadPicture() (string, error)
	}

	SourceData interface {
		DownloadPicture() (string, error)
	}

	source struct {
		sourceClient SourceData
	}

	SourceOption func(*source)

	dumbClientSourceData struct {
		log *log.Logger
	}
)

func defaultClient() SourceData {
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

func WithSourceDataClient(data SourceData) SourceOption {
	return func(s *source) {
		s.sourceClient = data
	}
}

func NewResourceManager(opts ...SourceOption) ResourceManager {
	sd := defaultSourceOption()

	for _, opt := range opts {
		opt(sd)
	}

	return sd
}
