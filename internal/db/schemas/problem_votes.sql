CREATE TABLE problem_votes (
    problem_id BIGINT              NOT NULL,
    voter_id   BIGINT              NOT NULL,
    type       ENUM ('up', 'down') NOT NULL,

    PRIMARY KEY (problem_id, voter_id)
);
