-- +goose Up

CREATE TABLE reset_token (
    token UUID PRIMARY KEY,
    expiration TIMESTAMP NOT NULL,
    email TEXT NOT NULL REFERENCES users(email)
);

-- +goose Down
DROP TABLE reset_token;