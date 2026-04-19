package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/events"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	eventType := chi.URLParam(r, "eventType")
	resource, err := resourceForSSEEventType(eventType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	principal := authorizeOperation(w, r, resource, auth.OperationSubscribe)
	if principal == nil {
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()

	ctx := r.Context()

	switch eventType {

	case string(events.SDInstanceRegisteredEventType):
		var filter graphQLModel.SDInstanceRegisteredFilter
		json.NewDecoder(r.Body).Decode(&filter)
		sub, err := events.SubscribeSDInstanceRegistered(ctx, &filter, principal.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		streamSSE(w, flusher, sub, eventType)

	case string(events.RawDataPointReceivedEventType):
		var filter graphQLModel.RawDataPointArrivedFilter
		json.NewDecoder(r.Body).Decode(&filter)
		sub, err := events.SubscribeRawDataPointArrived(ctx, &filter, principal.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		streamSSE(w, flusher, sub, eventType)

	case string(events.KPIFulfillmentCheckedEventType):
		var filter graphQLModel.KPIFulfillmentCheckedFilter
		json.NewDecoder(r.Body).Decode(&filter)
		sub, err := events.SubscribeKPIFulfillmentChecked(ctx, &filter, principal.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		streamSSE(w, flusher, sub, eventType)

	default:
		http.Error(w, "unknown event type", http.StatusBadRequest)
	}
}

func resourceForSSEEventType(eventType string) (string, error) {
	switch eventType {
	case string(events.SDInstanceRegisteredEventType):
		return auth.ResourceSDInstances, nil
	case string(events.RawDataPointReceivedEventType):
		return auth.ResourceRawData, nil
	case string(events.KPIFulfillmentCheckedEventType):
		return auth.ResourceKPIResults, nil
	default:
		return "", fmt.Errorf("unknown event type")
	}
}

func streamSSE[T any](w http.ResponseWriter, flusher http.Flusher, sub *events.StreamSubscription[T], topic string) {
	defer sub.Close()
	for payload := range sub.Channel {
		data, err := json.Marshal(payload)
		if err != nil {
			continue
		}
		fmt.Fprintf(w, "event: %s\n", topic)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}
}
