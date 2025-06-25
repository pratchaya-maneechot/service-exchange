-- name: CreateUserRole :one
INSERT INTO user_roles (
    user_id,
    role_id
) VALUES (
    $1,
    $2
) RETURNING user_id, role_id;

-- name: DeleteUserRoles :exec
DELETE FROM user_roles
WHERE user_id = $1;