package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/squadcast_assignment/internal/model/domain"
	"github.com/squadcast_assignment/internal/serializer"
	"github.com/squadcast_assignment/internal/service"
)

type UpdateIncidenthandler struct {
	IncidentService service.IncidentService
}

func (h UpdateIncidenthandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)
	incidentID := vars["incident_id"]
	id, err := strconv.Atoi(incidentID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(domain.ErrToJSON(err, http.StatusBadRequest))
		return
	}

	var reqBody serializer.UpdateIncidentRequest

	err = decoder.Decode(&reqBody)
	if err != nil {
		e := fmt.Errorf("json decode failed: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(domain.ErrToJSON(e, http.StatusBadRequest))
		return
	}

	validateErr := h.validateParams(reqBody)
	if validateErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(domain.ErrToJSON(validateErr, http.StatusBadRequest))
		return
	}

	updateErr := h.IncidentService.UpdateIncident(int64(id), reqBody)
	if updateErr != nil {
		e := fmt.Errorf("update failed: %s", updateErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(domain.ErrToJSON(e, http.StatusInternalServerError))
		return
	}

	successResponse := serializer.UpdateIncidentResponse{
		Status: "INCIDENT UPDATED",
	}

	responseJSON, err := json.Marshal(successResponse)
	if err != nil {
		e := fmt.Errorf("could not parse resp json: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(domain.ErrToJSON(e, http.StatusInternalServerError))
	}

	w.Write(responseJSON)
}

func (h UpdateIncidenthandler) validateParams(reqBody serializer.UpdateIncidentRequest) *domain.Error {
	return nil
}
