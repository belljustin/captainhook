package cmd

import (
	"errors"
	"github.com/belljustin/captainhook/cmd/client"
	pb "github.com/belljustin/captainhook/proto/captainhook"
	"github.com/spf13/cobra"
)

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

func init() {
	subCmd.AddCommand(createSubCmd)

	subCmd.AddCommand(getSubCmd)

	rootCmd.AddCommand(subCmd)
}
