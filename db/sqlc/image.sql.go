// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: image.sql

package db

import (
	"context"
)

const createImage = `-- name: CreateImage :one
INSERT INTO "images" ("article_id", "url")
VALUES ($1, $2) RETURNING id, article_id, url
`

type CreateImageParams struct {
	ArticleID int32  `json:"article_id"`
	Url       string `json:"url"`
}

func (q *Queries) CreateImage(ctx context.Context, arg CreateImageParams) (Image, error) {
	row := q.db.QueryRowContext(ctx, createImage, arg.ArticleID, arg.Url)
	var i Image
	err := row.Scan(&i.ID, &i.ArticleID, &i.Url)
	return i, err
}

const getImageByArticleId = `-- name: GetImageByArticleId :many
SELECT id, article_id, url
FROM "images"
WHERE article_id = $1
`

func (q *Queries) GetImageByArticleId(ctx context.Context, articleID int32) ([]Image, error) {
	rows, err := q.db.QueryContext(ctx, getImageByArticleId, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Image{}
	for rows.Next() {
		var i Image
		if err := rows.Scan(&i.ID, &i.ArticleID, &i.Url); err != nil {
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
