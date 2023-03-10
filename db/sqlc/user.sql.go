// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO "users" ("full_name", "email", "hashed_password")
VALUES ($1, $2, $3) RETURNING id, full_name, email, hashed_password
`

type CreateUserParams struct {
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.FullName, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, full_name, email, hashed_password
FROM "users"
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, full_name, email, hashed_password
FROM "users"
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE "users"
SET full_name     = COALESCE($1, full_name),
    email    = COALESCE($2, email),
    hashed_password = COALESCE($3, hashed_password)
WHERE id = $4 RETURNING id, full_name, email, hashed_password
`

type UpdateUserParams struct {
	FullName       sql.NullString `json:"full_name"`
	Email          sql.NullString `json:"email"`
	HashedPassword sql.NullString `json:"hashed_password"`
	ID             int32          `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.FullName,
		arg.Email,
		arg.HashedPassword,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.HashedPassword,
	)
	return i, err
}
