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
	url, err := domainLogicLayer.StartTimeSeriesExportAggregateKPI(principal.UserID, input)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, map[string]string{
		"url": url,
	})
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
	url, err := domainLogicLayer.StartTimeSeriesExport(principal.UserID, input)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, map[string]string{
		"url": url,
	})
}
