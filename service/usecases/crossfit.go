package usecases

import (
	"context"
	"github.com/ervitis/crossfitAgenda/calendar"
	"github.com/ervitis/crossfitAgenda/credentials"
	"github.com/ervitis/crossfitAgenda/ocr"
	"github.com/ervitis/crossfitAgenda/ports"
	"github.com/ervitis/crossfitAgenda/service/domain"
	"sync"
)

type Crossfit interface {
	Book(ctx context.Context) error
	Status() domain.Status
}

type crossfit struct {
	rm  ports.ResourceManager
	mgr ports.IManager

	cache domain.Cache
}

func (c *crossfit) Book(ctx context.Context) error {
	status := domain.Working
	c.updateStatus(status)

	namePic, err := c.rm.DownloadPicture()
	if err != nil {
		status = domain.Failed
		c.updateStatus(status)
		return err
	}
	ocrClient := ocr.NewSourceReader(namePic)

	processor, err := ocrClient.Read()
	if err != nil {
		status = domain.Failed
		c.updateStatus(status)
		return err
	}

	monthWod := processor.Convert()
	credManager := credentials.New()
	_ = credManager.SetConfigWithScopes(calendar.Scope, calendar.EventsScope)
	svc, _ := calendar.New(ctx, credManager)
	events, err := svc.GetCrossfitEvents()
	if err != nil {
		status = domain.Failed
		c.updateStatus(status)
		return err
	}

	if err := svc.UpdateEvents(events, monthWod); err != nil {
		status = domain.Failed
		c.updateStatus(status)
		return err
	}

	status = domain.Finished
	c.updateStatus(status)
	return nil
}

func (c *crossfit) updateStatus(st domain.Status) {
	c.cache.Mtx.Lock()
	c.cache.Status = &st
	c.cache.Mtx.Unlock()
}

func (c *crossfit) Status() domain.Status {
	var st domain.Status
	c.cache.Mtx.Lock()
	st = *c.cache.Status
	c.cache.Mtx.Unlock()
	return st
}

func New(rm ports.ResourceManager, mgr ports.IManager) Crossfit {
	f := domain.Finished

	return &crossfit{
		rm,
		mgr,
		domain.Cache{
			Status: &f,
			Mtx:    sync.Mutex{},
		},
	}
}
