package engine

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wbhemingway/gocker/internal/db"
	"github.com/wbhemingway/gocker/internal/models"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) (*Engine, *sql.DB) {
	testDB, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	schemaPath := filepath.Join("..", "..", "sql", "schema")
	files, err := os.ReadDir(schemaPath)
	if err != nil {
		t.Fatalf("failed to read schema directory: %v", err)
	}

	for _, file := range files {
		content, err := os.ReadFile(filepath.Join(schemaPath, file.Name()))
		if err != nil {
			t.Fatalf("failed to read schema file %s: %v", file.Name(), err)
		}

		upQuery := strings.Split(string(content), "-- +goose Down")[0]

		_, err = testDB.Exec(upQuery)
		if err != nil {
			t.Fatalf("failed to execute schema %s: %v", file.Name(), err)
		}
	}

	queries := db.New(testDB)
	return NewEngine(queries, testDB), testDB
}

func TestEngine_StartTask(t *testing.T) {
	tests := []struct {
		name        string
		taskName    string
		rate        int64
		note        string
		tags        []string
		setup       func(eng *Engine)
		expectError bool
	}{
		{
			name:        "Start normal task",
			taskName:    "Test task",
			rate:        1000,
			note:        "note",
			tags:        []string{"tag1", "tag2"},
			expectError: false,
		},
		{
			name:     "Start task when one is already running",
			taskName: "Second task",
			setup: func(eng *Engine) {
				_ = eng.StartTask("First task", 0, "", nil)
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eng, dbObj := setupTestDB(t)
			defer dbObj.Close()

			if tc.setup != nil {
				tc.setup(eng)
			}

			err := eng.StartTask(tc.taskName, tc.rate, tc.note, tc.tags)
			if (err != nil) != tc.expectError {
				t.Errorf("Expected err: %v, got: %v", tc.expectError, err)
			}

			if !tc.expectError {
				active, err := eng.queries.GetActiveEntry(context.Background())
				if err != nil {
					t.Fatalf("Expected active task, got error: %v", err)
				}
				if active.TaskName != tc.taskName {
					t.Errorf("Expected task name %q, got %q", tc.taskName, active.TaskName)
				}
				if active.HourlyRate != tc.rate {
					t.Errorf("Expected rate %d, got %d", tc.rate, active.HourlyRate)
				}
				if active.Note != tc.note {
					t.Errorf("Expected note %q, got %q", tc.note, active.Note)
				}

				if tc.tags != nil {
					verifyTags(t, eng, active.ID, tc.tags)
				}
			}
		})
	}
}

func TestEngine_CreateFlatTask(t *testing.T) {
	tests := []struct {
		name        string
		taskName    string
		fee         int64
		note        string
		tags        []string
		expectError bool
	}{
		{
			name:        "Create flat task",
			taskName:    "Flat task",
			fee:         5000,
			note:        "flat note",
			tags:        []string{"flat1"},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eng, dbObj := setupTestDB(t)
			defer dbObj.Close()

			err := eng.CreateFlatTask(tc.taskName, tc.fee, tc.note, tc.tags)
			if (err != nil) != tc.expectError {
				t.Errorf("Expected err: %v, got: %v", tc.expectError, err)
			}

			if !tc.expectError {
				entries, err := eng.queries.ListRecentEntries(context.Background(), 1)
				if err != nil {
					t.Fatalf("Failed to list entries: %v", err)
				}
				if len(entries) == 0 {
					t.Fatal("Expected an entry to be created")
				}
				entry := entries[0]
				if entry.TaskName != tc.taskName {
					t.Errorf("Expected task name %q, got %q", tc.taskName, entry.TaskName)
				}
				if entry.FlatFee != tc.fee {
					t.Errorf("Expected fee %d, got %d", tc.fee, entry.FlatFee)
				}
				if entry.Status != "completed" {
					t.Errorf("Expected status 'completed', got %q", entry.Status)
				}

				if tc.tags != nil {
					verifyTags(t, eng, entry.ID, tc.tags)
				}
			}
		})
	}
}

func verifyTags(t *testing.T, eng *Engine, entryID int64, expectedTags []string) {
	dbTags, err := eng.queries.GetTagsForEntry(context.Background(), entryID)
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}
	if len(dbTags) != len(expectedTags) {
		t.Errorf("Expected %d tags, got %d", len(expectedTags), len(dbTags))
		return
	}
	tagMap := make(map[string]bool)
	for _, t := range dbTags {
		tagMap[t.Name] = true
	}
	for _, tag := range expectedTags {
		if !tagMap[tag] {
			t.Errorf("Expected tag %q to be present", tag)
		}
	}
}

func TestEngine_StopTask(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(eng *Engine)
		expectError bool
	}{
		{
			name: "Stop active task",
			setup: func(eng *Engine) {
				_ = eng.StartTask("Task", 1000, "Note", nil)
			},
			expectError: false,
		},
		{
			name:        "Stop when no task is running",
			setup:       func(eng *Engine) {},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eng, dbObj := setupTestDB(t)
			defer dbObj.Close()

			if tc.setup != nil {
				tc.setup(eng)
			}

			err := eng.StopTask()
			if (err != nil) != tc.expectError {
				t.Errorf("Expected err: %v, got: %v", tc.expectError, err)
			}

			if !tc.expectError {
				_, err := eng.queries.GetActiveEntry(context.Background())
				if err == nil {
					t.Errorf("Expected no active task after stopping")
				}
			}
		})
	}
}

func TestEngine_CancelTask(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(eng *Engine)
		expectError bool
	}{
		{
			name: "Cancel active task",
			setup: func(eng *Engine) {
				_ = eng.StartTask("Task", 1000, "Note", nil)
			},
			expectError: false,
		},
		{
			name:        "Cancel when no task is running",
			setup:       func(eng *Engine) {},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eng, dbObj := setupTestDB(t)
			defer dbObj.Close()

			if tc.setup != nil {
				tc.setup(eng)
			}

			err := eng.CancelTask()
			if (err != nil) != tc.expectError {
				t.Errorf("Expected err: %v, got: %v", tc.expectError, err)
			}

			if !tc.expectError {
				_, err := eng.queries.GetActiveEntry(context.Background())
				if err == nil {
					t.Errorf("Expected no active task after cancelling")
				}
			}
		})
	}
}

func TestEngine_ToggleBreak(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(eng *Engine)
		expectError bool
		checkStatus func(t *testing.T, eng *Engine)
	}{
		{
			name: "Start break",
			setup: func(eng *Engine) {
				_ = eng.StartTask("Task", 1000, "Note", nil)
			},
			expectError: false,
			checkStatus: func(t *testing.T, eng *Engine) {
				status, _ := eng.GetStatus()
				if !status.IsOnBreak {
					t.Errorf("Expected task to be on break")
				}
			},
		},
		{
			name: "Stop break",
			setup: func(eng *Engine) {
				_ = eng.StartTask("Task", 1000, "Note", nil)
				_ = eng.ToggleBreak()
			},
			expectError: false,
			checkStatus: func(t *testing.T, eng *Engine) {
				status, _ := eng.GetStatus()
				if status.IsOnBreak {
					t.Errorf("Expected task to not be on break")
				}
			},
		},
		{
			name:        "Toggle break when no task is running",
			setup:       func(eng *Engine) {},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eng, dbObj := setupTestDB(t)
			defer dbObj.Close()

			if tc.setup != nil {
				tc.setup(eng)
			}

			err := eng.ToggleBreak()
			if (err != nil) != tc.expectError {
				t.Errorf("Expected err: %v, got: %v", tc.expectError, err)
			}

			if !tc.expectError && tc.checkStatus != nil {
				tc.checkStatus(t, eng)
			}
		})
	}
}

func TestEngine_GetStatus(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(eng *Engine)
		expectError bool
		verify      func(t *testing.T, status *models.TaskStatus)
	}{
		{
			name: "Status of active task",
			setup: func(eng *Engine) {
				_ = eng.StartTask("ActiveTask", 5000, "", nil)
			},
			expectError: false,
			verify: func(t *testing.T, status *models.TaskStatus) {
				if status.TaskName != "ActiveTask" {
					t.Errorf("Expected task name %q, got %q", "ActiveTask", status.TaskName)
				}
				if status.IsOnBreak {
					t.Errorf("Expected task to not be on break")
				}
				if status.TotalDuration < 0 {
					t.Errorf("Expected positive total duration")
				}
				if status.PaidDuration < 0 {
					t.Errorf("Expected positive paid duration")
				}
			},
		},
		{
			name: "Status of active task on break",
			setup: func(eng *Engine) {
				_ = eng.StartTask("BreakTask", 5000, "", nil)
				_ = eng.ToggleBreak()
			},
			expectError: false,
			verify: func(t *testing.T, status *models.TaskStatus) {
				if status.TaskName != "BreakTask" {
					t.Errorf("Expected task name %q, got %q", "BreakTask", status.TaskName)
				}
				if !status.IsOnBreak {
					t.Errorf("Expected task to be on break")
				}
			},
		},
		{
			name:        "Status when no task is running",
			setup:       func(eng *Engine) {},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			eng, dbObj := setupTestDB(t)
			defer dbObj.Close()

			if tc.setup != nil {
				tc.setup(eng)
			}

			status, err := eng.GetStatus()
			if (err != nil) != tc.expectError {
				t.Errorf("Expected err: %v, got: %v", tc.expectError, err)
			}

			if !tc.expectError && tc.verify != nil {
				tc.verify(t, status)
			}
		})
	}
}
