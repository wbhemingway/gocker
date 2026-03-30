-- +goose Up
CREATE TABLE entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_name TEXT NOT NULL,
    hourly_rate REAL NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    status TEXT NOT NULL DEFAULT 'active',
    breaks_json TEXT DEFAULT '[]',
    note TEXT
);

-- +goose Down
DROP TABLE entries;
