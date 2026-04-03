package cli

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var rateInput string
var noteInput string

var startCmd = &cobra.Command{
	Use:   "start [task name]",
	Short: "Start a new tracking session",
	Run: func(cmd *cobra.Command, args []string) {
		taskName := args[0]

		rateInCents, err := parseRate(rateInput)
		if err != nil {
			log.Fatalf("Invalid rate flag: %v", err)
		}
		err = tracker.StartTask(taskName, rateInCents, noteInput)
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

}

func parseRate(input string) (int64, error) {
	if input == "" {
		return 0, nil
	}

	parts := strings.Split(input, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid format: too many decimals")
	}

	dollars, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid dollar amount")
	}

	cents := int64(0)
	if len(parts) == 2 {
		centStr := parts[1]
		if len(centStr) > 2 {
			return 0, fmt.Errorf("cannot have more than two decimal places")
		}
		if len(centStr) == 1 {
			centStr += "0"
		}

		cents, err = strconv.ParseInt(centStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid cent amount")
		}
	}

	return (dollars * 100) + cents, nil
}
