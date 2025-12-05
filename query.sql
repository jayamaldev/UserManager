-- name: GetUser :one
SELECT * FROM users
WHERE userId = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY firstName;

-- name: CreateUser :one
INSERT INTO users (
  firstName, lastName, email, phone, age, user_status
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
  set 
  firstName = $2,
  lastName = $3,
  email = $4,
  phone = $5,
  age = $6,
  user_status = $7
WHERE userId = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE userId = $1;