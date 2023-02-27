-- name: CreateCategory :one
INSERT INTO "categories" ("name")
VALUES ($1) RETURNING *;

-- name: FindCategory :one
SELECT *
FROM "categories"
WHERE id = $1;

-- name: CountCategory :one
SELECT COUNT(*) FROM "categories";