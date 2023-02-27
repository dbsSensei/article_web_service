-- name: CreateView :exec
INSERT INTO "views" ("article_id", "view_date")
VALUES ($1, $2) RETURNING *;

-- name: CountArticleViewsByArticleId :one
SELECT COUNT(*)
FROM "views"
WHERE "article_id" = $1;