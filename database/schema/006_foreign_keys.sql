-- +goose Up
ALTER TABLE submissions 
ADD CONSTRAINT fk_submissions_question
ADD FOREIGN KEY(question_id) REFERENCES questions(id)
ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE testcases
ADD CONSTRAINT fk_testcases
ADD FOREIGN KEY(question_id) REFERENCES questions(id)
ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE submissions
ADD CONSTRAINT fk_submissions_user
ADD FOREIGN KEY(user_id) REFERENCES users(id)
ON UPDATE NO ACTION ON DELETE CASCADE;

-- +goose Down
ALTER TABLE submissions DROP CONSTRAINT fk_submissions_question;
ALTER TABLE testcases DROP CONSTRAINT fk_testcases;
ALTER TABLE submissions DROP CONSTRAINT fk_submissions_user;