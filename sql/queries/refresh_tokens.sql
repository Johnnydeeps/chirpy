-- name: CreateRefresgToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
select users.*, refresh_tokens.expires_at, refresh_tokens.revoked_at
from refresh_tokens
inner join users on refresh_tokens.user_id = users.id
where token = $1;

-- name: RevokeRefreshToken :exec
update refresh_tokens
set revoked_at = $1, updated_at = $2
where token = $3;
