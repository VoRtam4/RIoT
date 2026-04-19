package gql2dll

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func ToDLLTimeSeriesReadKPIRequest(input graphQLModel.TimeSeriesReadAggregateKPIInput, sdTypeUID string, sdInstanceUIDs []string) sharedModel.TimeSeriesReadRequest {
	var from *time.Time
	if input.From != nil {
		if t, err := time.Parse(time.RFC3339Nano, *input.From); err == nil {
			from = &t
		} else if t, err := time.Parse(time.RFC3339, *input.From); err == nil {
			from = &t
		}
	}
	var to *time.Time
	if input.To != nil {
		if t, err := time.Parse(time.RFC3339Nano, *input.To); err == nil {
			to = &t
		} else if t, err := time.Parse(time.RFC3339, *input.To); err == nil {
			to = &t
		}
	}
	var cursor *sharedModel.TimeSeriesCursor
	if input.Cursor != nil {
		var t time.Time
		if parsed, err := time.Parse(time.RFC3339Nano, input.Cursor.Time); err == nil {
			t = parsed
		} else if parsed, err := time.Parse(time.RFC3339, input.Cursor.Time); err == nil {
			t = parsed
		}
		cursor = &sharedModel.TimeSeriesCursor{
			Time:            t,
			SDInstanceUID:   input.Cursor.SdInstanceUID,
			KPIDefinitionID: input.Cursor.KpiDefinitionID,
		}
	}
	sortDesc := false
	return sharedModel.TimeSeriesReadRequest{
		Type:             sharedModel.TimeSeriesTypeKPIResult,
		SDTypeUID:        sdTypeUID,
		SDInstanceUIDs:   sdInstanceUIDs,
		KPIDefinitionIDs: input.KpiDefinitionIDs,
		From:             from,
		To:               to,
		AggregateSeconds: &input.AggregateSeconds,
		SortDesc:         &sortDesc,
		Limit:            input.Limit,
		Batch:            input.Batch,
		Filters:          toDLLFilterNode(input.Filters),
		Cursor:           cursor,
	}
}

func ToDLLTimeSeriesReadRequest(input graphQLModel.TimeSeriesReadInput, sdTypeUID string, sdInstanceUIDs []string) sharedModel.TimeSeriesReadRequest {
	var from *time.Time
	if input.From != nil {
		if t, err := time.Parse(time.RFC3339Nano, *input.From); err == nil {
			from = &t
		} else if t, err := time.Parse(time.RFC3339, *input.From); err == nil {
			from = &t
		}
	}
	var to *time.Time
	if input.To != nil {
		if t, err := time.Parse(time.RFC3339Nano, *input.To); err == nil {
			to = &t
		} else if t, err := time.Parse(time.RFC3339, *input.To); err == nil {
			to = &t
		}
	}
	var cursor *sharedModel.TimeSeriesCursor
	if input.Cursor != nil {
		var t time.Time
		if parsed, err := time.Parse(time.RFC3339Nano, input.Cursor.Time); err == nil {
			t = parsed
		} else if parsed, err := time.Parse(time.RFC3339, input.Cursor.Time); err == nil {
			t = parsed
		}
		cursor = &sharedModel.TimeSeriesCursor{
			Time:            t,
			SDInstanceUID:   input.Cursor.SdInstanceUID,
			KPIDefinitionID: input.Cursor.KpiDefinitionID,
		}
	}
	return sharedModel.TimeSeriesReadRequest{
		Type:             sharedModel.TimeSeriesType(input.Type),
		SDTypeUID:        sdTypeUID,
		SDInstanceUIDs:   sdInstanceUIDs,
		KPIDefinitionIDs: input.KpiDefinitionIDs,
		From:             from,
		To:               to,
		Limit:            input.Limit,
		SortDesc:         input.SortDesc,
		Batch:            input.Batch,
		Filters:          toDLLFilterNode(input.Filters),
		Cursor:           cursor,
	}
}

func ToDLLTimeSeriesDistinctTagValuesRequest(input graphQLModel.TimeSeriesDistinctTagValuesInput, sdTypeUID string, sdInstanceUIDs []string) sharedModel.TimeSeriesDistinctTagValuesRequest {
	var from *time.Time
	if input.From != nil {
		if t, err := time.Parse(time.RFC3339Nano, *input.From); err == nil {
			from = &t
		} else if t, err := time.Parse(time.RFC3339, *input.From); err == nil {
			from = &t
		}
	}
	var to *time.Time
	if input.To != nil {
		if t, err := time.Parse(time.RFC3339Nano, *input.To); err == nil {
			to = &t
		} else if t, err := time.Parse(time.RFC3339, *input.To); err == nil {
			to = &t
		}
	}
	return sharedModel.TimeSeriesDistinctTagValuesRequest{
		Type:             sharedModel.TimeSeriesType(input.Type),
		SDTypeUID:        sdTypeUID,
		SDInstanceUIDs:   sdInstanceUIDs,
		KPIDefinitionIDs: input.KpiDefinitionIDs,
		From:             from,
		To:               to,
		Tag:              input.Tag,
		Filters:          toDLLFilterNode(input.Filters),
	}
}

func toDLLFilterNode(node *graphQLModel.FilterNodeInput) *sharedModel.FilterNode {
	if node == nil {
		return nil
	}
	result := &sharedModel.FilterNode{
		Type: sharedModel.FilterNodeType(node.Type),
	}
	if node.Operator != nil {
		result.Operator = sharedModel.LogicalOperator(*node.Operator)
	}
	if node.Rule != nil {
		result.Rule = &sharedModel.FilterRule{
			Tag:      node.Rule.Tag,
			Operator: sharedModel.FilterOperator(node.Rule.Operator),
			Value:    node.Rule.Value,
		}
	}
	if len(node.Nodes) > 0 {
		result.Nodes = make([]sharedModel.FilterNode, 0, len(node.Nodes))
		for i := range node.Nodes {
			child := toDLLFilterNode(&node.Nodes[i])
			if child != nil {
				result.Nodes = append(result.Nodes, *child)
			}
		}
	}
	return result
}
