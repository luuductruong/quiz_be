package domain

import (
	"errors"

	"github.com/quiz_be/services/core/context"
	model "github.com/quiz_be/services/core/domain/quiz"
)

func (d *domain) SubmitAnswer(ctx context.Context, inp *model.SubmitAnswerReq) (*model.Quiz, error) {
	d.logger.DebugCtx(ctx, "SubmitAnswer")
	// validate
	if inp.UserID == "" || inp.QuizID == "" || inp.QuestionID == "" || inp.AnswerTitle == "" {
		d.logger.DebugCtx(ctx, "invalid input")
		return nil, errors.New("invalid input")
	}
	// find existsQuiz
	existsQuiz, err := d.quizRepo.Query(ctx).
		ByQuizID(inp.QuizID).
		WithQuizQuestion(inp.QuestionID). // embed a question to quiz_question
		Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query fail")
		return nil, err
	}
	if existsQuiz == nil {
		d.logger.DebugCtx(ctx, "existsQuiz id not found")
		return nil, errors.New("existsQuiz id not found")
	}
	// find question
	var existsQuestion *model.Question
	for _, quizQuestion := range existsQuiz.QuizQuestions {
		if quizQuestion.QuestionID == inp.QuestionID && quizQuestion.Question != nil {
			existsQuestion = quizQuestion.Question
			break
		}
	}
	if existsQuestion == nil {
		d.logger.DebugCtx(ctx, "not found question")
		return nil, errors.New("not found question")
	}
	// find answer
	var existsAnswer *model.Answer
	for _, answer := range existsQuestion.Answers {
		if answer.Title == inp.AnswerTitle {
			existsAnswer = answer
			break
		}
	}
	if existsAnswer == nil {
		d.logger.DebugCtx(ctx, "not found answer")
		return nil, errors.New("not found answer")
	}
	// check correct response
	if existsAnswer.Title != existsQuestion.CorrectAnswer {
		d.logger.DebugCtx(ctx, "wrong answer")
		return nil, errors.New("wrong answer")
	}
	// increase score

	//

	return existsQuiz, nil
}
