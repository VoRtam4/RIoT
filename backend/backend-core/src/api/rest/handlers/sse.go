package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/events"
)

func EventStream(w http.ResponseWriter, r *http.Request) {

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")

	for {
		select {

		case ev := <-events.SDInstanceRegistered:
			send(w, flusher, ev)

		case ev := <-events.KPIFulfillmentChecked:
			send(w, flusher, ev)
		}
	}
}

func send(w http.ResponseWriter, f http.Flusher, v any) {
	data, _ := json.Marshal(v)
	fmt.Fprintf(w, "data: %s\n\n", data)
	f.Flush()
}
