CREATE TABLE IF NOT EXISTS quiz (
    quiz_id TEXT NOT NULL,
    content TEXT NOT NULL,
    CONSTRAINT quiz_pk PRIMARY KEY (quiz_id)
);

CREATE TABLE IF NOT EXISTS question (
    question_id TEXT NOT NULL,
    content TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    answers JSONB NOT NULL,
    CONSTRAINT question_pk PRIMARY KEY (question_id)
);

CREATE TABLE IF NOT EXISTS quiz_question (
    quiz_id TEXT NOT NULL,
    question_id TEXT NOT NULL,
    CONSTRAINT quiz_question_pkey PRIMARY KEY (quiz_id, question_id),
    CONSTRAINT quiz_question_quiz_id_fkey FOREIGN KEY (quiz_id) REFERENCES quiz,
    CONSTRAINT quiz_question_question_id_fkey FOREIGN KEY (question_id) REFERENCES question
);
