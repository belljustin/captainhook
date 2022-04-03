package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc"

	"github.com/belljustin/captainhook/captainhook"
	pb "github.com/belljustin/captainhook/proto/captainhook"
	"github.com/belljustin/captainhook/storage/postgres"
)

var (
	defaultTenantID, _ = uuid.FromBytes([]byte("default"))
)

type Server struct {
	pb.UnimplementedCaptainhookServer

	storage     captainhook.Storage
	asynqClient *asynq.Client

	port int
}

func New(port int, redisAddr string) *Server {
	storage := postgres.NewStorage()
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})

	return &Server{
		storage:     storage,
		asynqClient: asynqClient,

		port: port,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCaptainhookServer(grpcServer, s)
	log.Printf("server listening at %v", lis.Addr())

	return grpcServer.Serve(lis)
}

func (s *Server) CreateApplication(ctx context.Context, createApp *pb.CreateApplicationRequest) (*pb.Application, error) {
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

func (s *Server) GetApplication(ctx context.Context, getApp *pb.GetApplicationRequest) (*pb.Application, error) {
	sID := getApp.GetId()
	id, err := uuid.Parse(sID)
	if err != nil {
		return nil, err
	}

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

func (s *Server) CreateMessage(ctx context.Context, createMsg *pb.CreateMessageRequest) (*pb.MessageReceipt, error) {
	tenantID, err := parseTenantIDString(createMsg.GetTenantId())
	if err != nil {
		return nil, err
	}

	appID, err := uuid.Parse(createMsg.GetApplicationId())
	if err != nil {
		return nil, err
	}

	id, err := captainhook.CreateMessage(s.asynqClient, tenantID, appID, createMsg.GetType(), createMsg.GetData())
	if err != nil {
		return nil, err
	}
	return &pb.MessageReceipt{
		TenantId:      tenantID.String(),
		Id:            id.String(),
		ApplicationId: appID.String(),
	}, nil
}

func (s *Server) CreateSubscription(ctx context.Context, createSub *pb.CreateSubscriptionRequest) (*pb.SubscriptionReceipt, error) {
	tenantID, err := parseTenantIDString(createSub.GetTenantId())
	if err != nil {
		return nil, err
	}

	appID, err := uuid.Parse(createSub.GetApplicationId())
	if err != nil {
		return nil, err
	}

	if createSub.GetEndpoint() == "" {
		return nil, errors.New("'Endpoint' is a required field")
	}
	endpoint, err := url.Parse(createSub.GetEndpoint())
	if err != nil {
		return nil, err
	}

	for _, subType := range createSub.GetTypes() {
		if strings.Contains(subType, ",") {
			return nil, errors.New("'types' must not contain a ',' character")
		}
	}

	id, err := captainhook.CreateSubscription(s.asynqClient, tenantID, appID, createSub.GetName(), createSub.GetTypes(), endpoint)
	if err != nil {
		return nil, err
	}
	return &pb.SubscriptionReceipt{
		TenantId:      tenantID.String(),
		ApplicationId: appID.String(),
		Id:            id.String(),
	}, nil
}

type PaginationOpt struct {
	token string
	size  int32
}

func (opt *PaginationOpt) GetPageToken() string {
	return opt.token
}

func (opt *PaginationOpt) GetPageSize() int32 {
	if opt.size == 0 {
		return 20
	}
	return opt.size
}

func (s *Server) GetSubscriptions(ctx context.Context, getSubs *pb.GetSubscriptionsRequest) (*pb.SubscriptionCollection, error) {
	tenantID, err := parseTenantIDString(getSubs.GetTenantId())
	if err != nil {
		return nil, err
	}

	appID, err := uuid.Parse(getSubs.GetApplicationId())
	if err != nil {
		return nil, err
	}

	pageOpt := &PaginationOpt{getSubs.Page, getSubs.Size}

	subCol, err := captainhook.GetSubscriptions(s.storage, tenantID, appID, pageOpt)
	if err != nil {
		return nil, err
	}
	return subCol.ToProtobuf(), nil
}

func parseTenantIDString(sTenantID string) (uuid.UUID, error) {
	if sTenantID == "" {
		return defaultTenantID, nil
	} else {
		return uuid.Parse(sTenantID)
	}
}
