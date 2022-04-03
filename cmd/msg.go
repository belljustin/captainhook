package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/belljustin/captainhook/internal/client"
	pb "github.com/belljustin/captainhook/proto/captainhook"
)

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

func init() {
	msgCmd.AddCommand(createMsgCmd)

	rootCmd.AddCommand(msgCmd)
}
