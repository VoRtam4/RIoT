package handlers

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func StreamTimeSeries(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	input, err := parsePayload[graphQLModel.TimeSeriesReadInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	err = domainLogicLayer.StreamTimeSeries(principal.UserID, input,
		func(batch graphQLModel.TimeSeriesReadResponse) error {
			sendSuccess(c, msg.ID, batch)
			return nil
		},
	)
	if err != nil {
		sendError(c, msg.ID, err.Error())
	}
}

func StreamTimeSeriesAggregateKpi(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	input, err := parsePayload[graphQLModel.TimeSeriesReadAggregateKPIInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	err = domainLogicLayer.StreamTimeSeriesAggregateKPI(principal.UserID, input,
		func(batch graphQLModel.TimeSeriesReadResponse) error {
			sendSuccess(c, msg.ID, batch)
			return nil
		},
	)
	if err != nil {
		sendError(c, msg.ID, err.Error())
	}
}

func StartTimeSeriesExportAggregateKpi(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	input, err := parsePayload[graphQLModel.TimeSeriesReadAggregateKPIInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	exportJob, err := domainLogicLayer.StartTimeSeriesExportAggregateKPI(principal.UserID, input)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, exportJob)
}

func StartTimeSeriesExport(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	input, err := parsePayload[graphQLModel.TimeSeriesReadInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	exportJob, err := domainLogicLayer.StartTimeSeriesExport(principal.UserID, input)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, exportJob)
}

func CancelTimeSeriesExport(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	id, ok := parseID(payload)
	if !ok {
		sendError(c, msg.ID, "invalid id")
		return
	}
	exportJob, err := domainLogicLayer.CancelTimeSeriesExport(principal.UserID, id)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, exportJob)
}

func GetTimeSeriesExport(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	id, ok := parseID(payload)
	if !ok {
		sendError(c, msg.ID, "invalid id")
		return
	}
	exportJob, err := domainLogicLayer.GetTimeSeriesExport(principal.UserID, id)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, exportJob)
}

func GetTimeSeriesDistinctTagValues(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	input, err := parsePayload[graphQLModel.TimeSeriesDistinctTagValuesInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	result := domainLogicLayer.DistinctTimeSeriesTagValues(principal.UserID, input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}
