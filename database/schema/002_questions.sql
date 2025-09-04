-- +goose Up
CREATE TABLE questions (
	id UUID NOT NULL,
	description TEXT NOT NULL,
	title TEXT NOT NULL,
    qType TEXT NOT NULL,
    isBountyActive BOOLEAN NOT NULL DEFAULT false,
	input_format TEXT[],
	points INTEGER NOT NULL,
	round INTEGER NOT NULL,
	constraints TEXT[] NOT NULL,
	output_format TEXT[] NOT NULL,
    sample_test_input TEXT[],
    sample_test_output TEXT[],
    explanation TEXT[],
	PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE questions;