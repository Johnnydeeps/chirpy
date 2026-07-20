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

-- name: UpdateUserHashedPasswordOrEmail :one
update users
set hashed_password = $1, email = $2, updated_at =$3
where id = $4
RETURNING *;

-- name: UgradeUserRedChirp :one
update users
set is_chirpy_red = true
where id = $1
RETURNING *;
