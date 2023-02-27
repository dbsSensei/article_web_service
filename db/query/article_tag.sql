-- name: CreateArticleTag :one
INSERT INTO "article_tags" ("article_id", "tag_id")
VALUES ($1, $2) RETURNING *;

-- name: GetArticleTagByArticleId :many
SELECT *
FROM "article_tags"
         JOIN "tags" ON "article_tags".tag_id = "tags".id
WHERE "article_tags".article_id = $1;
