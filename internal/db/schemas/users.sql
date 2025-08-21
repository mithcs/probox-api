CREATE TABLE users (
    id       BIGSERIAL   PRIMARY KEY,
    username VARCHAR(16) UNIQUE NOT NULL,
    password BINARY(60)  NOT NULL
);
