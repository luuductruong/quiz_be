package domain

import (
	"errors"
	"github.com/quiz_be/services/core/helper"
	"time"

	"github.com/quiz_be/services/core/context"
	model "github.com/quiz_be/services/core/domain/quiz"
)

// SubmitAnswer processes a user's answer submission for a quiz question and validates inputs, user, quiz, question, and answer.
// Returns the updated Quiz object or an error if the submission process fails.
func (d *domain) SubmitAnswer(ctx context.Context, inp *model.SubmitAnswerReq) (*model.Quiz, error) {
	d.logger.DebugCtx(ctx, "SubmitAnswer")
	// validate
	if inp.UserID == "" || inp.QuizID == "" || inp.QuestionID == "" || inp.AnswerTitle == "" {
		d.logger.DebugCtx(ctx, "invalid input")
		return nil, errors.New("invalid input")
	}
	// find exists user
	user, err := d.userRepo.Query(ctx).ByUserID(inp.UserID).Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query user failed")
		return nil, err
	}
	if user == nil {
		d.logger.DebugCtx(ctx, "user not found")
		return nil, errors.New("user not found")
	}
	// find existsQuiz
	existsQuiz, err := d.quizRepo.Query(ctx).
		ByQuizID(inp.QuizID).
		WithQuizQuestion(""). // embed a question to quiz_question, insert "" to get all QuizQuestion
		Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query fail")
		return nil, err
	}
	if existsQuiz == nil {
		d.logger.DebugCtx(ctx, "quiz not found")
		return nil, errors.New("quiz not found")
	}
	// find question
	var existsQuestion *model.Question
	for _, quizQuestion := range existsQuiz.QuizQuestions {
		if quizQuestion.Question != nil && quizQuestion.Question.QuestionID == inp.QuestionID {
			existsQuestion = quizQuestion.Question
			break
		}
	}
	if existsQuestion == nil {
		d.logger.DebugCtx(ctx, "not found question")
		return nil, errors.New("not found question")
	}
	// find answer
	var selectedAnswer *model.Answer
	for _, answer := range existsQuestion.Answers {
		if answer.Title == inp.AnswerTitle {
			selectedAnswer = answer
			break
		}
	}
	if selectedAnswer == nil {
		d.logger.DebugCtx(ctx, "not found answer")
		return nil, errors.New("not found answer")
	}
	err = d.processUserAnswer(ctx, inp.UserID, inp.QuizID, inp.AnswerTitle, existsQuestion)
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "processUserAnswer failed")
		return nil, err
	}
	return existsQuiz, nil
}

// processUserAnswer processes the user's selected answer for a quiz question and updates the score and user answer history.
// It determines correctness by checking the selected answer against the question's correct answer and adjusts scores.
func (d *domain) processUserAnswer(
	ctx context.Context,
	userID, quizID, selectedAnswer string,
	question *model.Question,
) error {
	d.logger.DebugCtx(ctx, "processUserAnswer")
	isCorrect := selectedAnswer == question.CorrectAnswer
	now := time.Now()
	// find user answer history
	userAnswer, err := d.userAnswerRepo.Query(ctx).
		ByUserID(userID).
		ByQuizID(quizID).
		ByQuestionID(question.QuestionID).
		Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query user answer failed")
		return err
	}
	prevCorrect := false
	scoreChange := 0
	if userAnswer == nil {
		userAnswer = &model.UserAnswer{
			ID:         helper.NewStringUUID(),
			UserID:     userID,
			QuizID:     quizID,
			QuestionID: question.QuestionID,
			CreatedAt:  now,
		}
		scoreChange = 0
		if isCorrect {
			scoreChange = int(question.Score)
		}
	} else {
		prevCorrect = userAnswer.IsCorrect
		if prevCorrect != isCorrect {
			if isCorrect {
				scoreChange = int(question.Score) // from incorrect to correct
			} else {
				scoreChange = -int(question.Score) // from correct to incorrect
			}
		}
	}

	// update user answer
	userAnswer.SelectedAnswer = selectedAnswer
	userAnswer.IsCorrect = isCorrect
	userAnswer.UpdatedAt = now
	if err := d.userAnswerRepo.Upsert(ctx, userAnswer); err != nil {
		d.logger.ErrorCtx(ctx, err, "update userAnswer failed")
		return err
	}
	// find score for user and quiz
	score, err := d.scoreRepo.Query(ctx).ByQuizID(quizID).ByUserID(userID).Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query score failed")
		return err
	}
	if score == nil {
		// first time
		score = &model.Score{
			ID:        helper.NewStringUUID(),
			UserID:    userID,
			QuizID:    quizID,
			Score:     max(scoreChange, 0), // không âm nếu lần đầu
			CreatedAt: now,
		}
	} else {
		// update score
		score.Score += scoreChange
		if score.Score < 0 {
			score.Score = 0
		}
	}
	score.UpdatedAt = now
	if err := d.scoreRepo.Upsert(ctx, score); err != nil {
		d.logger.ErrorCtx(ctx, err, "update score failed")
		return err
	}
	return nil
}

func (d *domain) GetLeaderboard(ctx context.Context, inp *model.GetLeaderboardReq) ([]*model.Score, int32, error) {
	d.logger.DebugCtx(ctx, "GetLeaderboard")
	if inp.QuizID == "" || inp.Page < 0 || inp.Limit < 0 {
		d.logger.DebugCtx(ctx, "invalid input")
		return nil, 0, errors.New("invalid input")
	}
	offset := helper.GetOffset(inp.Page, inp.Limit)
	scores, err := d.scoreRepo.Query(ctx).
		ByQuizID(inp.QuizID).
		Offset(offset).Limit(inp.Limit).
		OrderByScore(true).
		WithUser().
		ResultList()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query failed")
		return nil, 0, err
	}
	// get count
	total, err := d.scoreRepo.Query(ctx).ByQuizID(inp.QuizID).Count()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "count query failed")
		return nil, 0, err
	}
	return scores, int32(total), nil
}
