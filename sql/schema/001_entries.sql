-- +goose Up
CREATE TABLE entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_name TEXT NOT NULL,
    hourly_rate INTEGER NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    status TEXT NOT NULL DEFAULT 'active',
    breaks_json TEXT NOT NULL DEFAULT '[]',
    note TEXT NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE entries;
