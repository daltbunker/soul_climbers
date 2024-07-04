-- +goose Up
CREATE TABLE blog (
    blog_id SERIAL PRIMARY KEY,
    body BYTEA NOT NULL, 
    title TEXT UNIQUE NOT NULL,
    created_by INT NOT NULL REFERENCES users(users_id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE blog_img (
    img_name TEXT NOT NULL,
    img BYTEA NOT NULL, 
    blog_id INT NOT NULL REFERENCES blog(blog_id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY(img_name, blog_id)
);

-- +goose Down
DROP TABLE blogs;
DROP TABLE blog_img;