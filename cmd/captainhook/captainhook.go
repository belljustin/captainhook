package main

import (
	"errors"
	"fmt"
	pb "github.com/belljustin/captainhook/proto/captainhook"
	"os"

	"github.com/spf13/cobra"

	"github.com/belljustin/captainhook/cmd/client"
)

var rootCmd = &cobra.Command{
	Use:   "captainhook",
	Short: "Captainhook is a webhook delivery service",
	Long: `A webhook delivery service built with love by Justin Bell. Complete documentation is available at
					https://github.com/belljustin/captainhook`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	// captainhook app
	{
		// captainhook app create
		createApplicationCmd.Flags().StringVarP(&appName, "name", "n", "", "Human readable name that will be associated with the application")
		applicationCmd.AddCommand(createApplicationCmd)
		// captainhook app get
		applicationCmd.AddCommand(getApplicationCmd)
	}
	rootCmd.AddCommand(applicationCmd)

	// captainhook sub
	{
		// captainhook sub create
		subCmd.AddCommand(createSubCmd)
		// captainhook sub get
		subCmd.AddCommand(getSubCmd)
	}
	rootCmd.AddCommand(subCmd)

	// captainhook msg
	{
		// captainhook msg create
		msgCmd.AddCommand(createMsgCmd)
	}
	rootCmd.AddCommand(msgCmd)

	// captainhook version
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Captainhook",
	Long:  `All software has versions. This is Captainhook's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Captainhook webhook delivery service v0.1 -- HEAD")
	},
}

var (
	serverAddr = "localhost:50051"

	appName string
)

var applicationCmd = &cobra.Command{
	Use: "app",
}

var createApplicationCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(serverAddr)
		client.CreateApplication(c, &pb.CreateApplicationRequest{Name: appName})
	},
}

var getApplicationCmd = &cobra.Command{
	Use: "get appId",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a appId argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(serverAddr)
		client.GetApplication(c, &pb.GetApplicationRequest{Id: args[0]})
	},
}

var subCmd = &cobra.Command{
	Use: "sub",
}

var createSubCmd = &cobra.Command{
	Use: "create appId endpoint",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires appId and endpoint arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(serverAddr)
		client.CreateSubscription(c, &pb.CreateSubscriptionRequest{ApplicationId: args[0], Endpoint: args[1]})
	},
}

var getSubCmd = &cobra.Command{
	Use: "get appId",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires appId argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(serverAddr)
		client.GetSubscriptions(c, &pb.GetSubscriptionsRequest{ApplicationId: args[0]})
	},
}

var msgCmd = &cobra.Command{
	Use: "msg",
}

var createMsgCmd = &cobra.Command{
	Use: "create appId data",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires appId and data arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(serverAddr)
		data := []byte(args[1])
		client.CreateMessage(c, &pb.CreateMessageRequest{ApplicationId: args[0], Data: data})
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
