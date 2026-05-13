/**
 * @file fluxBuilder.go
 * @brief Sestavování Flux dotazů pro čtení historie, agregace a přepočet KPI.
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
	"fmt"
	"strings"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func (c Influx2Client) buildReprocessFlux(req sharedModel.TimeSeriesReprocessReadRequest) string {
	return c.buildReprocessFluxWindow(req, time.Unix(0, 0).UTC(), req.To.UTC())
}

func (c Influx2Client) buildReprocessFluxWindow(req sharedModel.TimeSeriesReprocessReadRequest, from time.Time, to time.Time) string {
	flux := fmt.Sprintf(`
from(bucket: "%s")
|> range(start: %s, stop: %s)
|> filter(fn: (r) => r["_measurement"] == "%s")
`, c.bucket, from.UTC().Format(time.RFC3339Nano), to.UTC().Format(time.RFC3339Nano), fmt.Sprintf("%s_%s", string(sharedModel.TimeSeriesTypeRaw), req.SDTypeUID))
	if len(req.SDInstanceUIDs) > 0 {
		parts := make([]string, 0, len(req.SDInstanceUIDs))
		for _, id := range req.SDInstanceUIDs {
			parts = append(parts, fmt.Sprintf(`r["sdInstanceUID"] == "%s"`, id))
		}
		flux += fmt.Sprintf(`
|> filter(fn: (r) => %s)
`, strings.Join(parts, " or "))
	}
	return flux
}

func (c Influx2Client) buildReadFlux(plan sharedModel.QueryPlan) string {
	return c.buildReadFluxForRange(plan, plan.From, plan.To)
}

func (c Influx2Client) buildReadFluxForRange(plan sharedModel.QueryPlan, from time.Time, to time.Time) string {
	flux := fmt.Sprintf(`
import "strings"
from(bucket: "%s")
|> range(start: %s, stop: %s)
`, c.bucket, from.UTC().Format(time.RFC3339Nano), to.UTC().Format(time.RFC3339Nano))
	flux += plan.MeasurementFilterFlux
	if plan.HasInstanceFilter {
		flux += plan.InstanceFilterFlux
	}
	if plan.IsKPI && plan.HasKPI && plan.KPIFilterFlux != "" {
		flux += plan.KPIFilterFlux
	}
	if plan.HasTagFilter {
		flux += plan.TagFilterFlux
	}
	if plan.UseSort {
		flux += fmt.Sprintf(`
|> group()
|> sort(columns: ["_time"], desc: %t)
`, plan.SortDesc)
	}
	return flux
}

func (c Influx2Client) buildSnapshotFlux(plan sharedModel.QueryPlan) string {
	flux := fmt.Sprintf(`
import "strings"
from(bucket: "%s")
|> range(start: 0, stop: %s)
`, c.bucket, plan.From.Format(time.RFC3339Nano))
	flux += plan.MeasurementFilterFlux
	if plan.HasInstanceFilter {
		flux += plan.InstanceFilterFlux
	}
	if plan.IsKPI && plan.HasKPI && plan.KPIFilterFlux != "" {
		flux += plan.KPIFilterFlux
	}
	if plan.HasTagFilter {
		flux += plan.TagFilterFlux
	}
	flux += `
|> last()
`
	if plan.NeedGrouping && plan.UseSort {
		flux += `
|>group()
`
	}
	if plan.UseSort {
		flux += fmt.Sprintf(`
|> sort(columns: ["_time"], desc: %t)
`, plan.SortDesc)
	}
	return flux
}

func (c Influx2Client) buildDistinctTagValuesFlux(plan sharedModel.QueryPlan, tag string) string {
	return c.buildDistinctTagValuesFluxForRange(plan, tag, plan.From, plan.To)
}

func (c Influx2Client) buildDistinctTagValuesFluxForRange(plan sharedModel.QueryPlan, tag string, from time.Time, to time.Time) string {
	flux := fmt.Sprintf(`
import "strings"
from(bucket: "%s")
|> range(start: %s, stop: %s)
`, c.bucket, from.UTC().Format(time.RFC3339Nano), to.UTC().Format(time.RFC3339Nano))
	flux += plan.MeasurementFilterFlux
	if plan.HasInstanceFilter {
		flux += plan.InstanceFilterFlux
	}
	if plan.IsKPI && plan.HasKPI && plan.KPIFilterFlux != "" {
		flux += plan.KPIFilterFlux
	}
	if plan.HasTagFilter {
		flux += plan.TagFilterFlux
	}
	flux += fmt.Sprintf(`
|> keep(columns: ["%s"])
|> group()
|> distinct(column: "%s")
	`, tag, tag)
	return flux
}

func (c Influx2Client) buildFirstPointFlux(plan sharedModel.QueryPlan, stop time.Time) string {
	flux := fmt.Sprintf(`
from(bucket: "%s")
|> range(start: 0, stop: %s)
`, c.bucket, stop.UTC().Format(time.RFC3339Nano))
	flux += plan.MeasurementFilterFlux
	if plan.HasInstanceFilter {
		flux += plan.InstanceFilterFlux
	}
	if plan.IsKPI && plan.HasKPI && plan.KPIFilterFlux != "" {
		flux += plan.KPIFilterFlux
	}
	if plan.HasTagFilter {
		flux += plan.TagFilterFlux
	}
	flux += `
|> first()
`
	if plan.NeedGrouping {
		flux += `
|> keep(columns: ["_time"])
|> group()
|> sort(columns: ["_time"], desc: false)
|> limit(n: 1)
`
	}
	return flux
}

func (c Influx2Client) buildReprocessFirstFlux(req sharedModel.TimeSeriesReprocessReadRequest, stop time.Time, needGrouping bool) string {
	flux := fmt.Sprintf(`
from(bucket: "%s")
|> range(start: 0, stop: %s)
|> filter(fn: (r) => r["_measurement"] == "%s")
`, c.bucket, stop.UTC().Format(time.RFC3339Nano), fmt.Sprintf("%s_%s", string(sharedModel.TimeSeriesTypeRaw), req.SDTypeUID))
	if len(req.SDInstanceUIDs) > 0 {
		parts := make([]string, 0, len(req.SDInstanceUIDs))
		for _, id := range req.SDInstanceUIDs {
			parts = append(parts, fmt.Sprintf(`r["sdInstanceUID"] == "%s"`, id))
		}
		flux += fmt.Sprintf(`
|> filter(fn: (r) => %s)
`, strings.Join(parts, " or "))
	}
	flux += `
|> first()
`
	if needGrouping {
		flux += `
|> keep(columns: ["_time"])
|> group()
|> sort(columns: ["_time"], desc: false)
|> limit(n: 1)
`
	}
	return flux
}
