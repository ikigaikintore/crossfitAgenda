package ports

import (
	"context"
	"github.com/ervitis/crossfitAgenda/domain"
	"net/http"
	"time"
)

type (
	ResourceManager interface {
		DownloadPicture() (string, error)
	}

	SourceReader interface {
		Read(ctx context.Context) (domain.RawProcessor, error)
		SetFile(path string)
	}

	SourceData interface {
		DownloadPicture() (string, error)
	}

	IAgendaService interface {
		GetCrossfitEvents() (ICalendar, error)
		UpdateEvents(ICalendar, domain.MonthWodExercises) error
	}

	IManager interface {
		SetConfigWithScopes(scopes ...string) error
		GetClient(ctx context.Context) *http.Client
	}

	IBook interface {
		GetEventID() string
		GetDay() int
		GetStartDate() time.Time
		GetEndDate() time.Time
	}

	ICalendar interface {
		GetID() string
		GetDaysBooked() map[int]IBook
		GetMonth() time.Month
	}
)
