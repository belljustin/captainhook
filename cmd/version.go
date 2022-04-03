package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Captainhook",
	Long:  `All software has versions. This is Captainhook's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Captainhook webhook delivery service v0.1 -- HEAD")
	},
}

func init() {
	// captainhook version
	rootCmd.AddCommand(versionCmd)
}
