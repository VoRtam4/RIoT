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

type SDTypeRegistrationRequestISCMessage struct {
	SDTypeSpecification string        `json:"sdTypeSpecification"`
	Parameters          []SDParameter `json:"parameters"`
}

type SDParameter struct {
	Denotation string `json:"denotation"`
	Type       string `json:"type"`
}

type KPIReprocessRequestISCMessage struct {
	KPIDefinitionID     uint32    `json:"kpiDefinitionID"`
	SDTypeSpecification string    `json:"sdTypeSpecification"`
	From                time.Time `json:"from"`
}

type KPIDeleteResultsRequestISCMessage struct {
	KPIDefinitionID uint32 `json:"kpiDefinitionID"`
}
