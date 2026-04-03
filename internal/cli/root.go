/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wbhemingway/gocker/internal/models"
)

type TimeTracker interface {
	StartTask(name string, rate int64, note string, tags []string) error
	StopTask() error
	CancelTask() error
	ToggleBreak() error
	GetStatus() (*models.TaskStatus, error)
}

var tracker TimeTracker

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gocker",
	Short: "A straightforward time and financial tracker for freelance work",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.

func Init(t TimeTracker) {
	tracker = t
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
