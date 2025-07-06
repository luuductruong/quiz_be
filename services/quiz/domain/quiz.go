package domain

import (
	"time"

	"github.com/quiz_be/services/core/context"
	model "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/errors"
	"github.com/quiz_be/services/core/helper"
)

// SubmitAnswer processes a user's answer submission for a quiz question and validates inputs, user, quiz, question, and answer.
// Returns the updated Quiz object or an error if the submission process fails.
func (d *domain) SubmitAnswer(ctx context.Context, inp *model.SubmitAnswerReq) (*model.Quiz, error) {
	d.logger.DebugCtx(ctx, "SubmitAnswer")
	// validate
	if inp.UserID == "" || inp.QuizID == "" || len(inp.SelectAnswers) == 0 {
		d.logger.DebugCtx(ctx, "invalid input")
		return nil, errors.InvalidArgument(ctx, errors.LocKeyInvalidArgumentError)
	}
	// find exists user
	user, err := d.userRepo.Query(ctx).ByUserID(inp.UserID).Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query user failed")
		return nil, errors.InternalDefault(ctx)
	}
	if user == nil {
		d.logger.DebugCtx(ctx, "user not found")
		return nil, errors.NotFound(ctx, model.LocKeyUserNotFound)
	}
	// find existsQuiz
	existsQuiz, err := d.quizRepo.Query(ctx).
		ByQuizID(inp.QuizID).
		WithQuizQuestion(""). // embed a question to quiz_question, insert "" to get all QuizQuestion
		Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query fail")
		return nil, errors.InternalDefault(ctx)
	}
	if existsQuiz == nil {
		d.logger.DebugCtx(ctx, "quiz not found")
		return nil, errors.NotFound(ctx, model.LocKeyQuizNotFound)
	}
	err = d.processUserAnswer(ctx, inp.UserID, existsQuiz, inp.SelectAnswers)
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
	userID string,
	existsQuiz *model.Quiz,
	selectedAnswer []*model.SelectedQuestionAnswer,
) error {
	d.logger.DebugCtx(ctx, "processUserAnswer")
	now := time.Now()
	// find user answer history
	userAnswers, err := d.userAnswerRepo.Query(ctx).
		ByUserID(userID).
		ByQuizID(existsQuiz.QuizID).
		ResultList()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query user answer failed")
		return errors.InternalDefault(ctx)
	}
	questions := helper.MapList(existsQuiz.QuizQuestions, func(q *model.QuizQuestion) *model.Question {
		return q.Question
	})
	totalScoreChange := 0
	for _, ans := range selectedAnswer {
		question, _ := helper.Find(questions, func(q *model.Question) bool {
			return q.QuestionID == ans.QuestionID
		})
		if question == nil {
			d.logger.DebugCtx(ctx, "question not found")
			return errors.NotFound(ctx, model.LocKeyQuestionNotFound)
		}
		isCorrect := ans.AnswerTitle == question.CorrectAnswer
		prevCorrect := false
		userAnswer, _ := helper.Find(userAnswers, func(ua *model.UserAnswer) bool {
			return ua.QuestionID == question.QuestionID
		})
		scoreChange := 0
		if userAnswer == nil {
			userAnswer = &model.UserAnswer{
				ID:         helper.NewStringUUID(),
				UserID:     userID,
				QuizID:     existsQuiz.QuizID,
				QuestionID: ans.QuestionID,
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
		userAnswer.SelectedAnswer = ans.AnswerTitle
		userAnswer.IsCorrect = isCorrect
		userAnswer.UpdatedAt = now
		if err := d.userAnswerRepo.Upsert(ctx, userAnswer); err != nil {
			d.logger.ErrorCtx(ctx, err, "update userAnswer failed")
			return errors.InternalDefault(ctx)
		}
		totalScoreChange += scoreChange
	}

	// find score for user and quiz
	score, err := d.scoreRepo.Query(ctx).ByQuizID(existsQuiz.QuizID).ByUserID(userID).Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query score failed")
		return errors.InternalDefault(ctx)
	}
	if score == nil {
		// first time
		score = &model.Score{
			ID:        helper.NewStringUUID(),
			UserID:    userID,
			QuizID:    existsQuiz.QuizID,
			Score:     max(totalScoreChange, 0), // không âm nếu lần đầu
			CreatedAt: now,
		}
	} else {
		// update score
		score.Score += totalScoreChange
		if score.Score < 0 {
			score.Score = 0
		}
	}
	score.UpdatedAt = now
	if err := d.scoreRepo.Upsert(ctx, score); err != nil {
		d.logger.ErrorCtx(ctx, err, "update score failed")
		return errors.InternalDefault(ctx)
	}
	return nil
}

func (d *domain) GetLeaderboard(ctx context.Context, inp *model.GetLeaderboardReq) ([]*model.Score, int32, error) {
	d.logger.DebugCtx(ctx, "GetLeaderboard")
	if inp.QuizID == "" || inp.Page < 0 || inp.Limit < 0 {
		d.logger.DebugCtx(ctx, "invalid input")
		return nil, 0, errors.InvalidArgument(ctx, errors.LocKeyInvalidArgumentError)
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
		return nil, 0, errors.InternalDefault(ctx)
	}
	// get count
	total, err := d.scoreRepo.Query(ctx).ByQuizID(inp.QuizID).Count()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "count query failed")
		total = len(scores)
	}
	return scores, int32(total), nil
}

func (d *domain) GetListQuiz(ctx context.Context, page int, limit int) ([]*model.Quiz, int32, error) {
	d.logger.DebugCtx(ctx, "GetListQuiz")
	offset := helper.GetOffset(page, limit)
	listQuiz, err := d.quizRepo.Query(ctx).
		Offset(offset).
		Limit(limit).
		ResultList()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query failed")
		return nil, 0, errors.InternalDefault(ctx)
	}
	count, err := d.quizRepo.Query(ctx).Count()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "count query failed")
	}
	return listQuiz, int32(count), nil
}

func (d *domain) GetQuizDetail(ctx context.Context, quizID string) (*model.Quiz, error) {
	d.logger.DebugCtx(ctx, "GetQuizDetail")
	if quizID == "" {
		d.logger.DebugCtx(ctx, "invalid input")
		return nil, errors.InvalidArgument(ctx, errors.LocKeyInvalidArgumentError)
	}
	quiz, err := d.quizRepo.Query(ctx).ByQuizID(quizID).WithQuizQuestion("").Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query failed")
		return nil, errors.InternalDefault(ctx)
	}
	if quiz == nil {
		d.logger.DebugCtx(ctx, "quiz not found")
		return nil, errors.NotFound(ctx, model.LocKeyQuizNotFound)
	}
	return quiz, nil
}

func (d *domain) GetUser(ctx context.Context) ([]*model.User, error) {
	d.logger.DebugCtx(ctx, "GetUser")
	users, err := d.userRepo.Query(ctx).ResultList()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query failed")
		return nil, errors.InternalDefault(ctx)
	}
	return users, nil
}
