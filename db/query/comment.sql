-- name: CreateComment :one
INSERT INTO "comments" ("article_id", "user_id", "comment_date", "content")
VALUES ($1, $2, $3, $4) RETURNING *;
