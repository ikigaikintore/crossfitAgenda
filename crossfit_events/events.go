package crossfit_events

import (
	"context"
	"fmt"
	"github.com/ervitis/crossfitAgenda/credentials"
	"github.com/ervitis/crossfitAgenda/domain"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log"
	"time"
)

type (
	agendaService struct {
		calendar *calendar.Service
	}

	book struct {
		EventID string
		Day     int
		Hour    int
		Minute  int
	}

	Calendar struct {
		ID         string
		DaysBooked map[int]*book
		Month      time.Month
	}

	IAgendaService interface {
		GetCrossfitEvents() (*Calendar, error)
		UpdateEvents(*Calendar, domain.MonthWodExercises) error
	}
)

func New(ctx context.Context, credManager *credentials.Manager) (IAgendaService, error) {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(credManager.GetClient()))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	return &agendaService{
		calendar: srv,
	}, nil
}

func (w *agendaService) UpdateEvents(cal *Calendar, wods domain.MonthWodExercises) error {
	for _, v := range cal.DaysBooked {
		for _, wod := range wods {
			if v.Day != wod.Day() {
				continue
			}

			w.calendar.Events.Update(cal.ID, v.EventID, &calendar.Event{
				Description: wod.ExerciseName().String(),
				Summary:     fmt.Sprintf("Crossfit class: %s", wod.ExerciseName().String()),
			})
		}
	}
	return nil
}

func (w *agendaService) getLocation() *time.Location {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	return loc
}

func (w *agendaService) GetCrossfitEvents() (*Calendar, error) {
	now := time.Now().In(w.getLocation())

	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, w.getLocation())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	primaryCalendar, err := w.calendar.Calendars.Get("primary").Do()
	if err != nil {
		return nil, err
	}

	events, err := w.calendar.
		Events.
		List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		Q("Class").
		TimeMin(firstOfMonth.Format(time.RFC3339)).
		TimeMax(lastOfMonth.Format(time.RFC3339)).
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	if len(events.Items) == 0 {
		log.Println("No upcoming events found.")
		return &Calendar{}, nil
	}

	myCalendar := &Calendar{Month: firstOfMonth.Month(), DaysBooked: make(map[int]*book), ID: primaryCalendar.Id}
	for _, item := range events.Items {
		startDateEvent, _ := time.Parse(time.RFC3339, item.Start.Date)
		dayEvent := startDateEvent.Day()
		myCalendar.DaysBooked[dayEvent].EventID = item.Id
		myCalendar.DaysBooked[dayEvent].Day = dayEvent
		myCalendar.DaysBooked[dayEvent].Hour = startDateEvent.Hour()
		myCalendar.DaysBooked[dayEvent].Minute = startDateEvent.Minute()
	}
	return myCalendar, nil
}
