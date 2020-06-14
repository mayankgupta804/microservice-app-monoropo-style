package webserver

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/squadcast_assignment/internal/api"
	"github.com/squadcast_assignment/internal/config"
	"github.com/squadcast_assignment/internal/service"
)

func NewCreateIncidentHandler(is service.IncidentService) api.CreateIncidentHandler {
	return api.CreateIncidentHandler{
		IncidentService: is,
	}
}

func NewGetIncidentHandler(is service.IncidentService) api.GetIncidentHandler {
	return api.GetIncidentHandler{
		IncidentService: is,
	}
}

func NewUpdateIncidentHandler(is service.IncidentService) api.UpdateIncidenthandler {
	return api.UpdateIncidenthandler{
		IncidentService: is,
	}
}

func NewDeleteIncidentHandler(is service.IncidentService) api.DeleteIncidentHandler {
	return api.DeleteIncidentHandler{
		IncidentService: is,
	}
}

func SetupRoutes(is service.IncidentService) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		str := `{"status": "OK"}`
		io.WriteString(w, str)
	}).Methods("GET")
	router.Handle("/incident", NewCreateIncidentHandler(is)).Methods(http.MethodPost)
	router.Handle("/incident/{incident_id}", NewUpdateIncidentHandler(is)).Methods(http.MethodPut)
	router.Handle("/incident/{incident_id}", NewGetIncidentHandler(is)).Methods(http.MethodGet)
	router.Handle("/incident/{incident_id}", NewDeleteIncidentHandler(is)).Methods(http.MethodDelete)

	return router

}

func StartServer(router *mux.Router, serverConf config.Server) error {
	log.Printf("Starting web server on port: %v", serverConf.Port)
	if err := http.ListenAndServe(":"+serverConf.Port, router); err != nil {
		return err
	}
	return nil
}
