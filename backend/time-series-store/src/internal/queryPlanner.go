package internal

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func BuildQueryPlan(req sharedModel.TimeSeriesReadRequest) sharedModel.QueryPlan {
	now := time.Now().UTC()
	var plan sharedModel.QueryPlan
	plan.Type = req.Type
	plan.IsKPI = plan.Type == sharedModel.TimeSeriesTypeKPIResult
	if req.SortDesc != nil {
		plan.SortDesc = *req.SortDesc
	}
	if req.To == nil {
		plan.To = now
	} else {
		plan.To = req.To.UTC()
	}
	if req.From != nil {
		plan.From = req.From.UTC()
	} else {
		plan.From = time.Time{}
	}
	if plan.From.After(plan.To) {
		temp := plan.From
		plan.From = plan.To
		plan.To = temp
	}
	plan.NeedInitial = false
	if plan.From.Equal(plan.To) {
		plan.NeedInitial = true
	}
	if req.Cursor != nil {
		plan.Cursor = req.Cursor
		plan.UseCursor = true
		if !plan.UseAggregation {
			if plan.SortDesc {
				plan.To = req.Cursor.Time.UTC().Add(time.Nanosecond)
			} else {
				plan.From = req.Cursor.Time.UTC()
			}
		}
	}
	if !plan.From.IsZero() && plan.From.After(plan.To) {
		plan.From = plan.To
	}
	if !plan.SortDesc && !plan.From.IsZero() {
		plan.NeedInitial = true
	}
	plan.SDTypeUID = req.SDTypeUID
	plan.SDInstanceUIDs = req.SDInstanceUIDs
	plan.KPIDefinitionIDs = req.KPIDefinitionIDs
	if plan.IsKPI {
		plan.Measurement = fmt.Sprintf("%s_%s",
			string(sharedModel.TimeSeriesTypeKPIResult),
			req.SDTypeUID,
		)
	} else {
		plan.Measurement = fmt.Sprintf("%s_%s",
			string(sharedModel.TimeSeriesTypeRaw),
			req.SDTypeUID,
		)
	}
	if req.AggregateSeconds != nil && *req.AggregateSeconds > 0 {
		plan.AggregateSeconds = req.AggregateSeconds
		plan.UseAggregation = true
	}
	if req.Limit != nil && *req.Limit > 0 {
		plan.Limit = *req.Limit
	}

	if req.Batch != nil && *req.Batch > 0 {
		plan.Batch = *req.Batch
	} else {
		plan.Batch = 1000
	}
	plan.MeasurementFilterFlux = fmt.Sprintf(
		`|> filter(fn: (r) => r["_measurement"] == "%s")`,
		plan.Measurement,
	)
	if len(plan.SDInstanceUIDs) > 0 {
		plan.HasInstanceFilter = true

		parts := make([]string, 0, len(plan.SDInstanceUIDs))
		for _, id := range plan.SDInstanceUIDs {
			parts = append(parts,
				fmt.Sprintf(`r["sdInstanceUID"] == "%s"`, id))
		}

		plan.InstanceFilterFlux = fmt.Sprintf(
			`|> filter(fn: (r) => %s)`,
			strings.Join(parts, " or "),
		)
	}
	if plan.IsKPI && len(plan.KPIDefinitionIDs) > 0 {
		plan.HasKPI = true
		parts := make([]string, 0, len(plan.KPIDefinitionIDs))
		for _, id := range plan.KPIDefinitionIDs {
			parts = append(parts,
				fmt.Sprintf(`r["kpiDefinitionID"] == "%d"`, id))
		}
		plan.KPIFilterFlux = fmt.Sprintf(`|> filter(fn: (r) => %s)`, strings.Join(parts, " or "))
	}
	if req.Filters != nil {
		plan.TagFilterFlux, plan.HasTagFilter =
			buildTagFilterFlux(req.Filters)
	}
	log.Printf("[TS][PLAN FULL] %+v", plan)
	return plan
}
