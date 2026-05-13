/**
 * @file client.go
 * @brief Inicializace klienta pro komunikaci s databází InfluxDB.
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
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewInflux2Client(endpoint string, token string, organization string, bucket string) Influx2Client {
	client := influxdb2.NewClientWithOptions(
		endpoint,
		token,
		influxdb2.DefaultOptions().
			SetBatchSize(5000).
			SetFlushInterval(1000).
			SetHTTPRequestTimeout(600),
	)
	return Influx2Client{
		endpoint:     endpoint,
		organization: organization,
		bucket:       bucket,
		client:       client,
		writeApi:     client.WriteAPI(organization, bucket),
		queryApi:     client.QueryAPI(organization),
	}
}

func (c Influx2Client) Close() {
	c.client.Close()
}
