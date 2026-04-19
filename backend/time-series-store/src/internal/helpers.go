package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func makeKey(sdType string, kpiID uint32) string {
	return fmt.Sprintf("%s|%d", sdType, kpiID)
}

func isActive(jobKey string, jobID string) bool {
	v, ok := activeReprocessJobs.Load(jobKey)
	if !ok {
		return false
	}
	job, ok := v.(string)
	return ok && job == jobID
}

func extractTags(values map[string]interface{}) map[string]string {
	tags := make(map[string]string)
	for k, v := range values {
		switch k {
		case "_time", "_value", "_field", "_measurement", "_start", "_stop", "_source_timestamp", "result", "table":
			continue
		}
		if str, ok := v.(string); ok {
			tags[k] = str
		}
	}
	return tags
}

func buildAggregateKey(tags map[string]string) string {
	return tags["sdInstanceUID"] + "|" + tags["kpiDefinitionID"]
}

func cloneTags(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}

func extractFulfilled(point sharedModel.TimeSeriesDataPoint) (bool, bool) {
	raw, ok := point.Data["fulfilled"]
	if !ok {
		return false, false
	}

	switch v := raw.(type) {
	case bool:
		return v, true
	case string:
		return strings.EqualFold(v, "true"), true
	default:
		return false, false
	}
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func buildTagFilterFlux(node *sharedModel.FilterNode) (string, bool) {
	if node == nil {
		return "", false
	}
	expr := buildFilterExpr(node)
	if expr == "" {
		return "", false
	}
	return fmt.Sprintf(`|> filter(fn: (r) => %s)`, expr), true
}

func buildFilterExpr(node *sharedModel.FilterNode) string {
	if node == nil {
		return ""
	}
	switch node.Type {
	case sharedModel.FilterNodeTypeRule:
		return buildRuleExpr(node.Rule)
	case sharedModel.FilterNodeTypeLogical:
		if node.Operator == sharedModel.LogicalNot {
			if len(node.Nodes) != 1 {
				return ""
			}
			inner := buildFilterExpr(&node.Nodes[0])
			if inner == "" {
				return ""
			}
			return fmt.Sprintf("not (%s)", inner)
		}
		parts := make([]string, 0, len(node.Nodes))
		for _, n := range node.Nodes {
			expr := buildFilterExpr(&n)
			if expr != "" {
				parts = append(parts, expr)
			}
		}
		if len(parts) == 0 {
			return ""
		}
		op := "and"
		if node.Operator == sharedModel.LogicalOr {
			op = "or"
		}
		return "(" + strings.Join(parts, " "+op+" ") + ")"
	}
	return ""
}

func buildRuleExpr(rule *sharedModel.FilterRule) string {
	if rule == nil {
		return ""
	}

	tag := fmt.Sprintf(`r["%s"]`, rule.Tag)

	switch rule.Operator {
	case sharedModel.OpEQ:
		return fmt.Sprintf(`%s == %q`, tag, rule.Value)

	case sharedModel.OpNEQ:
		return fmt.Sprintf(`%s != %q`, tag, rule.Value)

	case sharedModel.OpContains:
		return fmt.Sprintf(`strings.containsStr(v: %s, substr: %q)`, tag, rule.Value)

	case sharedModel.OpPrefix:
		return fmt.Sprintf(`strings.hasPrefix(v: %s, prefix: %q)`, tag, rule.Value)

	case sharedModel.OpSuffix:
		return fmt.Sprintf(`strings.hasSuffix(v: %s, suffix: %q)`, tag, rule.Value)

	case sharedModel.OpRegex:
		return fmt.Sprintf(`%s =~ /%s/`, tag, escapeFluxRegex(rule.Value))

	case sharedModel.OpIn:
		values := strings.Split(rule.Value, ",")
		parts := make([]string, 0, len(values))
		for _, v := range values {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			parts = append(parts, fmt.Sprintf(`%s == %q`, tag, v))
		}
		if len(parts) == 0 {
			return ""
		}
		return "(" + strings.Join(parts, " or ") + ")"
	}

	return ""
}

func escapeFluxRegex(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `/`, `\/`)
	return s
}

func (s *aggregateSeriesState) addDuration(from, to time.Time) {
	if !s.HasLastValue || !to.After(from) {
		return
	}

	d := to.Sub(from)
	if s.LastBool {
		s.TrueWindowTime += d
	} else {
		s.FalseWindowTime += d
	}
	s.LastTime = to
}

func (s *aggregateSeriesState) buildWindowPoint(windowStart time.Time) sharedModel.TimeSeriesDataPoint {
	result := s.LastBool

	if s.TrueWindowTime > s.FalseWindowTime {
		result = true
	} else if s.TrueWindowTime < s.FalseWindowTime {
		result = false
	}

	return sharedModel.TimeSeriesDataPoint{
		Time: windowStart.UTC(),
		Tags: cloneTags(s.Tags),
		Data: map[string]interface{}{
			"fulfilled": result,
		},
	}
}

func (s *aggregateSeriesState) resetWindow() {
	s.TrueWindowTime = 0
	s.FalseWindowTime = 0
	if s.HasLastValue {
		s.HasFullWindow = true
	}
}

func (c Influx2Client) iterateFluxPoints(fluxQuery string, plan sharedModel.QueryPlan, isSnapshot bool, onPoint func(sharedModel.TimeSeriesDataPoint) (bool, error)) (bool, error) {
	result, err := c.queryApi.Query(context.Background(), fluxQuery)
	if err != nil {
		return false, err
	}
	for result.Next() {
		values := result.Record().Values()
		t, ok := values["_time"].(time.Time)
		if !ok {
			continue
		}
		if isSnapshot {
			t = plan.From
		}
		field, _ := values["_field"].(string)
		if field == "" {
			continue
		}
		point := sharedModel.TimeSeriesDataPoint{
			Time: t.UTC(),
			Tags: extractTags(values),
			Data: map[string]interface{}{
				field: values["_value"],
			},
		}
		stop, err := onPoint(point)
		if err != nil {
			return false, err
		}
		if stop {
			return true, nil
		}
	}
	if result.Err() != nil {
		return false, result.Err()
	}
	return false, nil
}

func shouldSkipAggregatedPoint(plan sharedModel.QueryPlan, state *aggregateStreamState, point sharedModel.TimeSeriesDataPoint) bool {
	if !plan.UseCursor || plan.Cursor == nil {
		return false
	}
	if state.cursorPassed {
		return false
	}
	cursorKey := plan.Cursor.SDInstanceUID
	if plan.Cursor.KPIDefinitionID != nil {
		cursorKey += "|" + fmt.Sprintf("%d", *plan.Cursor.KPIDefinitionID)
	}
	pointKey := buildAggregateKey(point.Tags)
	if point.Time.Before(plan.Cursor.Time) {
		return true
	}
	if point.Time.After(plan.Cursor.Time) {
		state.cursorPassed = true
		return false
	}
	if pointKey < cursorKey {
		return true
	}
	if pointKey == cursorKey {
		state.cursorPassed = true
		return true
	}
	state.cursorPassed = true
	return false
}

func (c Influx2Client) flushAggregateWindow(plan sharedModel.QueryPlan, state *aggregateStreamState, onBatch func([]sharedModel.TimeSeriesDataPoint, bool, bool) error, hasMoreData bool) (bool, error) {
	for _, key := range state.seriesOrder {
		series := state.series[key]
		if !series.HasLastValue {
			continue
		}
		if !series.HasLastValue {
			continue
		}
		point := series.buildWindowPoint(state.windowStart)
		if shouldSkipAggregatedPoint(plan, state, point) {
			series.resetWindow()
			continue
		}
		state.pointsBatch = append(state.pointsBatch, point)
		state.totalSent++
		series.resetWindow()
		if plan.Limit > 0 && state.totalSent >= plan.Limit {
			if len(state.pointsBatch) > 0 {
				if err := onBatch(state.pointsBatch, false, hasMoreData); err != nil {
					return true, err
				}
				state.pointsBatch = state.pointsBatch[:0]
			}
			return true, nil
		}
		if plan.Batch > 0 && len(state.pointsBatch) >= plan.Batch {
			if err := onBatch(state.pointsBatch, true, hasMoreData); err != nil {
				return false, err
			}
			state.pointsBatch = state.pointsBatch[:0]
		}
	}
	return false, nil
}

func (c Influx2Client) advanceAggregateWindows(plan sharedModel.QueryPlan, state *aggregateStreamState, target time.Time, onBatch func([]sharedModel.TimeSeriesDataPoint, bool, bool) error) (bool, error) {
	for state.windowStart.Before(plan.To) {
		windowEnd := minTime(state.windowStart.Add(state.windowSize), plan.To)
		if target.Before(windowEnd) {
			return false, nil
		}
		for _, key := range state.seriesOrder {
			s := state.series[key]
			if !s.HasLastValue {
				continue
			}
			if s.LastTime.Before(windowEnd) {
				s.addDuration(s.LastTime, windowEnd)
			}
		}
		hitLimit, err := c.flushAggregateWindow(plan, state, onBatch, true)
		if err != nil || hitLimit {
			return hitLimit, err
		}
		state.windowStart = windowEnd
	}
	return false, nil
}
