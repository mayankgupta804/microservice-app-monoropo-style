package client

import (
	"github.com/squadcast_assignment/internal/config"
	"github.com/squadcast_assignment/internal/eventhandler/proto"
	"google.golang.org/grpc"
)

// GetGRPClient returns an instance of event handler client
func GetGRPClient(serverConf config.GRPCServer) (proto.EventHandlerClient, error) {
	conn, err := grpc.Dial(serverConf.Host+":"+serverConf.Port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := proto.NewEventHandlerClient(conn)
	return client, nil
}
