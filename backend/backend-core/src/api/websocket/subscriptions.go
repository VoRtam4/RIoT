package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/handlers"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/events"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func startForwarder[T any](
	c *connection.Client,
	id string,
	topic string,
	sub *events.StreamSubscription[T],
) {
	go func() {
		defer sub.Close()

		for payload := range sub.Channel {
			c.SafeSend(sharedModel.WebSocketMessage{
				Type:    sharedModel.MessageEvent,
				ID:      id,
				Topic:   topic,
				Payload: payload,
			})
		}
	}()
}

func subscribeClient(c *connection.Client, msg sharedModel.WebSocketMessage) error {

	switch msg.Topic {

	case string(events.SDInstanceRegisteredEventType):

		principal := handlers.AuthorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationSubscribe)
		if principal == nil {
			return fmt.Errorf("unauthorized")
		}

		var filter graphQLModel.SDInstanceRegisteredFilter
		if msg.Payload != nil {
			raw, _ := json.Marshal(msg.Payload)
			_ = json.Unmarshal(raw, &filter)
		}

		sub, err := events.SubscribeSDInstanceRegistered(c.Ctx, &filter, principal.UserID)
		if err != nil {
			return err
		}

		c.Subscriptions[msg.ID] = sub.Close
		startForwarder(c, msg.ID, msg.Topic, sub)
		return nil

	case string(events.RawDataPointReceivedEventType):

		principal := handlers.AuthorizeOperation(c, msg, auth.ResourceRawData, auth.OperationSubscribe)
		if principal == nil {
			return fmt.Errorf("unauthorized")
		}

		var filter graphQLModel.RawDataPointArrivedFilter
		if msg.Payload != nil {
			raw, _ := json.Marshal(msg.Payload)
			_ = json.Unmarshal(raw, &filter)
		}

		sub, err := events.SubscribeRawDataPointArrived(c.Ctx, &filter, principal.UserID)
		if err != nil {
			return err
		}

		c.Subscriptions[msg.ID] = sub.Close
		startForwarder(c, msg.ID, msg.Topic, sub)
		return nil

	case string(events.KPIFulfillmentCheckedEventType):

		principal := handlers.AuthorizeOperation(c, msg, auth.ResourceKPIResults, auth.OperationSubscribe)
		if principal == nil {
			return fmt.Errorf("unauthorized")
		}

		var filter graphQLModel.KPIFulfillmentCheckedFilter
		if msg.Payload != nil {
			raw, _ := json.Marshal(msg.Payload)
			_ = json.Unmarshal(raw, &filter)
		}

		sub, err := events.SubscribeKPIFulfillmentChecked(c.Ctx, &filter, principal.UserID)
		if err != nil {
			return err
		}

		c.Subscriptions[msg.ID] = sub.Close
		startForwarder(c, msg.ID, msg.Topic, sub)
		return nil
	}

	return fmt.Errorf("unsupported topic")
}
