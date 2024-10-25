-- User-related queries
-- name: GetUserByEmail :one
SELECT id, name, email, password FROM users WHERE email = $1;

-- name: CreateUser :exec
INSERT INTO users (name, email, password) VALUES ($1, $2, $3);

-- name: UpdateUserPassword :exec
UPDATE users SET password = $1 WHERE email = $2;

-- name: UpdateUserInfo :exec
UPDATE users SET name = $1, email = $2 WHERE id = $3;

-- name: SetUserVerificationCode :exec
UPDATE users SET verification_code = $1 WHERE email = $2;

-- name: VerifyUser :exec
UPDATE users SET is_verified = TRUE WHERE email = $1 AND verification_code = $2;

-- name: IsUserVerified :one
SELECT is_verified FROM users WHERE email = $1;

-- Desk-related queries
-- name: GetDeskById :one
SELECT id, title, description, user_id FROM desk WHERE id = $1;

-- name: GetDesksByUserId :many
SELECT id, title, description, user_id FROM desk WHERE user_id = $1;

-- name: CreateDesk :exec
INSERT INTO desk (title, description, user_id) VALUES ($1, $2, $3);

-- name: UpdateDesk :exec
UPDATE desk SET title = $1, description = $2 WHERE id = $3;

-- name: DeleteDesk :exec
DELETE FROM desk WHERE id = $1;

-- Card-related queries
-- name: GetCardById :one
SELECT id, language_1, language_2, description, desk_id FROM card WHERE id = $1;

-- name: GetCardsByDeskId :many
SELECT id, language_1, language_2, description, desk_id FROM card WHERE desk_id = $1;

-- name: CreateCard :exec
INSERT INTO card (language_1, language_2, description, desk_id) VALUES ($1, $2, $3, $4);

-- name: UpdateCard :exec
UPDATE card SET language_1 = $1, language_2 = $2, description = $3 WHERE id = $4;

-- name: DeleteCard :exec
DELETE FROM card WHERE id = $1;