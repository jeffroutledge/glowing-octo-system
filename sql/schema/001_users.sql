-- +goose Up
CREATE TABLE users(
    id UUID,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT
);

-- +goose Down
DROP TABLE users;
