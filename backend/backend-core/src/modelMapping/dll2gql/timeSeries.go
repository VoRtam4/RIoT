package dll2gql

import (
	"encoding/json"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func ToGraphQLTimeSeriesResponse(resp sharedModel.TimeSeriesReadResponse, params []graphQLModel.TimeSeriesParameter) graphQLModel.TimeSeriesReadResponse {
	data := make([]graphQLModel.TimeSeriesDataPoint, 0, len(resp.Data))
	for _, dp := range resp.Data {
		data = append(data, graphQLModel.TimeSeriesDataPoint{
			Time: dp.Time.Format(time.RFC3339Nano),
			Tags: stringToJSONString(dp.Tags),
			Data: toJSONString(dp.Data),
		})
	}
	var cursor *graphQLModel.TimeSeriesCursor
	if resp.NextCursor != nil {
		cursor = &graphQLModel.TimeSeriesCursor{
			Time:            resp.NextCursor.Time.Format(time.RFC3339Nano),
			SdInstanceUID:   resp.NextCursor.SDInstanceUID,
			KpiDefinitionID: resp.NextCursor.KPIDefinitionID,
		}
	}
	return graphQLModel.TimeSeriesReadResponse{
		Parameters:     params,
		Base:           stringToJSONString(resp.Base),
		Data:           data,
		HasMoreBatches: resp.HasMoreBatches,
		HasMoreData:    resp.HasMoreData,
		NextCursor:     cursor,
		Error:          toErrorPtr(resp.Error),
	}
}

func stringToJSONString(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

func toJSONString(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

func toErrorPtr(err string) *string {
	if err == "" {
		return nil
	}
	return &err
}
