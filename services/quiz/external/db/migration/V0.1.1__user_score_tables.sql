CREATE TABLE IF NOT EXISTS users (
    user_id TEXT NOT NULL,
    user_name TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    CONSTRAINT user_pk PRIMARY KEY (user_id),
    CONSTRAINT user_uk UNIQUE (phone_number)
);

CREATE TABLE IF NOT EXISTS scores (
    id TEXT NOT NULL,
    CONSTRAINT score_pk PRIMARY KEY (id),
    user_id TEXT NOT NULL,
    quiz_id TEXT NOT NULL,
    score INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT score_uk UNIQUE (user_id, quiz_id),
    CONSTRAINT score_user_fkey FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT score_quiz_fkey FOREIGN KEY (quiz_id) REFERENCES quiz(quiz_id)
);

CREATE INDEX IF NOT EXISTS idx_scores_user_id
    ON scores (user_id);

CREATE TABLE IF NOT EXISTS user_answer (
    id TEXT NOT NULL,
    CONSTRAINT user_answer_pk PRIMARY KEY (id),
    user_id TEXT NOT NULL,
    quiz_id TEXT NOT NULL,
    question_id TEXT NOT NULL,
    selected_answer TEXT NOT NULL,
    is_correct BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_answer_uk UNIQUE (user_id, quiz_id, question_id),
    CONSTRAINT user_answer_user_fkey FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT user_answer_quiz_fkey FOREIGN KEY (quiz_id) REFERENCES quiz(quiz_id),
    CONSTRAINT user_answer_question_fkey FOREIGN KEY (question_id) REFERENCES question(question_id)
);

ALTER TABLE question
    ADD COLUMN IF NOT EXISTS score INTEGER;

ALTER TABLE quiz
    RENAME COLUMN content TO title;