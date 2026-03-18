package events

type EventType string

const (
	SDInstanceRegisteredEventType  EventType = "sd_instance_registered"
	KPIFulfillmentCheckedEventType EventType = "kpi_fulfillment_checked"
)
