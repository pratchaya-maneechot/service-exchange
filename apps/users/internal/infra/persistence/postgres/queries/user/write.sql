-- name: CreateUser :one
INSERT INTO users (
    id, line_user_id, email, password_hash, status, created_at, updated_at, last_login_at
) VALUES (
    $1, $2, $3, $4, $5, NOW(), NOW(), $6
) RETURNING id, line_user_id, email, password_hash, status, created_at, updated_at, last_login_at;

-- name: CreateProfile :one
INSERT INTO profiles (
    user_id, display_name, first_name, last_name, bio, avatar_url, phone_number, address, preferences
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING user_id, display_name, first_name, last_name, bio, avatar_url, phone_number, address, preferences;

-- name: UpdateUser :one
UPDATE users
SET
    line_user_id = $2,
    email = $3,
    password_hash = $4,
    status = $5,
    updated_at = NOW(),
    last_login_at = $6
WHERE id = $1
RETURNING id, line_user_id, email, password_hash, status, created_at, updated_at, last_login_at;

-- name: UpdateProfile :one
UPDATE profiles
SET
    display_name = $2,
    first_name = $3,
    last_name = $4,
    bio = $5,
    avatar_url = $6,
    phone_number = $7,
    address = $8,
    preferences = $9
WHERE user_id = $1
RETURNING user_id, display_name, first_name, last_name, bio, avatar_url, phone_number, address, preferences;

-- name: UpdateUserLastLoginAt :execrows
UPDATE users
SET last_login_at = $2, updated_at = NOW()
WHERE id = $1;