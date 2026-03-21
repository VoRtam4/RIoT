package sharedModel

import "time"

type TimeSeriesRawRecord struct {
	EventTime           time.Time              `json:"eventTime"`
	SDInstanceUID       string                 `json:"sdInstanceUID"`
	SDTypeSpecification string                 `json:"sdTypeSpecification"`
	Fields              map[string]interface{} `json:"fields"`
	Tags                map[string]string      `json:"tags"`
}

type TimeSeriesKPIResultRecord struct {
	EventTime           time.Time         `json:"eventTime"`
	SDInstanceUID       string            `json:"sdInstanceUID"`
	SDTypeSpecification string            `json:"sdTypeSpecification"`
	KPIDefinitionID     uint32            `json:"kpiDefinitionID"`
	Fulfilled           bool              `json:"fulfilled"`
	Tags                map[string]string `json:"tags"`
}

type TimeSeriesType string

const (
	TimeSeriesTypeRaw       TimeSeriesType = "RAW"
	TimeSeriesTypeKPIResult TimeSeriesType = "KPI"
)

type TimeSeriesReprocessReadRequest struct {
	SDTypeSpecification string    `json:"sdTypeSpecification"`
	SDInstanceUIDs      []string  `json:"sdInstanceUIDs,omitempty"`
	To                  time.Time `json:"to"`
	Batch               int       `json:"batch"`
}

type TimeSeriesReprocessReadResponse struct {
	Data    []TimeSeriesDataPoint `json:"data"`
	HasMore bool                  `json:"hasMore"`
	Error   string                `json:"error,omitempty"`
}

type TimeSeriesReadRequest struct {
	Type                TimeSeriesType `json:"type"`
	SDInstanceUID       *string        `json:"sdInstanceUID,omitempty"`
	SDTypeSpecification *string        `json:"sdTypeSpecification,omitempty"`
	KPIDefinitionID     *uint32        `json:"kpiDefinitionID,omitempty"`
	From                *time.Time     `json:"from,omitempty"`
	To                  *time.Time     `json:"to,omitempty"`
	Normalize           *bool          `json:"normalize,omitempty"`
	AggregateMinutes    *int           `json:"aggregateMinutes,omitempty"`
	Limit               *int           `json:"limit,omitempty"`
	SortDesc            *bool          `json:"sortDesc,omitempty"`
	Batch               *int           `json:"batch,omitempty"`
}

type TimeSeriesReprocessRequest struct {
	SDTypeSpecification string    `json:"sdTypeSpecification"`
	SDInstanceUID       *string   `json:"sdInstanceUID,omitempty"`
	To                  time.Time `json:"to"`
	Batch               int       `json:"batch"`
}

type TimeSeriesDataPoint struct {
	Time time.Time              `json:"time"`
	Tags map[string]string      `json:"tags,omitempty"`
	Data map[string]interface{} `json:"data"`
}

type TimeSeriesReadResponse struct {
	Data    []TimeSeriesDataPoint `json:"data,omitempty"`
	HasMore bool                  `json:"hasMore,omitempty"`
	Error   string                `json:"error,omitempty"`
}
