package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

var (
	serverAddr = "localhost:50051"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
