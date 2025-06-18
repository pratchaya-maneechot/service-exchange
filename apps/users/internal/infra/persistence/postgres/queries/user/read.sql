-- name: FindUserByID :one
SELECT
    u.id, u.line_user_id, u.email, u.password_hash, u.status, u.created_at, u.updated_at, u.last_login_at,
    p.display_name, p.first_name, p.last_name, p.bio, p.avatar_url, p.phone_number, p.address, p.preferences
FROM users u
LEFT JOIN profiles p ON u.id = p.user_id
WHERE u.id = $1;

-- name: FindUserByLineUserID :one
SELECT
    u.id, u.line_user_id, u.email, u.password_hash, u.status, u.created_at, u.updated_at, u.last_login_at,
    p.display_name, p.first_name, p.last_name, p.bio, p.avatar_url, p.phone_number, p.address, p.preferences
FROM users u
LEFT JOIN profiles p ON u.id = p.user_id
WHERE u.line_user_id = $1;

-- name: UserExistsByLineUserID :one
SELECT EXISTS(SELECT 1 FROM users WHERE line_user_id = $1);

-- name: UserExistsByID :one
SELECT EXISTS(SELECT 1 FROM users u WHERE u.id = $1);

-- name: ProfileExistsByUserID :one
SELECT EXISTS(SELECT 1 FROM profiles p WHERE p.user_id = $1);

-- name: UserRoleExists :one
SELECT EXISTS(SELECT 1 FROM user_roles ur WHERE ur.user_id = $1 AND ur.role_id = $2);

-- name: GetUserRoles :many
SELECT
    r.id, r.name, r.description
FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1;