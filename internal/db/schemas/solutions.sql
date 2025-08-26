CREATE TABLE solutions (
    id           BIGSERIAL    PRIMARY KEY,
    problem_id   BIGINT       NOT NULL CHECK (problem_id > 0),
    title        VARCHAR(120) NOT NULL,
    description  TEXT         NOT NULL
);
