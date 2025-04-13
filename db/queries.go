package db

const createQuestionsTable = `
CREATE TABLE IF NOT EXISTS questions (
	id SERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	date DATE DEFAULT now() NOT NULL,
	likes INTEGER NOT NULL DEFAULT '0',
	dislikes INTEGER NOT NULL DEFAULT '0',
	created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);
CREATE INDEX IF NOT EXISTS ix_questions_date ON questions (date);
`
const createAnswersTable = `
CREATE TABLE IF NOT EXISTS answers (
	id SERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	question_id INTEGER,
	likes INTEGER NOT NULL DEFAULT '0',
	dislikes INTEGER NOT NULL DEFAULT '0',
	users_answered INTEGER NOT NULL DEFAULT '0',
	created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
	is_correct BOOLEAN DEFAULT 'f',
	FOREIGN KEY(question_id) REFERENCES questions (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS ix_answers_question_id ON answers (question_id);
`