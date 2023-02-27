-- name: CreateUser :one
INSERT INTO "users" ("full_name", "email", "hashed_password")
VALUES ($1, $2, $3) RETURNING *;

-- name: GetUser :one
SELECT *
FROM "users"
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM "users"
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE "users"
SET full_name     = COALESCE(sqlc.narg(full_name), full_name),
    email    = COALESCE(sqlc.narg(email), email),
    hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password)
WHERE id = sqlc.arg(id) RETURNING *;
