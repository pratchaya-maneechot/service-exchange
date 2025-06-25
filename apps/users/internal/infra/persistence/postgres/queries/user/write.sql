-- name: CreateUser :one
INSERT INTO users (
    id, line_user_id, email, password_hash, status, created_at, updated_at, last_login_at
) VALUES (
    $1, $2, $3, $4, $5, NOW(), NOW(), $6
) RETURNING id, line_user_id, email, password_hash, status, created_at, updated_at, last_login_at;

-- name: UpsertUserProfile :one
INSERT INTO profiles (
    user_id, display_name, first_name, last_name, bio, avatar_url, phone_number, address, preferences
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
ON CONFLICT (user_id)
DO UPDATE SET
    display_name = EXCLUDED.display_name,
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    bio = EXCLUDED.bio,
    avatar_url = EXCLUDED.avatar_url,
    phone_number = EXCLUDED.phone_number,
    address = EXCLUDED.address,
    preferences = EXCLUDED.preferences
RETURNING user_id, display_name, first_name, last_name, bio, avatar_url, phone_number, address, preferences;

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

-- name: UpdateUserLastLoginAt :execrows
UPDATE users
SET last_login_at = $2, updated_at = NOW()
WHERE id = $1;