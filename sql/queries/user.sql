-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_passwords)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE users.email = $1;