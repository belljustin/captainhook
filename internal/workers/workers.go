package workers

import (
	"github.com/belljustin/captainhook/captainhook"
	"github.com/belljustin/captainhook/storage/postgres"
	"github.com/hibiken/asynq"
)

type Workers struct {
	redisAddr string

	asynqClient asynq.Client
	storage     captainhook.Storage
}

func New(redisAddr string) *Workers {
	return &Workers{redisAddr: redisAddr}
}

func (w Workers) Run() error {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: w.redisAddr},
		asynq.Config{Concurrency: 10},
	)
	storage := postgres.NewStorage()
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: w.redisAddr})
	deliverer := captainhook.NewHttpDeliverer()
	defer asynqClient.Close()

	createSubscriptionTaskHandler := captainhook.CreateSubscriptionTaskHandler{Storage: storage}
	fanoutTaskHandler := captainhook.FanoutTaskHandler{Storage: storage, AsynqClient: asynqClient}
	deliveryTaskHandler := captainhook.DeliveryTaskHandler{Deliverer: deliverer}

	mux := asynq.NewServeMux()
	mux.HandleFunc(captainhook.TypeCreateSubscription, createSubscriptionTaskHandler.Handle)
	mux.HandleFunc(captainhook.TypeFanoutMessage, fanoutTaskHandler.Handle)
	mux.HandleFunc(captainhook.TypeDeliveryMessage, deliveryTaskHandler.Handle)

	return srv.Run(mux)
}
