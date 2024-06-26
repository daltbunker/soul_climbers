// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: reset_token.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createResetToken = `-- name: CreateResetToken :one
INSERT INTO reset_token(token, expiration, email)
VALUES($1, $2, $3)
RETURNING token, expiration, email
`

type CreateResetTokenParams struct {
	Token      uuid.UUID
	Expiration time.Time
	Email      string
}

func (q *Queries) CreateResetToken(ctx context.Context, arg CreateResetTokenParams) (ResetToken, error) {
	row := q.db.QueryRowContext(ctx, createResetToken, arg.Token, arg.Expiration, arg.Email)
	var i ResetToken
	err := row.Scan(&i.Token, &i.Expiration, &i.Email)
	return i, err
}

const getResetTokenByToken = `-- name: GetResetTokenByToken :one
SELECT token, expiration, email
FROM reset_token 
WHERE token = $1
`

func (q *Queries) GetResetTokenByToken(ctx context.Context, token uuid.UUID) (ResetToken, error) {
	row := q.db.QueryRowContext(ctx, getResetTokenByToken, token)
	var i ResetToken
	err := row.Scan(&i.Token, &i.Expiration, &i.Email)
	return i, err
}
