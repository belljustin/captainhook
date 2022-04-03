package cmd

import (
	"errors"
	"github.com/belljustin/captainhook/cmd/client"
	pb "github.com/belljustin/captainhook/proto/captainhook"
	"github.com/spf13/cobra"
)

var (
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

func init() {
	createApplicationCmd.Flags().StringVarP(&appName, "name", "n", "", "Human readable name that will be associated with the application")
	applicationCmd.AddCommand(createApplicationCmd)

	applicationCmd.AddCommand(getApplicationCmd)

	rootCmd.AddCommand(applicationCmd)
}
