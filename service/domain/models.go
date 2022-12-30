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

type ProcessStatus struct {
	ID       int64
	Date     time.Time
	Detail   string
	Status   Status
	Complete bool
}

type Cache struct {
	Mtx    sync.Mutex
	Status *Status
}
