package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/go-chi/chi/v5"
)

var address = sharedUtils.GetEnvironmentVariableValue("BACKEND_CORE_URL").GetPayloadOrDefault("http://localhost:9090")

func ReadTimeSeries(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	var input graphQLModel.TimeSeriesReadInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.ReadTimeSeries(principal.UserID, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func DistinctTimeSeriesTagValues(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	var input graphQLModel.TimeSeriesDistinctTagValuesInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.DistinctTimeSeriesTagValues(principal.UserID, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func ReadTimeSeriesAggregateKpi(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	var input graphQLModel.TimeSeriesReadAggregateKPIInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.ReadTimeSeriesAggregateKPI(principal.UserID, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func StartTimeSeriesExport(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	var input graphQLModel.TimeSeriesReadInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	url, err := domainLogicLayer.StartTimeSeriesExport(principal.UserID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"url": address + "/rest/time-series/export/" + url,
	})
}

func StartTimeSeriesExportAggregateKpi(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	var input graphQLModel.TimeSeriesReadAggregateKPIInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	url, err := domainLogicLayer.StartTimeSeriesExportAggregateKPI(principal.UserID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"url": address + "/rest/time-series/export/" + url,
	})
}

func TimeSeriesExport(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceTimeSeries, auth.OperationSubscribe); principal == nil {
		return
	}
	id := chi.URLParam(r, "id")
	job, ok := domainLogicLayer.GetExportJob(id)
	if !ok {
		http.Error(w, "not found", 404)
		return
	}
	<-job.Done
	if job.Err != nil {
		http.Error(w, job.Err.Error(), 500)
		return
	}
	file, err := os.Open(job.FilePath)
	if err != nil {
		http.Error(w, "file error", 500)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=export.csv")
	io.Copy(w, file)
}
