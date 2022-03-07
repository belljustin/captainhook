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

	signMessageTaskHandler := captainhook.SignMessageTaskHandler{Storage: storage}
	createSubscriptionTaskHandler := captainhook.CreateSubscriptionTaskHandler{Storage: storage}

	mux := asynq.NewServeMux()
	mux.HandleFunc(captainhook.TypeSignMessage, signMessageTaskHandler.Handle)
	mux.HandleFunc(captainhook.TypeCreateSubscription, createSubscriptionTaskHandler.Handle)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
