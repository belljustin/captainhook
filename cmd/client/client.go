package main

import (
	"context"
	"flag"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/belljustin/captainhook/proto/captainhook"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func createApplication(client pb.CaptainhookClient, createApp *pb.CreateApplicationRequest) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	app, err := client.CreateApplication(ctx, createApp)
	if err != nil {
		log.Fatalf("%v.CreateApplication(_) = _, %v: ", client, err)
	}
	log.Println(app)

	return app.Id
}

func getApplication(client pb.CaptainhookClient, getApp *pb.GetApplicationRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	app, err := client.GetApplication(ctx, getApp)
	if err != nil {
		log.Fatalf("%v.GetApplication(_) = _, %v: ", client, err)
	}
	log.Println(app)
}

func createMessage(client pb.CaptainhookClient, createMsg *pb.CreateMessageRequest) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	msg, err := client.CreateMessage(ctx, createMsg)
	if err != nil {
		log.Fatalf("%v.CreateMessage(_) = _, %v: ", client, err)
	}
	log.Println(msg)

	return msg.Id
}

func createSubscription(client pb.CaptainhookClient, createSub *pb.CreateSubscriptionRequest) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sub, err := client.CreateSubscription(ctx, createSub)
	if err != nil {
		log.Fatalf("%v.CreateSubscription(_) = _, %v: ", client, err)
	}
	log.Println(sub)

	return sub.Id
}

func main() {
	flag.Parse()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewCaptainhookClient(conn)

	appID := createApplication(client, &pb.CreateApplicationRequest{Name: "app_name"})
	getApplication(client, &pb.GetApplicationRequest{Id: appID})

	createMessage(client, &pb.CreateMessageRequest{
		ApplicationId: appID,
		Type:          "ch/test",
		Data:          []byte("hello world"),
	})

	createSubscription(client, &pb.CreateSubscriptionRequest{
		ApplicationId: appID,
		Name:          "testSubscription",
		Types:         []string{"app/payments", "app/withdrawals"},
	})
}
