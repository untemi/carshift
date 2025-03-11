--    Users     --
-- name: CreateUser :one
INSERT INTO users (
  username, firstname, lastname, passhash, phone, email 
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING id;

-- name: IsUsernameUsed :one
SELECT COUNT(id) FROM users
  WHERE username = ?
LIMIT 1;

-- name: FindUserById :one
SELECT * FROM users
  WHERE id = ?
LIMIT 1;

-- name: FindUserByUsername :one
SELECT * FROM users
  WHERE username = ?
LIMIT 1;

-- name: QueryUsers :many
SELECT * FROM users
  WHERE username LIKE ?
LIMIT ? OFFSET ?;

--    Cars     --
-- name: CreateCar :one
INSERT INTO cars (
  name, price, start_at, end_at, owner_id, district_id 
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING id;

-- name: QueryCars :many
SELECT c.* FROM cars c
  JOIN districts d ON c.district_id = d.id
WHERE c.name LIKE ?1
  AND d.name = @district_name
  AND (
    (start_at <= ?2 AND end_at >= ?3) 
    OR (?2 IS NULL AND ?3 IS NULL)
  )
LIMIT ?4 OFFSET ?5;
