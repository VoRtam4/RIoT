package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/events"
)

func EventsStream(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceEvents, auth.OperationSubscribe) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}
	query := r.URL.Query().Get("topics")
	var subscribedTypes []events.EventType
	if query == "" {
		subscribedTypes = []events.EventType{
			events.SDInstanceRegisteredEventType,
			events.KPIFulfillmentCheckedEventType,
		}
	} else {
		for _, t := range strings.Split(query, ",") {
			subscribedTypes = append(subscribedTypes, events.EventType(t))
		}
	}
	sub := events.GetEventBus().Subscribe(subscribedTypes, 32)
	defer events.GetEventBus().Unsubscribe(sub.ID)
	fmt.Fprintf(w, ": connected\n\n")
	flusher.Flush()
	for {
		select {
		case <-r.Context().Done():
			return

		case event, ok := <-sub.Channel:
			if !ok {
				return
			}

			data, err := json.Marshal(event)
			if err != nil {
				continue
			}

			fmt.Fprintf(w, "event: %s\n", event.Type)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
