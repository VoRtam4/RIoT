package rest

import (
	"log"
	"net/http"
)

func SetupRESTServer() {
	log.Println("[REST] running on :8080")
	http.ListenAndServe(":8080", SetupRouter())
}
