-- +goose Up
CREATE TABLE feeds(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT,
    url TEXT,
    user_id UUID,
    CONSTRAINT fk_feed_user_id
     FOREIGN KEY (user_id)
     REFERENCES users (id)
     ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;