-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL,
    title      VARCHAR(200) NOT NULL,
    slug       VARCHAR(250) NOT NULL UNIQUE ,
    content    TEXT         NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX idx_posts_slug ON posts (slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd