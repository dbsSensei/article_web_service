-- name: CreateArticle :one
INSERT INTO "articles" ("title", "author_id", "content", "published_at")
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetArticles :many
SELECT "articles"."id", "articles"."title", "articles"."author_id", "articles"."content", "articles"."published_at",
       "users".FULL_NAME AS author_name,
       COUNT("views".*) AS total_views
FROM "articles"
         LEFT JOIN "views" ON "views".ARTICLE_ID = "articles".ID
         LEFT JOIN "users" ON "users".ID = "articles".AUTHOR_ID
GROUP BY "articles".ID,
         "users".FULL_NAME
HAVING ("users".FULL_NAME ILIKE CONCAT('%%',@author_name::text,'%%') OR @author_name = '')
AND (("articles".TITLE ILIKE CONCAT('%%',@search::text,'%%') OR @search = '')
					OR ("articles".CONTENT ILIKE CONCAT('%%',@search::text,'%%') OR @search = ''))
ORDER BY "articles"."id" DESC;

-- name: GetArticle :one
SELECT *
FROM "articles"
WHERE id = $1;