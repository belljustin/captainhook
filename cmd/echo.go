package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/belljustin/captainhook/internal/echo"
)

var (
	echoPort int
)

var echoCmd = &cobra.Command{
	Use:   "echo",
	Short: "Start a server that prints every webhook received.",
	Run: func(cmd *cobra.Command, args []string) {
		w := echo.New(echoPort)
		if err := w.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	echoCmd.Flags().IntVar(&echoPort, "port", 8081, "Echo server port")
	rootCmd.AddCommand(echoCmd)
}
