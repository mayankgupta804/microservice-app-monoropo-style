package webserver

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/squadcast_assignment/internal/api"
	"github.com/squadcast_assignment/internal/config"
	"github.com/squadcast_assignment/internal/eventhandler/proto"
	"github.com/squadcast_assignment/internal/service"
)

func NewCreateIncidentHandler(is service.IncidentService, grpcClient proto.EventHandlerClient) api.CreateIncidentHandler {
	return api.CreateIncidentHandler{
		IncidentService: is,
		GRPCClient:      grpcClient,
	}
}

func NewGetIncidentHandler(is service.IncidentService, grpcClient proto.EventHandlerClient) api.GetIncidentHandler {
	return api.GetIncidentHandler{
		IncidentService: is,
		GRPCClient:      grpcClient,
	}
}

func NewUpdateIncidentHandler(is service.IncidentService, grpcClient proto.EventHandlerClient) api.UpdateIncidenthandler {
	return api.UpdateIncidenthandler{
		IncidentService: is,
		GRPCClient:      grpcClient,
	}
}

func NewDeleteIncidentHandler(is service.IncidentService, grpcClient proto.EventHandlerClient) api.DeleteIncidentHandler {
	return api.DeleteIncidentHandler{
		IncidentService: is,
		GRPCClient:      grpcClient,
	}
}

func SetupRoutes(is service.IncidentService, grpcClient proto.EventHandlerClient) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		str := `{"status": "OK"}`
		io.WriteString(w, str)
	}).Methods(http.MethodGet)
	router.Handle("/incident", NewCreateIncidentHandler(is, grpcClient)).Methods(http.MethodPost)
	router.Handle("/incident/{incident_id}", NewUpdateIncidentHandler(is, grpcClient)).Methods(http.MethodPut)
	router.Handle("/incident/{incident_id}", NewGetIncidentHandler(is, grpcClient)).Methods(http.MethodGet)
	router.Handle("/incident/{incident_id}", NewDeleteIncidentHandler(is, grpcClient)).Methods(http.MethodDelete)

	return router
}

func StartServer(router *mux.Router, serverConf config.Server) error {
	log.Printf("Starting web server on port: %v", serverConf.Port)
	if err := http.ListenAndServe(":"+serverConf.Port, router); err != nil {
		return err
	}
	return nil
}
