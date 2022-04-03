package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/belljustin/captainhook/internal/workers"
)

var (
	redisAddr string
)

var workersCmd = &cobra.Command{
	Use:   "workers",
	Short: "Start instance of captainhook workers",
	Run: func(cmd *cobra.Command, args []string) {
		w := workers.New(redisAddr)
		if err := w.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	workersCmd.Flags().StringVar(&redisAddr, "redisAddr", "localhost:6379", "Address of redis supporting asynq")
	rootCmd.AddCommand(workersCmd)
}
