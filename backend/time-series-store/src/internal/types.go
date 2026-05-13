/**
 * @file types.go
 * @brief Datové typy a konfigurační struktury interní části modulu Time Series Store.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_time_series_store
 */
package internal

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type TimeSeriesStoreEnvironment struct {
	InfluxToken  string
	InfluxUrl    string
	InfluxOrg    string
	InfluxBucket string
	AmqpURLValue string
}

type Influx2Client struct {
	endpoint     string
	organization string
	bucket       string
	client       influxdb2.Client
	writeApi     api.WriteAPI
	queryApi     api.QueryAPI
}

type streamState struct {
	totalSent    int
	pointsBatch  []sharedModel.TimeSeriesDataPoint
	cursorPassed bool
	hadAnyData   bool
	needsTerminal bool
}

type aggregateSeriesState struct {
	Key  string
	Tags map[string]string

	LastBool     bool
	LastTime     time.Time
	HasLastValue bool

	TrueWindowTime  time.Duration
	FalseWindowTime time.Duration

	HasFullWindow bool
}

type aggregateStreamState struct {
	totalSent   int
	pointsBatch []sharedModel.TimeSeriesDataPoint

	cursorPassed bool
	needsTerminal bool

	series      map[string]*aggregateSeriesState
	seriesOrder []string

	windowStart time.Time
	windowSize  time.Duration
}
