package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/squadcast_assignment/internal/eventhandler/proto"
	"github.com/squadcast_assignment/internal/model/domain"
	"github.com/squadcast_assignment/internal/serializer"
	"github.com/squadcast_assignment/internal/service"
)

type GetIncidentHandler struct {
	IncidentService service.IncidentService
	GRPCClient      proto.EventHandlerClient
}

func (h GetIncidentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	incidentID := vars["incident_id"]
	id, err := strconv.Atoi(incidentID)

	if err != nil {
		e := fmt.Errorf("string to int coversion failed: %s", err.Error())
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

	incident, fetchErr := h.IncidentService.GetIncident(int64(id))
	if fetchErr != nil {
		e := fmt.Errorf("error fetching incident: %s", fetchErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(domain.ErrToJSON(e, http.StatusInternalServerError))
		return
	}

	req := &proto.Request{Id: int64(id), IncidentStatus: "INCIDENT_FETCHED"}
	ctx := context.Background()
	if resp, err := h.GRPCClient.EmitEvent(ctx, req); err == nil {
		log.Printf("Event handler service notified: %v", resp.Notify)
	} else {
		log.Printf("Event handler service did not send any response. Error: %v", err)
	}

	successResponse := serializer.ReadIncidentResponse{
		Status:   incident.Status,
		Ack:      incident.Ack,
		Message:  incident.Message,
		Comments: incident.Comment,
	}

	responseJSON, err := json.Marshal(successResponse)
	if err != nil {
		e := fmt.Errorf("could not parse resp json: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(domain.ErrToJSON(e, http.StatusInternalServerError))
	}

	w.Write(responseJSON)
}
