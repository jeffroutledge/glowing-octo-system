-- +goose Up
CREATE TABLE feed_follows(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID,
    CONSTRAINT fk_follow_feed_user_id
     FOREIGN KEY (user_id)
     REFERENCES users (id)
     ON DELETE CASCADE,
    feed_id UUID,
    CONSTRAINT fk_follow_feed_feed_id
     FOREIGN KEY (feed_id)
     REFERENCES feeds (id)
     ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed_follows;