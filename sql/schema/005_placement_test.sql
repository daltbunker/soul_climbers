-- +goose Up
CREATE TABLE test_question (
    test_question_id SERIAL PRIMARY KEY,
    question_text TEXT UNIQUE NOT NULL,
    input_type TEXT NOT NULL,
    answers TEXT NOT NULL, -- comma separated string
    answer_points TEXT NOT NULL, -- comma separated string
    points INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

create TABLE placement_test (
    placement_test_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(users_id),
    score INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE test_question;
DROP TABLE placement_test;