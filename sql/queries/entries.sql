-- name: CreateEntry :one
INSERT INTO entries (
task_name,
hourly_rate,
start_time,
note,
status
) VALUES (
 ?, ?, ?, ?, 'active'
 )
RETURNING *;

-- name: CreateFlatEntry :one
INSERT INTO entries (
    task_name,
    start_time,
    end_time,
    note,
    status,
    flat_fee
) VALUES (
    @task_name,
    @logged_time,
    @logged_time,
    @note,
    'completed',
    @flat_fee
)
RETURNING *;

-- name: GetActiveEntry :one
SELECT * FROM entries
WHERE status = 'active'
LIMIT 1;

-- name: EndEntry :one
UPDATE entries
SET
    end_time = ?,
    status = 'completed'
WHERE id = ?
RETURNING *;

-- name: CancelEntry :exec
UPDATE entries
SET status = 'cancelled'
WHERE id = ?;

-- name: ListRecentEntries :many
SELECT * FROM entries
ORDER BY start_time DESC
LIMIT ?;


-- name: UpdateEntryBreaks :exec
UPDATE entries
SET breaks_json = ?
WHERE id = ?;

-- name: UpdateEntryStatus :exec
UPDATE entries
SET status = ?
WHERE id = ?;
