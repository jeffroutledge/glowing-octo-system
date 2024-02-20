-- +goose Up
CREATE TABLE feeds(
    id UUID,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT,
    url TEXT,
    user_id UUID REFERENCES users (id)
);

-- +goose Down
DELETE TABLE feeds
ON CASCADE
