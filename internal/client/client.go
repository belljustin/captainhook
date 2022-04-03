package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/belljustin/captainhook/proto/captainhook"
)

func CreateApplication(client pb.CaptainhookClient, createApp *pb.CreateApplicationRequest) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	app, err := client.CreateApplication(ctx, createApp)
	if err != nil {
		log.Fatalf("%v.CreateApplication(_) = _, %v: ", client, err)
	}
	log.Println(app)

	return app.Id
}

func GetApplication(client pb.CaptainhookClient, getApp *pb.GetApplicationRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	app, err := client.GetApplication(ctx, getApp)
	if err != nil {
		log.Fatalf("%v.GetApplication(_) = _, %v: ", client, err)
	}
	log.Println(app)
}

func CreateMessage(client pb.CaptainhookClient, createMsg *pb.CreateMessageRequest) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	msg, err := client.CreateMessage(ctx, createMsg)
	if err != nil {
		log.Fatalf("%v.CreateMessage(_) = _, %v: ", client, err)
	}
	log.Println(msg)

	return msg.Id
}

func CreateSubscription(client pb.CaptainhookClient, createSub *pb.CreateSubscriptionRequest) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sub, err := client.CreateSubscription(ctx, createSub)
	if err != nil {
		log.Fatalf("%v.CreateSubscription(_) = _, %v: ", client, err)
	}
	log.Println("SubscriptionReceipt", sub)

	return sub.Id
}

func GetSubscriptions(client pb.CaptainhookClient, getSubs *pb.GetSubscriptionsRequest) (string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	subCollection, err := client.GetSubscriptions(ctx, getSubs)
	if err != nil {
		log.Fatalf("%v.GetSubscriptions(_) = _, %v: ", client, err)
	}
	log.Println("Subscriptions", subCollection)

	return subCollection.GetPrev(), subCollection.GetNext()
}

func New(serverAddr string) pb.CaptainhookClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return pb.NewCaptainhookClient(conn)
}
