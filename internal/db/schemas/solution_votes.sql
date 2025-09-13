CREATE TABLE solution_votes (
    solution_id BIGINT              NOT NULL,
    voter_id    BIGINT              NOT NULL,
    type        ENUM ('up', 'down') NOT NULL,

    PRIMARY KEY (solution_id, voter_id)
);
