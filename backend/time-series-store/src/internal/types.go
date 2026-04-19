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
	writeApi     api.WriteAPIBlocking
	queryApi     api.QueryAPI
}

type streamState struct {
	totalSent    int
	pointsBatch  []sharedModel.TimeSeriesDataPoint
	cursorPassed bool
	hadAnyData   bool
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

	series      map[string]*aggregateSeriesState
	seriesOrder []string

	windowStart time.Time
	windowSize  time.Duration
}
