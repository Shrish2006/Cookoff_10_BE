-- +goose Up
CREATE TABLE submission_results (
    id UUID NOT NULL UNIQUE,
	testcase_id UUID,
    submission_id UUID NOT NULL,
    runtime DECIMAL NOT NULL,
    memory NUMERIC NOT NULL,
    points_awarded INTEGER NOT NULL,
	status TEXT NOT NULL,
    description TEXT,
    PRIMARY KEY(id),
    FOREIGN KEY(submission_id) REFERENCES submissions(id)
    ON UPDATE NO ACTION ON DELETE CASCADE
);

-- +goose Down
DROP TABLE submission_results;