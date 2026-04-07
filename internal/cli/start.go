package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rateInput string
var noteInput string
var tagsInput string

var startCmd = &cobra.Command{
	Use:   "start [task name]",
	Short: "Start a new tracking session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskName := args[0]
		tags := parseTags(tagsInput)
		rateInCents, err := parseRate(rateInput)
		if err != nil {
			log.Fatalf("Invalid rate flag: %v", err)
		}
		err = tracker.StartTask(taskName, rateInCents, noteInput, tags)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		fmt.Printf("Started task: '%s' at $%d.%02d/hr\n", taskName, rateInCents/100, rateInCents%100)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&rateInput, "rate", "40.00", "Hourly rate in dollars (e.g. 40.00)")
	startCmd.Flags().StringVar(&noteInput, "note", "", "Optional note for the task")
	startCmd.Flags().StringVar(&tagsInput, "tags", "", "Optional tags to add to the task (e.g. \"golang, zebra, R&R\"")
}
