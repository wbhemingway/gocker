package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop current tracking session",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := tracker.StopTask()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		fmt.Printf("Stopped task\n")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
