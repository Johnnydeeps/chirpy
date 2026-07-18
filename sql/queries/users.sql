-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: ResetAllUsers :exec
delete from users;

-- name: GetUserByEmailLogin :one
select *
from users
where email = $1;
