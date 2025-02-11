-- +goose Up
CREATE TABLE blog (
    blog_id SERIAL PRIMARY KEY,
    body BYTEA NOT NULL, 
    title TEXT UNIQUE NOT NULL,
    excerpt TEXT NOT NULL,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_by INT NOT NULL REFERENCES users(users_id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE blog_img (
    img_name TEXT NOT NULL,
    img BYTEA NOT NULL, 
    blog_id INT NOT NULL REFERENCES blog(blog_id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY(img_name, blog_id) --TODO: don't need img_name as part of PK
);

-- +goose Down
DROP TABLE blog_img;
DROP TABLE blog;