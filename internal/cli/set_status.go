package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/wbhemingway/gocker/internal/models"
)

var setStatusCmd = &cobra.Command{
	Use:   "set-status [entry ID] [new status]",
	Short: "Update the billing status of a specific entry",
	Long:  "Valid statuses are: active, completed, billed, paid, cancelled",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		entryID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatalf("Invalid entry ID: %v", err)
		}

		rawStatusString := args[1]

		validStatus, err := models.ParseStatus(rawStatusString)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		err = tracker.UpdateStatus(entryID, validStatus)
		if err != nil {
			log.Fatalf("Database error: %v", err)
		}

		fmt.Printf("Successfully updated entry %d to '%s'\n", entryID, validStatus)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
