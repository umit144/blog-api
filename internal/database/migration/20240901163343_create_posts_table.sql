-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(250) NOT NULL,
    content TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd