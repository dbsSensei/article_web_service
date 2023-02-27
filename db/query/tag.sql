-- name: CreateTag :one
INSERT INTO "tags" ("name")
VALUES ($1) RETURNING *;

-- name: FindTag :one
SELECT * FROM "tags" WHERE id = $1;

-- name: CountTag :one
SELECT COUNT(*) FROM "tags";