package events

import (
	"context"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

type StreamSubscription[T any] struct {
	Channel <-chan T
	Close   func()
}

func SubscribeSDInstanceRegistered(ctx context.Context, filter *graphQLModel.SDInstanceRegisteredFilter, userID uint32) (*StreamSubscription[graphQLModel.SDInstance], error) {
	filterFn, err := BuildSDInstanceRegisteredFilter(userID, filter)
	if err != nil {
		return nil, err
	}
	busFilter := func(e Event) bool {
		payload, ok := e.Payload.(graphQLModel.SDInstance)
		return ok && filterFn(payload)
	}
	output := make(chan graphQLModel.SDInstance, 16)
	sub := GetEventBus().Subscribe([]EventType{SDInstanceRegisteredEventType}, busFilter, 16)
	go func() {
		defer close(output)
		defer GetEventBus().Unsubscribe(sub.ID)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-sub.Channel:
				if !ok {
					return
				}
				payload, ok := event.Payload.(graphQLModel.SDInstance)
				if !ok {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case output <- payload:
				}
			}
		}
	}()
	return &StreamSubscription[graphQLModel.SDInstance]{
		Channel: output,
		Close: func() {
			GetEventBus().Unsubscribe(sub.ID)
		},
	}, nil
}

func SubscribeRawDataPointArrived(ctx context.Context, filter *graphQLModel.RawDataPointArrivedFilter, userID uint32) (*StreamSubscription[[]graphQLModel.RawDataPoint], error) {
	filterFn, err := BuildRawDataPointArrivedFilter(userID, filter)
	if err != nil {
		return nil, err
	}
	busFilter := func(e Event) bool {
		payload, ok := e.Payload.([]graphQLModel.RawDataPoint)
		return ok && filterFn(payload)
	}
	output := make(chan []graphQLModel.RawDataPoint, 16)
	sub := GetEventBus().Subscribe([]EventType{RawDataPointReceivedEventType}, busFilter, 16)
	go func() {
		defer close(output)
		defer GetEventBus().Unsubscribe(sub.ID)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-sub.Channel:
				if !ok {
					return
				}
				payload, ok := event.Payload.([]graphQLModel.RawDataPoint)
				if !ok {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case output <- payload:
				}
			}
		}
	}()
	return &StreamSubscription[[]graphQLModel.RawDataPoint]{
		Channel: output,
		Close: func() {
			GetEventBus().Unsubscribe(sub.ID)
		},
	}, nil
}

func SubscribeKPIFulfillmentChecked(ctx context.Context, filter *graphQLModel.KPIFulfillmentCheckedFilter, userID uint32) (*StreamSubscription[[]graphQLModel.KPIFulfillmentCheckResult], error) {
	filterFn, err := BuildKPIFulfillmentCheckedFilter(userID, filter)
	if err != nil {
		return nil, err
	}
	busFilter := func(e Event) bool {
		payload, ok := e.Payload.([]graphQLModel.KPIFulfillmentCheckResult)
		return ok && filterFn(payload)
	}
	output := make(chan []graphQLModel.KPIFulfillmentCheckResult, 16)
	sub := GetEventBus().Subscribe([]EventType{KPIFulfillmentCheckedEventType}, busFilter, 16)
	go func() {
		defer close(output)
		defer GetEventBus().Unsubscribe(sub.ID)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-sub.Channel:
				if !ok {
					return
				}
				payload, ok := event.Payload.([]graphQLModel.KPIFulfillmentCheckResult)
				if !ok {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case output <- payload:
				}
			}
		}
	}()
	return &StreamSubscription[[]graphQLModel.KPIFulfillmentCheckResult]{
		Channel: output,
		Close: func() {
			GetEventBus().Unsubscribe(sub.ID)
		},
	}, nil
}
