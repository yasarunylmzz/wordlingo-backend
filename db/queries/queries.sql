-- User-related queries
-- name: GetUserByEmail :one
SELECT id, name, surname, username, email, is_verified
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, surname, username, email, password) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING id;

-- name: UpdateUser :exec
UPDATE users 
SET name = COALESCE($1, name),
    surname = COALESCE($2, surname),
    username = COALESCE($3, username),
    email = COALESCE($4, email)
WHERE id = $5;

-- name: UpdateUserPassword :exec
UPDATE users SET password = $1 WHERE id = $2;

-- Desk-related queries
-- name: CreateDesk :one
INSERT INTO desk (title, description, image_link, user_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: UpdateDesk :exec
UPDATE desk 
SET title = COALESCE($1,title),
    description = COALESCE($2,description) ,
    image_link = COALESCE($3,image_link)
WHERE id = $4 AND user_id=$5;

-- name: DeleteDesk :exec
DELETE FROM desk WHERE id = $1 AND user_id=$2;

-- name: GetAllDesksByUserId :many
SELECT * FROM desk WHERE user_id = $1;

-- Card-related queries
-- name: CreateCard :exec
INSERT INTO card (language_1, language_2, importance_value, desk_id) VALUES ($1, $2, $3, $4);

-- name: UpdateCard :exec
UPDATE card SET language_1 = $1, language_2 = $2 WHERE id = $3 AND desk_id = $4;

-- name: DeleteCard :exec
DELETE FROM card WHERE id = $1 AND desk_id = $2;

-- name: GetCardById :one
SELECT id, language_1, language_2 , desk_id FROM card WHERE id = $1;

-- name: GetCardsByDeskId :many
SELECT id, language_1, language_2, importance_value, desk_id FROM card WHERE desk_id = $1;

-- name: VerificationCodeCreate :one
INSERT INTO verification_codes (user_id, code, created_at, expires_at)
VALUES ($1, $2, NOW(), NOW() + INTERVAL '1 day')
RETURNING id;

-- name: IsUserVerified :one
SELECT is_verified, email FROM users WHERE id = $1;

-- name: VerifyUser :exec
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
  AND EXISTS (SELECT 1 FROM delete_verification);

-- name: GetHashPass :one
SELECT password 
FROM users 
WHERE email = $1;

-- name: GetVerificationCodeById :one
SELECT code, expires_at FROM verification_codes WHERE user_id = $1;