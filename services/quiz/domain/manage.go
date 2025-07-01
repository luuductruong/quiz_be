package domain

import (
	"errors"
	"fmt"
	"github.com/quiz_be/services/core/context"
	model "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/helper"
	"time"
)

func (d *domain) ManageQuiz(ctx context.Context, inp *model.ManageQuizReq) (*model.Quiz, error) {
	return &model.Quiz{}, nil
}

func (d *domain) ManageQuestion(ctx context.Context, inp *model.ManageQuestionReq) (*model.Question, error) {
	d.logger.DebugCtx(ctx, "ManageQuestion")
	//validate input
	if inp.Content == "" {
		d.logger.DebugCtx(ctx, "invalid input content")
		return nil, errors.New("invalid input content")
	}
	if len(inp.Answers) == 0 {
		d.logger.DebugCtx(ctx, "invalid input answers")
		return nil, errors.New("invalid input answers")
	}
	if inp.CorrectAnswer == "" {
		d.logger.DebugCtx(ctx, "invalid input correct answer")
		return nil, errors.New("invalid input correct answer")
	}
	if inp.Score <= 0 {
		d.logger.DebugCtx(ctx, "invalid input score")
		return nil, errors.New("invalid input score")
	}
	var inpCorrectAnswer *model.Answer
	validCorrectAnswer := false
	mapAnswer := make(map[string]struct{})
	for _, answer := range inp.Answers {
		// check empty
		if answer.Content == "" || answer.Title == "" {
			d.logger.DebugCtx(ctx, "invalid input answers content or title")
			return nil, errors.New("invalid input answers content or title")
		}
		// check duplicate title
		if _, ok := mapAnswer[answer.Title]; ok {
			d.logger.DebugCtx(ctx, "duplicate answer title")
			return nil, errors.New("duplicate answer title")
		}
		mapAnswer[answer.Title] = struct{}{}
		// check the correct answer title
		if answer.Title == inp.CorrectAnswer {
			validCorrectAnswer = true
			inpCorrectAnswer = answer
		}
	}
	if !validCorrectAnswer {
		d.logger.DebugCtx(ctx, "invalid input correct answer")
		return nil, errors.New("invalid input correct answer")
	}
	now := time.Now()
	if inp.QuestionID == "" {
		// create new question
		totalQuestion, err := d.questionRepo.Query(ctx).Count()
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "count question failed")
			return nil, err
		}
		newID := fmt.Sprintf("question_id_%03d", totalQuestion+1)
		newQuestion := &model.Question{
			QuestionID:    newID,
			Content:       inp.Content,
			Score:         inp.Score,
			Answers:       inp.Answers,
			CorrectAnswer: inp.CorrectAnswer,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		err = d.questionRepo.Upsert(ctx, newQuestion)
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "upsert question failed")
			return nil, err
		}
		return newQuestion, nil
	}
	// update exists question
	existsQuestion, err := d.questionRepo.Query(ctx).ByQuestionID(inp.QuestionID).Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query question failed")
		return nil, err
	}
	if existsQuestion == nil {
		d.logger.DebugCtx(ctx, "question not found")
		return nil, errors.New("question not found")
	}
	// check content by title, score, correct answer
	isContentChanged := inp.Content != existsQuestion.Content ||
		inp.Score != existsQuestion.Score ||
		inp.CorrectAnswer != existsQuestion.CorrectAnswer
	// check content of correct answer
	if !isContentChanged {
		existsCorrectAnswer, _ := helper.Find(existsQuestion.Answers, func(answer *model.Answer) bool {
			return answer.Title == existsQuestion.CorrectAnswer
		})
		if existsCorrectAnswer == nil {
			isContentChanged = true
		} else {
			isContentChanged = existsCorrectAnswer.Content != inpCorrectAnswer.Content
		}
	}
	if isContentChanged {
		// check leaderboard history
		existsAnswer, err := d.userAnswerRepo.Query(ctx).ByQuestionID(inp.QuestionID).Result()
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "query user answer failed")
			return nil, err
		}
		if existsAnswer != nil {
			d.logger.DebugCtx(ctx, "question has leaderboard history, can't update")
			return nil, errors.New("question has leaderboard history, can't update")
		}
	}
	// update
	existsQuestion.Content = inp.Content
	existsQuestion.Score = inp.Score
	existsQuestion.CorrectAnswer = inp.CorrectAnswer
	existsQuestion.Answers = inp.Answers
	existsQuestion.UpdatedAt = now
	err = d.questionRepo.Upsert(ctx, existsQuestion)
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "upsert question failed")
		return nil, err
	}
	return existsQuestion, nil
}

func (d *domain) clearLeaderboard(ctx context.Context, quizId string) error {
	d.logger.DebugCtx(ctx, "clearLeaderboard")
	if quizId == "" {
		d.logger.DebugCtx(ctx, "invalid input, return nil")
		return nil
	}
	err := d.userAnswerRepo.Query(ctx).ByQuizID(quizId).Delete()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "delete user answer failed")
		return err
	}
	err = d.scoreRepo.Query(ctx).ByQuizID(quizId).Delete()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "delete score failed")
		return err
	}
	return nil
}
