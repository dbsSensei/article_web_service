-- name: CreateLike :one
INSERT INTO "likes" ("article_id", "user_id")
VALUES ($1, $2) RETURNING *;
