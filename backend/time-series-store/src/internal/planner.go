/**
 * @file planner.go
 * @brief Tvorba interního plánu pro vykonání požadavků na čtení časových dat.
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

func BuildQueryPlan(req sharedModel.TimeSeriesReadRequest) sharedModel.QueryPlan {
	now := time.Now().UTC()
	var plan sharedModel.QueryPlan
	plan.Type = req.Type
	plan.IsKPI = plan.Type == sharedModel.TimeSeriesTypeKPIResult
	if req.SortDesc != nil {
		plan.UseSort = true
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
		if plan.UseSort {
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
	}
	if !plan.From.IsZero() && plan.From.After(plan.To) {
		plan.From = plan.To
	}
	if plan.UseSort && !plan.SortDesc && !plan.From.IsZero() {
		plan.NeedInitial = true
	}
	plan.SDTypeUID = req.SDTypeUID
	plan.SDInstanceUIDs = req.SDInstanceUIDs
	plan.KPIDefinitionUIDs = req.KPIDefinitionUIDs
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
	if plan.IsKPI {
		plan.NeedGrouping = len(plan.KPIDefinitionUIDs) != 1 || len(plan.SDInstanceUIDs) != 1
	} else {
		plan.NeedGrouping = len(plan.SDInstanceUIDs) != 1
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
	if plan.IsKPI && len(plan.KPIDefinitionUIDs) > 0 {
		plan.HasKPI = true
		parts := make([]string, 0, len(plan.KPIDefinitionUIDs))
		for _, uid := range plan.KPIDefinitionUIDs {
			parts = append(parts,
				fmt.Sprintf(`r["kpiDefinitionUID"] == "%s"`, uid))
		}
		plan.KPIFilterFlux = fmt.Sprintf(`|> filter(fn: (r) => %s)`, strings.Join(parts, " or "))
	}
	if req.Filters != nil {
		plan.TagFilterFlux, plan.HasTagFilter =
			buildTagFilterFlux(req.Filters)
	}
	return plan
}
