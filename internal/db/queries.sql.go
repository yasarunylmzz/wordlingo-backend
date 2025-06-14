// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createCard = `-- name: CreateCard :exec
INSERT INTO card (language_1, language_2, importance_value, desk_id) VALUES ($1, $2, $3, $4)
`

type CreateCardParams struct {
	Language1       string
	Language2       string
	ImportanceValue int32
	DeskID          uuid.UUID
}

// Card-related queries
func (q *Queries) CreateCard(ctx context.Context, arg CreateCardParams) error {
	_, err := q.db.ExecContext(ctx, createCard,
		arg.Language1,
		arg.Language2,
		arg.ImportanceValue,
		arg.DeskID,
	)
	return err
}

const createDesk = `-- name: CreateDesk :one
INSERT INTO desk (title, description, image_link, user_id)
VALUES ($1, $2, $3, $4)
RETURNING id
`

type CreateDeskParams struct {
	Title       string
	Description string
	ImageLink   sql.NullString
	UserID      uuid.UUID
}

// Desk-related queries
func (q *Queries) CreateDesk(ctx context.Context, arg CreateDeskParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createDesk,
		arg.Title,
		arg.Description,
		arg.ImageLink,
		arg.UserID,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, surname, username, email, password) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING id
`

type CreateUserParams struct {
	Name     string
	Surname  string
	Username string
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Name,
		arg.Surname,
		arg.Username,
		arg.Email,
		arg.Password,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteCard = `-- name: DeleteCard :exec
DELETE FROM card WHERE id = $1 AND desk_id = $2
`

type DeleteCardParams struct {
	ID     uuid.UUID
	DeskID uuid.UUID
}

func (q *Queries) DeleteCard(ctx context.Context, arg DeleteCardParams) error {
	_, err := q.db.ExecContext(ctx, deleteCard, arg.ID, arg.DeskID)
	return err
}

const deleteDesk = `-- name: DeleteDesk :exec
DELETE FROM desk WHERE id = $1 AND user_id=$2
`

type DeleteDeskParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteDesk(ctx context.Context, arg DeleteDeskParams) error {
	_, err := q.db.ExecContext(ctx, deleteDesk, arg.ID, arg.UserID)
	return err
}

const getAllDesksByUserId = `-- name: GetAllDesksByUserId :many
SELECT id, user_id, image_link, title, description FROM desk WHERE user_id = $1
`

func (q *Queries) GetAllDesksByUserId(ctx context.Context, userID uuid.UUID) ([]Desk, error) {
	rows, err := q.db.QueryContext(ctx, getAllDesksByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Desk
	for rows.Next() {
		var i Desk
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ImageLink,
			&i.Title,
			&i.Description,
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

const getCardById = `-- name: GetCardById :one
SELECT id, language_1, language_2 , desk_id FROM card WHERE id = $1
`

type GetCardByIdRow struct {
	ID        uuid.UUID
	Language1 string
	Language2 string
	DeskID    uuid.UUID
}

func (q *Queries) GetCardById(ctx context.Context, id uuid.UUID) (GetCardByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getCardById, id)
	var i GetCardByIdRow
	err := row.Scan(
		&i.ID,
		&i.Language1,
		&i.Language2,
		&i.DeskID,
	)
	return i, err
}

const getCardsByDeskId = `-- name: GetCardsByDeskId :many
SELECT id, language_1, language_2, importance_value, desk_id FROM card WHERE desk_id = $1
`

type GetCardsByDeskIdRow struct {
	ID              uuid.UUID
	Language1       string
	Language2       string
	ImportanceValue int32
	DeskID          uuid.UUID
}

func (q *Queries) GetCardsByDeskId(ctx context.Context, deskID uuid.UUID) ([]GetCardsByDeskIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getCardsByDeskId, deskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCardsByDeskIdRow
	for rows.Next() {
		var i GetCardsByDeskIdRow
		if err := rows.Scan(
			&i.ID,
			&i.Language1,
			&i.Language2,
			&i.ImportanceValue,
			&i.DeskID,
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

const getHashPass = `-- name: GetHashPass :one
SELECT password 
FROM users 
WHERE email = $1
`

func (q *Queries) GetHashPass(ctx context.Context, email string) (string, error) {
	row := q.db.QueryRowContext(ctx, getHashPass, email)
	var password string
	err := row.Scan(&password)
	return password, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, surname, username, email, is_verified
FROM users
WHERE email = $1
`

type GetUserByEmailRow struct {
	ID         uuid.UUID
	Name       string
	Surname    string
	Username   string
	Email      string
	IsVerified sql.NullBool
}

// User-related queries
func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Surname,
		&i.Username,
		&i.Email,
		&i.IsVerified,
	)
	return i, err
}

const getVerificationCodeById = `-- name: GetVerificationCodeById :one
SELECT code, expires_at FROM verification_codes WHERE user_id = $1
`

type GetVerificationCodeByIdRow struct {
	Code      string
	ExpiresAt sql.NullTime
}

func (q *Queries) GetVerificationCodeById(ctx context.Context, userID uuid.NullUUID) (GetVerificationCodeByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getVerificationCodeById, userID)
	var i GetVerificationCodeByIdRow
	err := row.Scan(&i.Code, &i.ExpiresAt)
	return i, err
}

const isUserVerified = `-- name: IsUserVerified :one
SELECT is_verified, email FROM users WHERE id = $1
`

type IsUserVerifiedRow struct {
	IsVerified sql.NullBool
	Email      string
}

func (q *Queries) IsUserVerified(ctx context.Context, id uuid.UUID) (IsUserVerifiedRow, error) {
	row := q.db.QueryRowContext(ctx, isUserVerified, id)
	var i IsUserVerifiedRow
	err := row.Scan(&i.IsVerified, &i.Email)
	return i, err
}

const updateCard = `-- name: UpdateCard :exec
UPDATE card SET language_1 = $1, language_2 = $2 WHERE id = $3 AND desk_id = $4
`

type UpdateCardParams struct {
	Language1 string
	Language2 string
	ID        uuid.UUID
	DeskID    uuid.UUID
}

func (q *Queries) UpdateCard(ctx context.Context, arg UpdateCardParams) error {
	_, err := q.db.ExecContext(ctx, updateCard,
		arg.Language1,
		arg.Language2,
		arg.ID,
		arg.DeskID,
	)
	return err
}

const updateDesk = `-- name: UpdateDesk :exec
UPDATE desk 
SET title = COALESCE($1,title),
    description = COALESCE($2,description) ,
    image_link = COALESCE($3,image_link)
WHERE id = $4 AND user_id=$5
`

type UpdateDeskParams struct {
	Title       string
	Description string
	ImageLink   sql.NullString
	ID          uuid.UUID
	UserID      uuid.UUID
}

func (q *Queries) UpdateDesk(ctx context.Context, arg UpdateDeskParams) error {
	_, err := q.db.ExecContext(ctx, updateDesk,
		arg.Title,
		arg.Description,
		arg.ImageLink,
		arg.ID,
		arg.UserID,
	)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users 
SET name = COALESCE($1, name),
    surname = COALESCE($2, surname),
    username = COALESCE($3, username),
    email = COALESCE($4, email)
WHERE id = $5
`

type UpdateUserParams struct {
	Name     string
	Surname  string
	Username string
	Email    string
	ID       uuid.UUID
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.Name,
		arg.Surname,
		arg.Username,
		arg.Email,
		arg.ID,
	)
	return err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users SET password = $1 WHERE id = $2
`

type UpdateUserPasswordParams struct {
	Password string
	ID       uuid.UUID
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.Password, arg.ID)
	return err
}

const verificationCodeCreate = `-- name: VerificationCodeCreate :one
INSERT INTO verification_codes (user_id, code, created_at, expires_at)
VALUES ($1, $2, NOW(), NOW() + INTERVAL '1 day')
RETURNING id
`

type VerificationCodeCreateParams struct {
	UserID uuid.NullUUID
	Code   string
}

func (q *Queries) VerificationCodeCreate(ctx context.Context, arg VerificationCodeCreateParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, verificationCodeCreate, arg.UserID, arg.Code)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const verifyUser = `-- name: VerifyUser :exec
WITH delete_verification AS (
  DELETE FROM verification_codes
  WHERE user_id = $1
    AND code = $3
    AND expires_at > NOW()
  RETURNING user_id
)
UPDATE users
SET is_verified = TRUE
WHERE users.id = $1 
  AND users.email = $2 
  AND EXISTS (SELECT 1 FROM delete_verification)
`

type VerifyUserParams struct {
	ID    uuid.UUID
	Email string
	Code  string
}

func (q *Queries) VerifyUser(ctx context.Context, arg VerifyUserParams) error {
	_, err := q.db.ExecContext(ctx, verifyUser, arg.ID, arg.Email, arg.Code)
	return err
}
