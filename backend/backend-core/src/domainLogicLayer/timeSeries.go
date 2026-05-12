package domainLogicLayer

import (
	"fmt"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/rabbitmq/amqp091-go"
)

func StreamTimeSeriesAggregateKPI(userID uint32, input graphQLModel.TimeSeriesReadAggregateKPIInput, onBatch func(graphQLModel.TimeSeriesReadResponse) error) error {
	req, err := buildTimeSeriesAggregateReadRequest(userID, input)
	if err != nil {
		return err
	}
	params := loadParametersFromDB(graphQLModel.TimeSeriesTypeKpi, input.SdTypeID)
	return StreamTimeSeriesBase(userID, req, params, onBatch)
}

func StreamTimeSeries(userID uint32, input graphQLModel.TimeSeriesReadInput, onBatch func(graphQLModel.TimeSeriesReadResponse) error) error {
	req, err := buildTimeSeriesReadRequest(userID, input)
	if err != nil {
		return err
	}
	params := loadParametersFromDB(input.Type, input.SdTypeID)
	return StreamTimeSeriesBase(userID, req, params, onBatch)
}

func StreamTimeSeriesBase(userID uint32, req sharedModel.TimeSeriesReadRequest, params []graphQLModel.TimeSeriesParameter, onBatch func(graphQLModel.TimeSeriesReadResponse) error) error {
	session, err := newGeneratedTimeSeriesRPCSession()
	if err != nil {
		return err
	}
	defer session.Close()
	if err := session.PublishReadRequest(req); err != nil {
		return err
	}
	return rabbitmq.ConsumeRPCStream[sharedModel.TimeSeriesReadResponse](session.messages, session.correlationID, 3*time.Minute, func(resp sharedModel.TimeSeriesReadResponse, msg amqp091.Delivery) (bool, error) {
		if resp.Error != "" {
			return true, onBatch(graphQLModel.TimeSeriesReadResponse{
				Error: toErrorPtr(resp.Error),
			})
		}
		gqlResp := mapToGraphQLResponse(resp, params)
		if err := onBatch(gqlResp); err != nil {
			return true, err
		}
		if !resp.HasMoreBatches {
			return true, nil
		}
		return false, nil
	})
}

func ReadTimeSeriesAggregateKPI(userID uint32, input graphQLModel.TimeSeriesReadAggregateKPIInput) sharedUtils.Result[graphQLModel.TimeSeriesReadResponse] {
	req, err := buildTimeSeriesAggregateReadRequest(userID, input)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesReadResponse](err)
	}
	batch := 200
	req.Batch = &batch
	tsResp, err := executeTimeSeriesQuery(req)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesReadResponse](err)
	}
	params := loadParametersFromDB(graphQLModel.TimeSeriesTypeKpi, input.SdTypeID)
	gqlResp := mapToGraphQLResponse(tsResp, params)
	return sharedUtils.NewSuccessResult(gqlResp)
}

func ReadTimeSeries(userID uint32, input graphQLModel.TimeSeriesReadInput) sharedUtils.Result[graphQLModel.TimeSeriesReadResponse] {
	req, err := buildTimeSeriesReadRequest(userID, input)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesReadResponse](err)
	}
	batch := 200
	req.Batch = &batch
	tsResp, err := executeTimeSeriesQuery(req)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesReadResponse](err)
	}
	params := loadParametersFromDB(input.Type, input.SdTypeID)
	gqlResp := mapToGraphQLResponse(tsResp, params)
	return sharedUtils.NewSuccessResult(gqlResp)
}

func DistinctTimeSeriesTagValues(userID uint32, input graphQLModel.TimeSeriesDistinctTagValuesInput) sharedUtils.Result[graphQLModel.TimeSeriesDistinctTagValuesResponse] {
	req, err := buildTimeSeriesDistinctTagValuesRequest(userID, input)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesDistinctTagValuesResponse](err)
	}
	resp, err := executeTimeSeriesDistinctTagValuesQuery(req)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.TimeSeriesDistinctTagValuesResponse](err)
	}
	return sharedUtils.NewSuccessResult(graphQLModel.TimeSeriesDistinctTagValuesResponse{
		Values: resp.Values,
		Error:  toErrorPtr(resp.Error),
	})
}

func executeTimeSeriesQuery(req sharedModel.TimeSeriesReadRequest) (sharedModel.TimeSeriesReadResponse, error) {
	session, err := newGeneratedTimeSeriesRPCSession()
	if err != nil {
		return sharedModel.TimeSeriesReadResponse{}, err
	}
	defer session.Close()
	if err := session.PublishReadRequest(req); err != nil {
		return sharedModel.TimeSeriesReadResponse{}, err
	}
	var finalResponse sharedModel.TimeSeriesReadResponse
	err = rabbitmq.ConsumeRPCStream[sharedModel.TimeSeriesReadResponse](session.messages, session.correlationID, 3*time.Minute, func(resp sharedModel.TimeSeriesReadResponse, msg amqp091.Delivery) (bool, error) {
		if resp.Error != "" {
			return true, fmt.Errorf("%s", resp.Error)
		}
		if len(resp.Data) > 0 {
			finalResponse.Data = append(finalResponse.Data, resp.Data...)
		}
		finalResponse.HasMoreData = resp.HasMoreData
		finalResponse.NextCursor = resp.NextCursor
		if !resp.HasMoreBatches {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return sharedModel.TimeSeriesReadResponse{}, err
	}
	return finalResponse, nil
}

func executeTimeSeriesDistinctTagValuesQuery(req sharedModel.TimeSeriesDistinctTagValuesRequest) (sharedModel.TimeSeriesDistinctTagValuesResponse, error) {
	session, err := newGeneratedTimeSeriesRPCSession()
	if err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesResponse{}, err
	}
	defer session.Close()
	if err := session.PublishDistinctTagValuesRequest(req); err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesResponse{}, err
	}
	var finalResponse sharedModel.TimeSeriesDistinctTagValuesResponse
	err = rabbitmq.ConsumeRPCStream[sharedModel.TimeSeriesDistinctTagValuesResponse](session.messages, session.correlationID, 3*time.Minute, func(resp sharedModel.TimeSeriesDistinctTagValuesResponse, msg amqp091.Delivery) (bool, error) {
		if resp.Error != "" {
			return true, fmt.Errorf("%s", resp.Error)
		}
		finalResponse = resp
		return true, nil
	})
	if err != nil {
		return sharedModel.TimeSeriesDistinctTagValuesResponse{}, err
	}
	return finalResponse, nil
}
