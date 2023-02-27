-- name: CreateImage :one
INSERT INTO "images" ("article_id", "url")
VALUES ($1, $2) RETURNING *;

-- name: GetImageByArticleId :many
SELECT *
FROM "images"
WHERE article_id = $1;