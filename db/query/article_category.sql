-- name: CreateArticleCategory :one
INSERT INTO "article_categories" ("article_id", "category_id")
VALUES ($1, $2)
RETURNING *;

-- name: GetArticleCategoryByArticleId :many
SELECT * FROM "article_categories"
JOIN "categories" ON "article_categories".category_id = "categories".id
WHERE "article_categories".article_id = $1;