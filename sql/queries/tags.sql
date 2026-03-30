-- name: CreateTag :one
INSERT INTO tags (name)
VALUES (?)
ON CONFLICT(name) DO UPDATE SET name=excluded.name
RETURNING *;

-- name: LinkTagToEntry :exec
INSERT INTO entry_tags (entry_id, tag_id)
VALUES (?, ?);

-- name: UnlinkTagFromEntry :exec
DELETE FROM entry_tags
WHERE entry_id = ? AND tag_id = ?;

-- name: GetTagsForEntry :many
SELECT t.* FROM tags t
JOIN entry_tags et ON t.id = et.tag_id
WHERE et.entry_id = ?;

-- name: ListAllTags :many
SELECT * FROM tags ORDER BY name ASC;
