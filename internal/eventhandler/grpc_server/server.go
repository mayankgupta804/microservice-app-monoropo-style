package server

import (
	"context"
	"log"
	"net"

	"github.com/squadcast_assignment/internal/config"
	"github.com/squadcast_assignment/internal/eventhandler/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) EmitEvent(ctc context.Context, request *proto.Request) (*proto.Response, error) {
	id, incidentStatus := request.GetId(), request.GetIncidentStatus()
	log.Println(id, incidentStatus)
	return &proto.Response{Notify: "yes"}, nil
}

func StartGRPCServer(serverConf config.GRPCServer) error {
	listener, err := net.Listen("tcp", ":"+serverConf.Port)
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	proto.RegisterEventHandlerServer(srv, &server{})
	reflection.Register(srv)
	if err = srv.Serve(listener); err != nil {
		return err
	}
	return nil
}
