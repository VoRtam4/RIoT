package sharedModel

import "time"

type KPIFulfillmentCheckResultTupleISCMessage []KPIFulfillmentCheckResultISCMessage

type KPIFulfillmentCheckResultISCMessage struct {
	EventTime       time.Time `json:"eventTime"`
	SDInstanceUID   string    `json:"sdInstanceUID"`
	KPIDefinitionID uint32    `json:"kpiDefinitionID"`
	Fulfilled       bool      `json:"fulfilled"`
}

type KPIFulfillmentCheckRequestISCMessage struct {
	EventTime           time.Time `json:"eventTime"`
	SDInstanceUID       string    `json:"sdInstanceUID"`
	SDTypeSpecification string    `json:"sdTypeSpecification"`
	Parameters          any       `json:"parameters"`
}

type SDInstanceRegistrationRequestISCMessage struct {
	EventTime           time.Time `json:"eventTime"`
	SDInstanceUID       string    `json:"sdInstanceUID"`
	SDTypeSpecification string    `json:"sdTypeSpecification"`
}

type SDTypeConfigurationUpdateISCMessage []string

type SDInstanceInfo struct {
	SDInstanceUID   string `json:"sdInstanceUID"`
	ConfirmedByUser bool   `json:"confirmedByUser"`
}

type SDInstanceConfigurationUpdateISCMessage []SDInstanceInfo

type KPIConfigurationUpdateISCMessage map[string][]KPIDefinition

type MessageProcessingUnitConnectionNotification struct{}
