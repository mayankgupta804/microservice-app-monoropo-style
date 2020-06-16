package server

import (
	"context"
	"encoding/json"
	"log"
	"net"

	"github.com/squadcast_assignment/internal/config"
	"github.com/squadcast_assignment/internal/eventhandler/proto"
	queue "github.com/squadcast_assignment/internal/infrastructure/workqueue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	queue.QueueClient
}

func (s *server) EmitEvent(ctc context.Context, request *proto.Request) (*proto.Response, error) {
	id, incidentStatus := request.GetId(), request.GetIncidentStatus()

	data := struct {
		ID             int64  `json:"id"`
		IncidentStatus string `json:"incident_status"`
	}{
		ID:             id,
		IncidentStatus: incidentStatus,
	}
	msg, err := json.Marshal(&data)
	if err != nil {
		log.Printf("error marshalling json: %v", err)
		return &proto.Response{Notify: "no"}, err
	}
	if err := s.QueueClient.Publish(config.App.Queue.Exchange, "", msg); err != nil {
		return &proto.Response{Notify: "no"}, err

	}
	return &proto.Response{Notify: "yes"}, nil
}

// StartGRPCServer starts the GRPC server at the given port
func StartGRPCServer(q queue.QueueClient, serverConf config.GRPCServer) error {
	listener, err := net.Listen("tcp", ":"+serverConf.Port)
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	proto.RegisterEventHandlerServer(srv, &server{q})
	reflection.Register(srv)
	log.Printf("Starting GRPC server on port: %s", serverConf.Port)
	if err = srv.Serve(listener); err != nil {
		return err
	}
	return nil
}
