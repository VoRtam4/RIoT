package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func (c Influx2Client) buildReprocessFlux(req sharedModel.TimeSeriesReprocessReadRequest) string {
	quoted := make([]string, len(req.SDInstanceUIDs))
	for i, v := range req.SDInstanceUIDs {
		quoted[i] = fmt.Sprintf(`"%s"`, v)
	}
	ids := "[" + strings.Join(quoted, ", ") + "]"
	flux := fmt.Sprintf(`
from(bucket: "%s")
|> range(start: 0, stop: %s)
|> filter(fn: (r) => r["_measurement"] == "%s")
`, c.bucket, req.To.UTC().Format(time.RFC3339Nano), fmt.Sprintf("%s_%s", string(sharedModel.TimeSeriesTypeRaw), req.SDTypeUID))
	if len(req.SDInstanceUIDs) > 0 {
		flux += fmt.Sprintf(`
|> filter(fn: (r) => contains(value: r["sdInstanceUID"], set: %v))
`, ids)
	}
	flux += `
|> sort(columns: ["_time"], desc: false)
`
	return flux
}

func (c Influx2Client) buildReadFlux(plan sharedModel.QueryPlan) string {
	flux := fmt.Sprintf(`
import "strings"
from(bucket: "%s")
|> range(start: %s, stop: %s)
`, c.bucket, plan.From.Format(time.RFC3339Nano), plan.To.Format(time.RFC3339Nano))
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
|> group()
|> sort(columns: ["_time"], desc: %t)
`, plan.SortDesc)
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
	if plan.IsKPI {
		flux += `
|> group(columns: ["sdInstanceUID", "kpiDefinitionID"])
`
	} else {
		flux += `
|> group(columns: ["sdInstanceUID"])
`
	}
	flux += `
|> last()
|> group()
`
	flux += fmt.Sprintf(`
|> sort(columns: ["_time"], desc: %t)
`, plan.SortDesc)
	return flux
}

func (c Influx2Client) buildDistinctTagValuesFlux(plan sharedModel.QueryPlan, tag string) string {
	flux := fmt.Sprintf(`
import "strings"
from(bucket: "%s")
|> range(start: %s, stop: %s)
`, c.bucket, plan.From.Format(time.RFC3339Nano), plan.To.Format(time.RFC3339Nano))
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
