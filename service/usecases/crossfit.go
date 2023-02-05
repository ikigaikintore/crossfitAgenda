package usecases

import (
	"context"
	"github.com/ervitis/crossfitAgenda/calendar"
	"github.com/ervitis/crossfitAgenda/credentials"
	"github.com/ervitis/crossfitAgenda/ocr"
	"github.com/ervitis/crossfitAgenda/ports"
	"github.com/ervitis/crossfitAgenda/service/domain"
	"sync"
	"time"
)

type AgendaCrossfit interface {
	Book(ctx context.Context) error
	Status() domain.ProcessStatus
	GetCredentialsUrl() string
	SaveToken(ctx context.Context, token string) error
}

type crossfit struct {
	rm  ports.ResourceManager
	mgr ports.IManager

	cred *credentials.Manager

	cache domain.Cache
}

func (c *crossfit) GetCredentialsUrl() string {
	return c.cred.GetUrlLogin()
}

func (c *crossfit) SaveToken(ctx context.Context, token string) error {
	return c.cred.SaveToken(ctx, token)
}

func (c *crossfit) Book(ctx context.Context) error {
	if c.cache.Process.IsRunning() {
		return nil
	}

	if c.cred.IsTokenExpired() {
		c.updateStatus(domain.Failed)
		return domain.ErrTokenCredentialsExpired
	}

	wg := sync.WaitGroup{}
	type result struct {
		finished bool
		err      error
	}
	processEnd := make(chan result, 1)

	bookProcess := func(processEnd chan<- result) {
		c.updateStatus(domain.Working)

		namePic, err := c.rm.DownloadPicture()
		if err != nil {
			c.updateStatus(domain.Failed)
			processEnd <- result{finished: true, err: err}
		}
		ocrClient := ocr.NewSourceReader(namePic)

		processor, err := ocrClient.Read(ctx)
		if err != nil {
			c.updateStatus(domain.Failed)
			processEnd <- result{finished: true, err: err}
		}

		monthWod := processor.Convert()

		svc, _ := calendar.New(ctx, c.cred)
		events, err := svc.GetCrossfitEvents()
		if err != nil {
			c.updateStatus(domain.Failed)
			processEnd <- result{finished: true, err: err}
		}

		if err = svc.UpdateEvents(events, monthWod); err != nil {
			c.updateStatus(domain.Failed)
			processEnd <- result{finished: true, err: err}
		}
		processEnd <- result{finished: true}
	}

	go bookProcess(processEnd)

	var errResult error
	wg.Add(1)

	go func() {
		ret := <-processEnd
		errResult = ret.err
		wg.Done()
	}()
	wg.Wait()
	status := domain.Finished
	if errResult != nil {
		status = domain.Failed
	}
	c.updateStatus(status)
	return errResult
}

func (c *crossfit) updateStatus(st domain.Status) {
	c.cache.Mtx.Lock()
	c.cache.Process.Status = &st
	c.cache.Mtx.Unlock()
}

func (c *crossfit) Status() domain.ProcessStatus {
	c.cache.Mtx.Lock()
	st := c.cache.Process
	c.cache.Mtx.Unlock()
	return st
}

func New(rm ports.ResourceManager, mgr ports.IManager) AgendaCrossfit {
	f := domain.Finished

	cred := credentials.New()
	_ = cred.SetConfigWithScopes(calendar.Scope, calendar.EventsScope)

	return &crossfit{
		rm,
		mgr,
		cred,
		domain.Cache{
			Process: domain.ProcessStatus{
				ID:       1,
				Date:     time.Now().UTC(),
				Detail:   "starting",
				Status:   &f,
				Complete: false,
			},
			Mtx: sync.Mutex{},
		},
	}
}
