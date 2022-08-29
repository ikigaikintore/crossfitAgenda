package event_handler

import (
	"sync"
	"time"
)

type (
	Metadata struct {
		Timestamp time.Time
		ID        string
	}

	Message struct {
		Payload  []byte
		Metadata Metadata
	}

	gobbit struct {
		subs  map[string]chan *Message
		mtx   sync.Mutex
		order uint64
	}
)
