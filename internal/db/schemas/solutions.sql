CREATE TABLE solutions (
    id           BIGSERIAL    PRIMARY KEY,
    problemId    BIGINT       NOT NULL CHECK (problemId > 0),
    title        VARCHAR(120) NOT NULL,
    description  TEXT         NOT NULL
);
