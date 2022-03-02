package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/belljustin/captainhook"
	pb "github.com/belljustin/captainhook/proto/captainhook"
	"github.com/belljustin/captainhook/storage/postgres"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")

	defaultTenantID, _ = uuid.FromBytes([]byte("default"))
)

type server struct {
	pb.UnimplementedCaptainhookServer

	storage captainhook.Storage
}

func (s *server) CreateApplication(ctx context.Context, createApp *pb.CreateApplicationRequest) (*pb.Application, error) {
	tenantID, err := parseTenantIDString(createApp.GetTenantId())
	if err != nil {
		return nil, err
	}

	application, err := captainhook.NewApplication(s.storage, tenantID, createApp.GetName())
	if err != nil {
		return nil, err
	}
	return application.ToProtobuf(), nil
}

func (s *server) GetApplication(ctx context.Context, getApp *pb.GetApplicationRequest) (*pb.Application, error) {
	sID := getApp.GetId()
	id, _ := uuid.Parse(sID)

	tenantID, err := parseTenantIDString(getApp.GetTenantId())
	if err != nil {
		return nil, err
	}

	app, err := captainhook.GetApplication(s.storage, tenantID, id)
	if err != nil {
		return nil, err
	}
	return app.ToProtobuf(), nil
}

func parseTenantIDString(sTenantID string) (uuid.UUID, error) {
	if sTenantID == "" {
		return defaultTenantID, nil
	} else {
		return uuid.Parse(sTenantID)
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	storage := postgres.NewStorage()
	pb.RegisterCaptainhookServer(s, &server{storage: storage})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
