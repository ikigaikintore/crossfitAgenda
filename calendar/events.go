package calendar

import (
	"context"
	"fmt"
	"github.com/ervitis/crossfitAgenda/ports"
	"log"
	"time"

	"github.com/ervitis/crossfitAgenda/credentials"
	"github.com/ervitis/crossfitAgenda/domain"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type (
	agendaService struct {
		calendar *calendar.Service
	}

	book struct {
		EventID   string
		Day       int
		StartDate time.Time
		EndDate   time.Time
	}

	Calendar struct {
		ID         string
		DaysBooked map[int]ports.IBook
		Month      time.Month
	}
)

func (b book) GetEventID() string {
	return b.EventID
}

func (b book) GetDay() int {
	return b.Day
}

func (b book) GetStartDate() time.Time {
	return b.StartDate
}

func (b book) GetEndDate() time.Time {
	return b.EndDate
}

func (c Calendar) GetID() string {
	return c.ID
}

func (c Calendar) GetDaysBooked() map[int]ports.IBook {
	return c.DaysBooked
}

func (c Calendar) GetMonth() time.Month {
	return c.Month
}

func New(ctx context.Context, credManager *credentials.Manager) (ports.IAgendaService, error) {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(credManager.GetClient(ctx)))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	return &agendaService{
		calendar: srv,
	}, nil
}

func (w *agendaService) UpdateEvents(cal ports.ICalendar, wods domain.MonthWodExercises) error {
	for day, wod := range wods {
		book, exists := cal.GetDaysBooked()[day]
		if !exists {
			continue
		}

		if book.GetStartDate().Month() != wod.Month() {
			continue
		}

		if wod.ExerciseName().String() == "" {
			continue
		}

		if _, err := w.calendar.Events.Update(cal.GetID(), book.GetEventID(), &calendar.Event{
			Description: wod.ExerciseName().String(),
			Summary:     fmt.Sprintf("Crossfit class: %s", wod.ExerciseName().String()),
			End:         &calendar.EventDateTime{DateTime: book.GetEndDate().Format(time.RFC3339)},
			Start:       &calendar.EventDateTime{DateTime: book.GetStartDate().Format(time.RFC3339)},
		}).Do(); err != nil {
			return err
		}
	}
	return nil
}

func (w *agendaService) getLocation() *time.Location {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	return loc
}

func (w *agendaService) GetCrossfitEvents() (ports.ICalendar, error) {
	now := time.Now().In(w.getLocation())

	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, w.getLocation())
	lastOfMonth := firstOfMonth.AddDate(0, 2, -1)

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
		TimeMax(lastOfMonth.AddDate(0, 0, 1).Format(time.RFC3339)).
		Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	if len(events.Items) == 0 {
		log.Println("No upcoming events found.")
		return &Calendar{}, nil
	}

	myCalendar := &Calendar{Month: firstOfMonth.Month(), DaysBooked: make(map[int]ports.IBook), ID: primaryCalendar.Id}
	for _, item := range events.Items {
		startDateEvent, _ := time.Parse(time.RFC3339, item.Start.DateTime)
		endDateEvent, _ := time.Parse(time.RFC3339, item.End.DateTime)
		dayEvent := startDateEvent.Day()
		myCalendar.DaysBooked[dayEvent] = &book{
			EventID:   item.Id,
			Day:       dayEvent,
			StartDate: startDateEvent,
			EndDate:   endDateEvent,
		}
	}
	return myCalendar, nil
}
