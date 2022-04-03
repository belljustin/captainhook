package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/belljustin/captainhook/internal/server"
)

var (
	port int
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start instance of captainhook workers",
	Run: func(cmd *cobra.Command, args []string) {
		w := server.New(port, redisAddr)
		if err := w.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	serverCmd.Flags().StringVar(&redisAddr, "redisAddr", "localhost:6379", "Address of redis supporting asynq")
	serverCmd.Flags().IntVar(&port, "port", 50051, "Captainhook server port")
	rootCmd.AddCommand(serverCmd)
}
