-- name: CreateUserRole :one
INSERT INTO user_roles (
    user_id,
    role_id
) VALUES (
    $1,
    $2
) RETURNING user_id, role_id;

-- name: DeleteUserRole :exec
-- Deletes a specific user-role association.
DELETE FROM user_roles
WHERE user_id = $1 AND role_id = $2;