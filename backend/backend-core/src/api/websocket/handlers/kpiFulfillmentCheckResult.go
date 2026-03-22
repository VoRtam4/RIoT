package handlers

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetKPIResults(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceKPIResults, auth.OperationRead); principal == nil {
		return
	}

	result := domainLogicLayer.GetKPIFulfillmentCheckResults()
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}

	sendSuccess(c, msg.ID, result.GetPayload())
}
