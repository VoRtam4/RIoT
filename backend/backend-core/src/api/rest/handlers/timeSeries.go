package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

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
	exportJob, err := domainLogicLayer.StartTimeSeriesExport(principal.UserID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(exportJob)
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
	exportJob, err := domainLogicLayer.StartTimeSeriesExportAggregateKPI(principal.UserID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(exportJob)
}

func CancelTimeSeriesExport(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	exportJob, err := domainLogicLayer.CancelTimeSeriesExport(principal.UserID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(exportJob)
}

func GetTimeSeriesExport(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	exportJob, err := domainLogicLayer.GetTimeSeriesExport(principal.UserID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(exportJob)
}

func TimeSeriesExport(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceTimeSeries, auth.OperationRead)
	if principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if _, err := domainLogicLayer.GetTimeSeriesExport(principal.UserID, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
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
