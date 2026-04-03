package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var breakCmd = &cobra.Command{
	Use:   "break",
	Short: "Toggle break status for current tracking session",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := tracker.ToggleBreak()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		fmt.Printf("Toggled break status\n")
	},
}

func init() {
	rootCmd.AddCommand(breakCmd)
}
