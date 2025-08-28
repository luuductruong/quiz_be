package domain

import (
	"encoding/json"
	"github.com/quiz_be/services/core/infra/job"
	"time"

	"github.com/quiz_be/services/core/context"
	model "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/errors"
	"github.com/quiz_be/services/core/helper"
)

func (d *domain) ManageQuiz(ctx context.Context, inp *model.ManageQuizReq) (*model.Quiz, error) {
	d.logger.DebugCtx(ctx, "ManageQuiz: ", inp.QuizID, " localize: ", ctx.GetLocale())
	if inp.Title == "" {
		d.logger.DebugCtx(ctx, "invalid input title")
		return nil, errors.InvalidArgument(ctx, model.LocKeyQuizTitleEmpty)
	}
	inp.Title = helper.UpperFirstLetter(inp.Title)
	if len(inp.QuestionIDs) == 0 {
		d.logger.DebugCtx(ctx, "invalid input questions")
		return nil, errors.InvalidArgument(ctx, model.LocKeyInputEmptyListQuestion)
	}
	// find questions
	inpQuestion := len(inp.QuestionIDs)
	inp.QuestionIDs = helper.Unique(inp.QuestionIDs)
	if len(inp.QuestionIDs) != inpQuestion {
		d.logger.DebugCtx(ctx, "duplicated input questions")
		return nil, errors.InvalidArgument(ctx, model.LocKeyInputDuplicatedQuestion)
	}
	existsQuestions, err := d.questionRepo.Query(ctx).ByQuestionIDs(inp.QuestionIDs).ResultList()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "can't query questions by ids")
		return nil, errors.InternalDefault(ctx)
	}
	if len(existsQuestions) != inpQuestion {
		d.logger.DebugCtx(ctx, "not found questions")
		return nil, errors.NotFound(ctx, model.LocKeyQuestionNotFound)
	}
	topics := []string{}
	eventName := ""
	var handledQuiz *model.Quiz
	defer func() {
		if handledQuiz == nil || err != nil {
			return
		}
		d.pushJob(ctx, eventName, topics, &model.Quiz{
			QuizID:        handledQuiz.QuizID,
			Title:         handledQuiz.Title,
			CreatedAt:     handledQuiz.CreatedAt,
			UpdatedAt:     handledQuiz.UpdatedAt,
			QuizQuestions: nil, // set null to quiz questions
		})
		topics = nil
		eventName = ""
	}()

	now := time.Now()
	if inp.QuizID == "" {
		// create a new quiz
		topics = append(topics, job.TOPIC_SEARCH, job.TOPIC_LEADERBOARD)
		eventName = job.EVENT_QUIZ_CREATED
		totalQuiz, err := d.quizRepo.Query(ctx).Count()
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "count quiz failed")
			return nil, errors.InternalDefault(ctx)
		}
		newID := helper.GenQuizID(totalQuiz)

		handledQuiz = &model.Quiz{
			QuizID: newID,
			Title:  inp.Title,
			QuizQuestions: helper.MapList(existsQuestions, func(ques *model.Question) *model.QuizQuestion {
				return &model.QuizQuestion{
					QuizID:     newID,
					QuestionID: ques.QuestionID,
					Question:   ques,
				}
			}),
			CreatedAt: now,
			UpdatedAt: now,
		}
		//insert quiz
		err = d.quizRepo.Upsert(ctx, handledQuiz)
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "upsert quiz failed")
			return nil, errors.InternalDefault(ctx)
		}
		//insert quiz question
		err = d.quizQuestionRepo.BulkUpsert(ctx, handledQuiz.QuizQuestions)
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "bulk upsert quiz question failed")
			return nil, errors.InternalDefault(ctx)
		}
		return handledQuiz, nil
	}
	// updating
	topics = append(topics, job.TOPIC_LEADERBOARD)
	eventName = job.EVENT_QUIZ_UPDATED
	handledQuiz, err = d.quizRepo.Query(ctx).ByQuizID(inp.QuizID).WithQuizQuestion("").Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query quiz failed")
		return nil, errors.InternalDefault(ctx)
	}
	if handledQuiz == nil {
		d.logger.DebugCtx(ctx, "not found quiz")
		return nil, errors.NotFound(ctx, model.LocKeyQuizNotFound)
	}
	// find question that removed from quiz
	removedQuestions := make([]*model.QuizQuestion, 0)
	for _, qQ := range handledQuiz.QuizQuestions {
		found := false
		for _, questionId := range inp.QuestionIDs {
			if qQ.QuestionID == questionId {
				found = true
				break
			}
		}
		if !found {
			removedQuestions = append(removedQuestions, qQ)
		}
	}

	// check removed questions for answer history
	for _, qQ := range removedQuestions {
		existsAnswer, err := d.userAnswerRepo.Query(ctx).ByQuizID(inp.QuizID).ByQuestionID(qQ.QuestionID).Result()
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "query user answer failed")
			return nil, errors.InternalDefault(ctx)
		}
		if existsAnswer != nil {
			d.logger.DebugCtx(ctx, "question has answer history, can't remove")
			return nil, errors.FailedPreCondition(ctx, model.LocKeyQuizQuestionHaveHistory)
		}
	}

	// update quiz with new questions
	handledQuiz.Title = inp.Title
	handledQuiz.QuizQuestions = helper.MapList(existsQuestions, func(ques *model.Question) *model.QuizQuestion {
		return &model.QuizQuestion{
			QuizID:     handledQuiz.QuizID,
			QuestionID: ques.QuestionID,
			Question:   ques,
		}
	})
	handledQuiz.UpdatedAt = now

	// update quiz and quiz questions
	err = d.quizRepo.Upsert(ctx, handledQuiz)
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "upsert quiz failed")
		return nil, errors.InternalDefault(ctx)
	}

	err = d.quizQuestionRepo.BulkUpsert(ctx, handledQuiz.QuizQuestions)
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "bulk upsert quiz question failed")
		return nil, errors.InternalDefault(ctx)
	}
	return handledQuiz, nil
}

func (d *domain) ManageQuestion(ctx context.Context, inp *model.ManageQuestionReq) (*model.Question, error) {
	d.logger.DebugCtx(ctx, "ManageQuestion: ", inp.QuestionID)
	//validate input
	if inp.Content == "" {
		d.logger.DebugCtx(ctx, "invalid input content")
		return nil, errors.InvalidArgument(ctx, model.LocKeyQuestionContentEmpty)
	}
	if len(inp.Answers) == 0 {
		d.logger.DebugCtx(ctx, "invalid input answers")
		return nil, errors.InvalidArgument(ctx, model.LocKeyEmptyListAnswer)
	}
	if inp.CorrectAnswer == "" {
		d.logger.DebugCtx(ctx, "invalid input correct answer")
		return nil, errors.InvalidArgument(ctx, model.LocKeyEmptyCorrectAnswer)
	}
	if inp.Score <= 0 {
		d.logger.DebugCtx(ctx, "invalid input score")
		return nil, errors.InvalidArgument(ctx, model.LocKeyInvalidScore)
	}
	var inpCorrectAnswer *model.Answer
	mapAnswer := make(map[string]struct{})
	for _, answer := range inp.Answers {
		// check empty
		if answer.Content == "" || answer.Title == "" {
			d.logger.DebugCtx(ctx, "invalid input answers content or title")
			return nil, errors.InvalidArgument(ctx, model.LocKeyInvalidInputAnswer)
		}
		// check duplicate title
		if _, ok := mapAnswer[answer.Title]; ok {
			d.logger.DebugCtx(ctx, "duplicate answer title")
			return nil, errors.InvalidArgument(ctx, model.LocKeyDuplicatedAnswerTitle)
		}
		mapAnswer[answer.Title] = struct{}{}
		// check the correct answer title
		if answer.Title == inp.CorrectAnswer {
			inpCorrectAnswer = answer
		}
	}
	if inpCorrectAnswer == nil {
		d.logger.DebugCtx(ctx, "invalid input correct answer")
		return nil, errors.InvalidArgument(ctx, model.LocKeyNotFoundCorrectAnswer)
	}
	var err error
	topics := []string{}
	eventName := ""
	handledQuestion := &model.Question{}
	defer func() {
		if handledQuestion == nil || err != nil {
			return
		}
		// push job
		d.pushJob(ctx, eventName, topics, handledQuestion)
		// reset value
		topics = nil
		eventName = ""
	}()
	now := time.Now()
	if inp.QuestionID == "" {
		topics = append(topics, job.TOPIC_SEARCH)
		eventName = job.EVENT_QUESTION_CREATED
		// create new question
		totalQuestion, err := d.questionRepo.Query(ctx).Count()
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "count question failed")
			return nil, errors.InternalDefault(ctx)
		}
		newID := helper.GenQuestionID(totalQuestion)
		handledQuestion = &model.Question{
			QuestionID:    newID,
			Content:       inp.Content,
			Score:         inp.Score,
			Answers:       inp.Answers,
			CorrectAnswer: inp.CorrectAnswer,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		err = d.questionRepo.Upsert(ctx, handledQuestion)
		if err != nil {
			d.logger.ErrorCtx(ctx, err, "upsert question failed")
			return nil, errors.InternalDefault(ctx)
		}
		return handledQuestion, nil
	}
	// update exists question
	eventName = job.EVENT_QUESTION_UPDATED
	handledQuestion, err = d.questionRepo.Query(ctx).ByQuestionID(inp.QuestionID).Result()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "query question failed")
		return nil, errors.InternalDefault(ctx)
	}
	if handledQuestion == nil {
		d.logger.DebugCtx(ctx, "question not found")
		return nil, errors.NotFound(ctx, model.LocKeyQuestionNotFound)
	}
	// check content by title, score, correct answer
	isContentChanged := inp.Content != handledQuestion.Content ||
		inp.Score != handledQuestion.Score ||
		inp.CorrectAnswer != handledQuestion.CorrectAnswer
	// check content of correct answer
	if !isContentChanged {
		existsCorrectAnswer, _ := helper.Find(handledQuestion.Answers, func(answer *model.Answer) bool {
			return answer.Title == handledQuestion.CorrectAnswer
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
			return nil, errors.InternalDefault(ctx)
		}
		if existsAnswer != nil {
			d.logger.DebugCtx(ctx, "question has leaderboard history, can't update")
			return nil, errors.FailedPreCondition(ctx, model.LocKeyQuizQuestionHaveHistory)
		}
	}
	// update
	handledQuestion.Content = inp.Content
	handledQuestion.Score = inp.Score
	handledQuestion.CorrectAnswer = inp.CorrectAnswer
	handledQuestion.Answers = inp.Answers
	handledQuestion.UpdatedAt = now
	err = d.questionRepo.Upsert(ctx, handledQuestion)
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "upsert question failed")
		return nil, err
	}
	return handledQuestion, nil
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
		return errors.InternalDefault(ctx)
	}
	err = d.scoreRepo.Query(ctx).ByQuizID(quizId).Delete()
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "delete score failed")
		return errors.InternalDefault(ctx)
	}
	return nil
}

func (d *domain) pushJob(ctx context.Context, name string, topics []string, data any) {
	d.logger.DebugCtx(ctx, "pushJob")
	if name == "" {
		d.logger.DebugCtx(ctx, "name is empty")
		return
	}
	if len(topics) == 0 {
		d.logger.DebugCtx(ctx, "topics is empty")
		return
	}
	payload, err := json.Marshal(data)
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "marshal payload failed")
		return
	}
	err = d.jobClient.PushJob(ctx, name, topics, payload)
	if err != nil {
		d.logger.ErrorCtx(ctx, err, "push job failed")
	}
	return
}
