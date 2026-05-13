/**
 * @file model.go
 * @brief Datový model interních událostí publikovaných Backend Core.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package events

import "time"

type EventType string

const (
	SDInstanceRegisteredEventType  EventType = "sd_instance_registered"
	RawDataPointReceivedEventType  EventType = "raw_data_point"
	KPIFulfillmentCheckedEventType EventType = "kpi_fulfillment_checked"
	TimeSeriesExportUpdatedEventType EventType = "time_series_export_updated"
)

type Event struct {
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Payload   any       `json:"payload"`
}
