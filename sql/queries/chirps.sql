-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
WHERE ($1::uuid IS NULL OR $1 = chirps.user_id)
ORDER BY chirps.created_at asc;

-- name: GetChirp :one
SELECT * FROM chirps
Where chirps.id = $1
LIMIT 1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1 and user_id = $2;