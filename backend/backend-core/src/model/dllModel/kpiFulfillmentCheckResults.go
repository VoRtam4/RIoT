package dllModel

import "time"

type KPIFulfillmentCheckResult struct {
	KPIDefinitionID uint32
	SDInstanceID    uint32
	Fulfilled       bool
	EventTime       time.Time
}
