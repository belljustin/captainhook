package main

import (
	"flag"
	"log"

	"github.com/hibiken/asynq"

	"github.com/belljustin/captainhook"
	"github.com/belljustin/captainhook/storage/postgres"
)

var (
	redisAddr = flag.String("redisAddr", "localhost:6379", "The redis address")
)

func main() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: *redisAddr},
		asynq.Config{Concurrency: 10},
	)
	storage := postgres.NewStorage()
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: *redisAddr})

	signMessageTaskHandler := captainhook.SignMessageTaskHandler{Storage: storage, AsynqClient: asynqClient}
	createSubscriptionTaskHandler := captainhook.CreateSubscriptionTaskHandler{Storage: storage}
	fanoutTaskHandler := captainhook.FanoutTaskHandler{Storage: storage, AsynqClient: asynqClient}
	deliveryTaskHandler := captainhook.DeliveryTaskHandler{}

	mux := asynq.NewServeMux()
	mux.HandleFunc(captainhook.TypeSignMessage, signMessageTaskHandler.Handle)
	mux.HandleFunc(captainhook.TypeCreateSubscription, createSubscriptionTaskHandler.Handle)
	mux.HandleFunc(captainhook.TypeFanoutMessage, fanoutTaskHandler.Handle)
	mux.HandleFunc(captainhook.TypeDeliveryMessage, deliveryTaskHandler.Handle)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
