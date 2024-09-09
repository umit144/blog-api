-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id              UUID PRIMARY KEY             DEFAULT uuid_generate_v4(),
    name            VARCHAR(100)        NOT NULL,
    lastname        VARCHAR(100),
    email           VARCHAR(255) UNIQUE NOT NULL,
    password        TEXT,
    google_id       VARCHAR(255) UNIQUE,
    profile_picture VARCHAR(255),
    auth_provider   VARCHAR(20)         NOT NULL DEFAULT 'local',
    created_at      TIMESTAMP WITH TIME ZONE     DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE     DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd