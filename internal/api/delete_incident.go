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

type DeleteIncidentHandler struct {
	IncidentService service.IncidentService
}

func (h DeleteIncidentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	incidentID := vars["incident_id"]
	id, err := strconv.Atoi(incidentID)

	if err != nil {
		e := fmt.Errorf("error convering string to int: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(domain.ErrToJSON(e, http.StatusBadRequest))
		return
	}

	if id <= 0 {
		e := fmt.Errorf("invalid incident id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(domain.ErrToJSON(e, http.StatusBadRequest))
		return
	}

	deletionErr := h.IncidentService.DeleteIncident(int64(id))
	if deletionErr != nil {
		e := fmt.Errorf("incident deletion failed: %s", deletionErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(domain.ErrToJSON(e, http.StatusInternalServerError))
		return
	}

	successResponse := serializer.DeleteIncidentResponse{
		Status: "INCIDENT DELETED",
	}

	responseJSON, err := json.Marshal(successResponse)
	if err != nil {
		e := fmt.Errorf("could not parse resp json: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(domain.ErrToJSON(e, http.StatusInternalServerError))
	}

	w.Write(responseJSON)
}

func (h DeleteIncidentHandler) validateParams(reqBody serializer.DeleteIncidentRequest) *domain.Error {
	return nil
}
