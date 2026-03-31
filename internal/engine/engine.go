package engine

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/wbhemingway/gocker/internal/db"
)

type Engine struct {
	queries *db.Queries
}

type Break struct {
	Start time.Time  `json:"start"`
	End   *time.Time `json:"end,omitempty"`
}

func NewEngine(queries *db.Queries) *Engine {
	return &Engine{queries: queries}
}

func (e *Engine) StartTask(name string, rate int64) error {

	_, err := e.queries.GetActiveEntry(context.Background())
	if err == nil {
		return fmt.Errorf("cannot start '%s': another task is already running", name)
	}

	if err != sql.ErrNoRows {
		return fmt.Errorf("database error checking active tasks: %w", err)
	}

	args := db.CreateEntryParams{
		TaskName:   name,
		HourlyRate: rate,
		StartTime:  time.Now(),
	}

	_, err = e.queries.CreateEntry(context.Background(), args)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) StopTask() error {
	curTask, err := e.queries.GetActiveEntry(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no active task to stop/cancel")
		}
		return fmt.Errorf("database error: %w", err)
	}

	args := db.EndEntryParams{
		EndTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: curTask.ID,
	}
	_, err = e.queries.EndEntry(context.Background(), args)
	if err != nil {
		return fmt.Errorf("cannot stop task: %w", err)
	}
	return nil
}

func (e *Engine) CancelTask() error {
	curTask, err := e.queries.GetActiveEntry(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no active task to stop/cancel")
		}
		return fmt.Errorf("database error: %w", err)
	}

	err = e.queries.CancelEntry(context.Background(), curTask.ID)
	if err != nil {
		return fmt.Errorf("cannot cancel task: %w", err)
	}
	return nil
}

func (e *Engine) ToggleBreak() error {
	curTask, err := e.queries.GetActiveEntry(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no active task to take a break from")
		}
		return fmt.Errorf("database error: %w", err)
	}

	var breaks []Break

	err = json.Unmarshal([]byte(curTask.BreaksJson), &breaks)
	if err != nil {
		return fmt.Errorf("failed to parse breaks: %w", err)
	}

	if len(breaks) == 0 || breaks[len(breaks)-1].End != nil {
		breaks = append(breaks, Break{Start: time.Now()})
	} else {
		now := time.Now()
		breaks[len(breaks)-1].End = &now
	}

	data, err := json.Marshal(breaks)
	if err != nil {
		return fmt.Errorf("failed to encode breaks: %w", err)
	}

	args := db.UpdateEntryBreaksParams{
		BreaksJson: string(data),
		ID:         curTask.ID,
	}
	err = e.queries.UpdateEntryBreaks(context.Background(), args)
	if err != nil {
		return fmt.Errorf("failed to update breaks: %w", err)
	}

	return nil
}
