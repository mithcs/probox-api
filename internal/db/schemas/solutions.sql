CREATE TABLE solutions (
    id        BIGSERIAL PRIMARY KEY,
    problemId BIGINT    NOT NULL CHECK (problemId > 0),
    solution  TEXT      NOT NULL
);
