package domain

import (
	"sync"
	"time"
)

type Status int

const (
	Working Status = iota + 1
	Finished
	Failed
)

func (st Status) IsComplete() bool {
	return st != Working
}

func (st Status) IsRunning() bool {
	return st == Working
}

func (ps ProcessStatus) IsComplete() bool {
	return ps.Status.IsComplete()
}

func (ps ProcessStatus) IsRunning() bool {
	return ps.Status.IsRunning()
}

type ProcessStatus struct {
	ID       int64
	Date     time.Time
	Detail   string
	Status   *Status
	Complete bool
}

type Cache struct {
	Mtx     sync.Mutex
	Process ProcessStatus
}
