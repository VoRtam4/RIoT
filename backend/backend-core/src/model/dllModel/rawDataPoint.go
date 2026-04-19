package dllModel

import "time"

type RawDataPoint struct {
	SDTypeID     uint32
	SDInstanceID uint32
	EventTime    time.Time
	Payload      []byte
}
