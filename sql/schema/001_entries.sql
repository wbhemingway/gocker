-- +goose Up
CREATE TABLE entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_name TEXT NOT NULL,
    hourly_rate INTEGER NOT NULL DEFAULT 0,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    flat_fee INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active',
    breaks_json TEXT NOT NULL DEFAULT '[]',
    note TEXT NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE entries;
