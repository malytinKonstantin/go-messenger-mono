-- name: CreateUser :one
INSERT INTO user_credentials (email, password_hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM user_credentials
WHERE email = $1 LIMIT 1;

-- name: UpdateUserVerificationStatus :exec
UPDATE user_credentials
SET is_verified = $1, updated_at = NOW()
WHERE user_id = $2;

-- name: CreateOAuthAccount :one
INSERT INTO oauth_accounts (user_id, provider, provider_user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOAuthAccount :one
SELECT * FROM oauth_accounts
WHERE provider = $1 AND provider_user_id = $2 LIMIT 1;

-- name: CreateResetPasswordToken :one
INSERT INTO reset_password_tokens (user_id, expires_at)
VALUES ($1, $2)
RETURNING *;

-- name: GetResetPasswordToken :one
SELECT * FROM reset_password_tokens
WHERE token = $1 AND expires_at > NOW() LIMIT 1;

-- name: DeleteResetPasswordToken :exec
DELETE FROM reset_password_tokens
WHERE token = $1;

-- name: UpdateUserPassword :exec
UPDATE user_credentials
SET password_hash = $1, updated_at = NOW()
WHERE user_id = $2;