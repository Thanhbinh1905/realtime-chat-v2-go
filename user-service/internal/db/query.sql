-- name: CreateUser :one
INSERT INTO users (id, email, username, avatar)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: SendFriendRequest :one
INSERT INTO friendships (id, requester_id, addressee_id, status)
VALUES ($1, $2, $3, 'pending')
RETURNING *;

-- name: AcceptFriendRequest :exec
UPDATE friendships
SET status = 'accepted'
WHERE requester_id = $1 AND addressee_id = $2;

-- name: RejectFriendRequest :exec
UPDATE friendships
SET status = 'rejected'
WHERE requester_id = $1 AND addressee_id = $2;

-- name: GetFriends :many
SELECT u.*
FROM users u
JOIN friendships f ON (f.addressee_id = u.id OR f.requester_id = u.id)
WHERE (f.requester_id = $1 OR f.addressee_id = $1)
  AND f.status = 'accepted'
  AND u.id != $1;