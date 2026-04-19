package dllModel

import "time"

type KPIFulfillmentCheckResult struct {
	SDTypeID        uint32
	KPIDefinitionID uint32
	SDInstanceID    uint32
	Fulfilled       bool
	EventTime       time.Time
}

type KPIFulfillmentCheckResultRequest struct {
	KPIDefinitionID uint32
	SDInstanceID    uint32
}
