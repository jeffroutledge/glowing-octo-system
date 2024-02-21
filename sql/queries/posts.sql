-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT p.*
FROM posts p, feed_follows ff
WHERE ff.user_id = $1
AND p.feed_id = ff.feed_id
ORDER BY p.published_at
LIMIT $2;