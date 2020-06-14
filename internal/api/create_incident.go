package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/squadcast_assignment/internal/model/domain"
	"github.com/squadcast_assignment/internal/serializer"
	"github.com/squadcast_assignment/internal/service"
)

type CreateIncidentHandler struct {
	IncidentService service.IncidentService
}

func (h CreateIncidentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var reqBody serializer.CreateIncidentRequest

	err := decoder.Decode(&reqBody)
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

	ID, creationErr := h.IncidentService.CreateIncident(reqBody)
	if creationErr != nil {
		e := fmt.Errorf("incident creation failed: %s", creationErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(domain.ErrToJSON(e, http.StatusInternalServerError))
		return
	}

	successResponse := serializer.CreateIncidentResponse{
		ID:     strconv.Itoa(int(ID)),
		Status: "CREATED",
	}

	responseJson, err := json.Marshal(successResponse)
	if err != nil {
		e := fmt.Errorf("could not parse resp json: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(domain.ErrToJSON(e, http.StatusInternalServerError))
	}

	w.Write(responseJson)
}

func (h CreateIncidentHandler) validateParams(reqBody serializer.CreateIncidentRequest) *domain.Error {
	return nil
}
