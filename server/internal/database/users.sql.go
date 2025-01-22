// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const checkUser = `-- name: CheckUser :one
SELECT id, email FROM users WHERE email = $1
`

type CheckUserRow struct {
	ID    uuid.UUID
	Email string
}

func (q *Queries) CheckUser(ctx context.Context, email string) (CheckUserRow, error) {
	row := q.db.QueryRowContext(ctx, checkUser, email)
	var i CheckUserRow
	err := row.Scan(&i.ID, &i.Email)
	return i, err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (id, email, username, hashedPassword, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateUserParams struct {
	ID             uuid.UUID
	Email          string
	Username       string
	Hashedpassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Username,
		arg.Hashedpassword,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, username, hashedpassword, created_at, updated_at FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Hashedpassword,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
