package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel current tracking session",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := tracker.CancelTask()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		fmt.Printf("Canceled task\n")
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)
}
