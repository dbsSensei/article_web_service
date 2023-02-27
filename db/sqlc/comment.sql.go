// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: comment.sql

package db

import (
	"context"
	"time"
)

const createComment = `-- name: CreateComment :one
INSERT INTO "comments" ("article_id", "user_id", "comment_date", "content")
VALUES ($1, $2, $3, $4) RETURNING id, article_id, user_id, comment_date, content
`

type CreateCommentParams struct {
	ArticleID   int32     `json:"article_id"`
	UserID      int32     `json:"user_id"`
	CommentDate time.Time `json:"comment_date"`
	Content     string    `json:"content"`
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, createComment,
		arg.ArticleID,
		arg.UserID,
		arg.CommentDate,
		arg.Content,
	)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.ArticleID,
		&i.UserID,
		&i.CommentDate,
		&i.Content,
	)
	return i, err
}
