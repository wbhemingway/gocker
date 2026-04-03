package cli

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of current tracking session",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := tracker.GetStatus()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		breakStr := ""
		if status.IsOnBreak {
			breakStr = " (ON BREAK)"
		}

		fmt.Printf("Task: %s%s\n", status.TaskName, breakStr)
		fmt.Printf("Total Duration: %s\n", status.TotalDuration.Round(time.Second))
		fmt.Printf("Paid Duration: %s\n", status.PaidDuration.Round(time.Second))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
