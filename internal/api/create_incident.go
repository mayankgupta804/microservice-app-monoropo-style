package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/squadcast_assignment/internal/eventhandler/proto"
	"github.com/squadcast_assignment/internal/model/domain"
	"github.com/squadcast_assignment/internal/serializer"
	"github.com/squadcast_assignment/internal/service"
)

type CreateIncidentHandler struct {
	IncidentService service.IncidentService
	GRPCClient      proto.EventHandlerClient
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
		log.Printf("validation error: %v", validateErr.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(domain.ErrToJSON(validateErr, http.StatusBadRequest))
		return
	}

	ID, creationErr := h.IncidentService.CreateIncident(reqBody)
	if creationErr != nil {
		e := errors.New("incident creation failed")
		log.Printf("db error: %v", creationErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(domain.ErrToJSON(e, http.StatusInternalServerError))
		return
	}

	req := &proto.Request{Id: ID, IncidentStatus: "INCIDENT_CREATED"}
	ctx := context.Background()
	if resp, err := h.GRPCClient.EmitEvent(ctx, req); err == nil {
		log.Printf("Event handler service notified: %v", resp.Notify)
	} else {
		log.Printf("Event handler service did not send any response. Error: %v", err)
	}

	successResponse := serializer.CreateIncidentResponse{
		ID:     strconv.Itoa(int(ID)),
		Status: "INCIDENT CREATED",
	}

	responseJSON, err := json.Marshal(successResponse)
	if err != nil {
		e := fmt.Errorf("could not parse resp json: %s", err.Error())
		log.Printf("parsing error: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(domain.ErrToJSON(e, http.StatusInternalServerError))
	}

	w.Write(responseJSON)
}

func (h CreateIncidentHandler) validateParams(reqBody serializer.CreateIncidentRequest) *domain.Error {
	if reqBody.Message == "" {
		return domain.NewError("'message' field must not be empty")
	}
	return nil
}
