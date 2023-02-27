// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: article.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createArticle = `-- name: CreateArticle :one
INSERT INTO "articles" ("title", "author_id", "content", "published_at")
VALUES ($1, $2, $3, $4) RETURNING id, title, author_id, content, published_at
`

type CreateArticleParams struct {
	Title       string    `json:"title"`
	AuthorID    int32     `json:"author_id"`
	Content     string    `json:"content"`
	PublishedAt time.Time `json:"published_at"`
}

func (q *Queries) CreateArticle(ctx context.Context, arg CreateArticleParams) (Article, error) {
	row := q.db.QueryRowContext(ctx, createArticle,
		arg.Title,
		arg.AuthorID,
		arg.Content,
		arg.PublishedAt,
	)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.AuthorID,
		&i.Content,
		&i.PublishedAt,
	)
	return i, err
}

const getArticle = `-- name: GetArticle :one
SELECT id, title, author_id, content, published_at
FROM "articles"
WHERE id = $1
`

func (q *Queries) GetArticle(ctx context.Context, id int32) (Article, error) {
	row := q.db.QueryRowContext(ctx, getArticle, id)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.AuthorID,
		&i.Content,
		&i.PublishedAt,
	)
	return i, err
}

const getArticles = `-- name: GetArticles :many
SELECT "articles"."id", "articles"."title", "articles"."author_id", "articles"."content", "articles"."published_at",
       "users".FULL_NAME AS author_name,
       COUNT("views".*) AS total_views
FROM "articles"
         LEFT JOIN "views" ON "views".ARTICLE_ID = "articles".ID
         LEFT JOIN "users" ON "users".ID = "articles".AUTHOR_ID
GROUP BY "articles".ID,
         "users".FULL_NAME
HAVING ("users".FULL_NAME ILIKE CONCAT('%%',$1::text,'%%') OR $1 = '')
AND (("articles".TITLE ILIKE CONCAT('%%',$2::text,'%%') OR $2 = '')
					OR ("articles".CONTENT ILIKE CONCAT('%%',$2::text,'%%') OR $2 = ''))
ORDER BY "articles"."id" DESC
`

type GetArticlesParams struct {
	AuthorName string `json:"author_name"`
	Search     string `json:"search"`
}

type GetArticlesRow struct {
	ID          int32          `json:"id"`
	Title       string         `json:"title"`
	AuthorID    int32          `json:"author_id"`
	Content     string         `json:"content"`
	PublishedAt time.Time      `json:"published_at"`
	AuthorName  sql.NullString `json:"author_name"`
	TotalViews  int64          `json:"total_views"`
}

func (q *Queries) GetArticles(ctx context.Context, arg GetArticlesParams) ([]GetArticlesRow, error) {
	rows, err := q.db.QueryContext(ctx, getArticles, arg.AuthorName, arg.Search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetArticlesRow{}
	for rows.Next() {
		var i GetArticlesRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.AuthorID,
			&i.Content,
			&i.PublishedAt,
			&i.AuthorName,
			&i.TotalViews,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}