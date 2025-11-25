-- +goose Up
CREATE TABLE ascent(
    PRIMARY KEY (climb_id, created_by),
    climb_id INT NOT NULL REFERENCES climb(climb_id),
    grade TEXT NOT NULL,
    rating TEXT NOT NULL,
    attempts TEXT NOT NULL,
    over_200_pounds BOOLEAN NOT NULL,
    comment TEXT,
    ascent_date TIMESTAMP NOT NULL,
    created_by INT NOT NULL REFERENCES users(users_id),
    created_at TIMESTAMP NOT NULL DEFAULT current_date,
    updated_at TIMESTAMP NOT NULL DEFAULT current_date
);

-- +goose Down 
DROP TABLE ascent;