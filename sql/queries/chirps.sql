-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetAllChirps :many
select *
from chirps
order by created_at ASC;

-- name: GetChirpByID :one
select *
from chirps
where chirps.id = $1;
