package models

import "time"

type TaskStatus struct {
	TaskName      string
	IsOnBreak     bool
	TotalDuration time.Duration
	PaidDuration  time.Duration
}
