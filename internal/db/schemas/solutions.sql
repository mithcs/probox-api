CREATE TABLE solutions (
    id           BIGSERIAL    PRIMARY KEY,
    problem_id   BIGINT       NOT NULL CHECK (problem_id > 0),
    title        VARCHAR(120) NOT NULL,
    description  TEXT         NOT NULL,
    owner_id     BIGINT       NOT NULL,
    owner_name   VARCHAR(24)  NOT NULL,

    FOREIGN KEY (problem_id) REFERENCES problems(id),
    FOREIGN KEY (owner_id, owner_name) REFERENCES users(id, full_name)
);
