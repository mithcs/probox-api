CREATE TABLE problems (
    id                   BIGSERIAL    PRIMARY KEY,
    title                VARCHAR(120) NOT NULL,
    description          TEXT         NOT NULL,
    owner_id             BIGINT       NOT NULL,
    owner_name           VARCHAR(24)  NOT NULL,
    accepted_solution_id BIGINT,
    created_at           TIMESTAMP    NOT NULL DEFAULT now(),

    FOREIGN KEY (owner_id, owner_name) REFERENCES users(id, full_name)
);
