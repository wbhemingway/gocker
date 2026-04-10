package engine

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/wbhemingway/gocker/internal/db"
	"github.com/wbhemingway/gocker/internal/models"
)

type Engine struct {
	queries *db.Queries
	db      *sql.DB
}

type Break struct {
	Start time.Time  `json:"start"`
	End   *time.Time `json:"end,omitempty"`
}

func NewEngine(queries *db.Queries, db *sql.DB) *Engine {
	return &Engine{queries: queries, db: db}
}

func (e *Engine) StartTask(name string, rate int64, note string, tags []string) error {

	_, err := e.queries.GetActiveEntry(context.Background())
	if err == nil {
		return fmt.Errorf("cannot start '%s': another task is already running", name)
	}

	if err != sql.ErrNoRows {
		return fmt.Errorf("database error checking active tasks: %w", err)
	}
	tx, err := e.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := e.queries.WithTx(tx)
	args := db.CreateEntryParams{
		TaskName:   name,
		HourlyRate: rate,
		Note:       note,
		StartTime:  time.Now().UTC(),
	}

	entry, err := qtx.CreateEntry(context.Background(), args)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		newTag, err := qtx.CreateTag(context.Background(), tag)
		if err != nil {
			return fmt.Errorf("cannot create tag '%s': %w", tag, err)
		}
		err = qtx.LinkTagToEntry(context.Background(), db.LinkTagToEntryParams{TagID: newTag.ID, EntryID: entry.ID})
		if err != nil {
			return fmt.Errorf("cannot add tag '%s' to entry: %w", tag, err)
		}
	}
	return tx.Commit()
}

func (e *Engine) CreateFlatTask(name string, fee int64, note string, tags []string) error {
	tx, err := e.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := e.queries.WithTx(tx)
	args := db.CreateFlatEntryParams{
		TaskName:   name,
		FlatFee:    fee,
		Note:       note,
		LoggedTime: time.Now().UTC(),
	}

	entry, err := qtx.CreateFlatEntry(context.Background(), args)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		newTag, err := qtx.CreateTag(context.Background(), tag)
		if err != nil {
			return fmt.Errorf("cannot create tag '%s': %w", tag, err)
		}
		err = qtx.LinkTagToEntry(context.Background(), db.LinkTagToEntryParams{TagID: newTag.ID, EntryID: entry.ID})
		if err != nil {
			return fmt.Errorf("cannot add tag '%s' to entry: %w", tag, err)
		}
	}
	return tx.Commit()
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
			Time:  time.Now().UTC(),
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
		breaks = append(breaks, Break{Start: time.Now().UTC()})
	} else {
		now := time.Now().UTC()
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

func (e *Engine) GetStatus() (*models.TaskStatus, error) {
	curTask, err := e.queries.GetActiveEntry(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no active task to get a status for")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	var breaks []Break
	err = json.Unmarshal([]byte(curTask.BreaksJson), &breaks)
	if err != nil {
		return nil, fmt.Errorf("failed to parse breaks: %w", err)
	}

	totalDuration := time.Since(curTask.StartTime)
	var breakDuration time.Duration
	isOnBreak := false

	for i, b := range breaks {
		if b.End != nil {
			breakDuration += b.End.Sub(b.Start)
		} else if i == len(breaks)-1 {
			isOnBreak = true
			breakDuration += time.Since(b.Start)
		}
	}

	paidDuration := totalDuration - breakDuration

	return &models.TaskStatus{
		TaskName:      curTask.TaskName,
		IsOnBreak:     isOnBreak,
		TotalDuration: totalDuration,
		PaidDuration:  paidDuration,
	}, nil
}

func (e *Engine) UpdateStatus(id int64, newStatus models.EntryStatus) error {
	args := db.UpdateEntryStatusParams{
		Status: string(newStatus),
		ID:     id,
	}

	err := e.queries.UpdateEntryStatus(context.Background(), args)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}
