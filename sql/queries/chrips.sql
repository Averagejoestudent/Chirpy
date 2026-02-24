-- name: UserChirp :one
INSERT INTO chirps (id, created_at, updated_at, body , user_id)
VALUES (
   gen_random_uuid(),
   NOW(),
   NOW(),
   $1,
   $2
)
RETURNING *;


-- name: GetAllChirps :many
SELECT * FROM chirps
Order by created_at asc;

-- name: GetChirpsByID :one
SELECT * FROM chirps
Where id = $1;

-- name: DelChirpsByID :exec
DELETE FROM chirps
Where id = $1 AND user_id = $2;