-- name: FindRoleByID :one
SELECT * FROM roles r WHERE r.id = $1;

-- name: RoleExistsByID :one
SELECT EXISTS(SELECT 1 FROM roles r WHERE r.id = $1);
