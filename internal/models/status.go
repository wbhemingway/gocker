package models

import (
	"fmt"
	"time"
)

type TaskStatus struct {
	TaskName      string
	IsOnBreak     bool
	TotalDuration time.Duration
	PaidDuration  time.Duration
}

type EntryStatus string

const (
	StatusActive    EntryStatus = "active"
	StatusCompleted EntryStatus = "completed"
	StatusBilled    EntryStatus = "billed"
	StatusPaid      EntryStatus = "paid"
	StatusCancelled EntryStatus = "cancelled"
)

func ParseStatus(input string) (EntryStatus, error) {
	status := EntryStatus(input)

	switch status {
	case StatusActive, StatusCompleted, StatusBilled, StatusPaid, StatusCancelled:
		return status, nil
	default:
		return "", fmt.Errorf("invalid status: '%s'. Must be active, completed, billed, paid, or cancelled", input)
	}
}
