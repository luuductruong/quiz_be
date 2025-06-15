package dto

import (
	"github.com/quiz_be/services/core/application/quiz/model"
	"github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/helper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapQuizFromDomain(source *quiz.Quiz) *model.Quiz {
	if source == nil {
		return nil
	}
	return &model.Quiz{
		QuizId:    source.QuizID,
		Title:     source.Title,
		Questions: helper.MapList(questionsFromQuizQuestions(source.QuizQuestions), mapQuestionFromDomain),
	}
}

func questionsFromQuizQuestions(source []*quiz.QuizQuestion) []*quiz.Question {
	rs := make([]*quiz.Question, 0)
	for _, quizQuestion := range source {
		if quizQuestion.Question != nil {
			rs = append(rs, quizQuestion.Question)
		}
	}
	return rs
}

func mapQuestionFromDomain(source *quiz.Question) *model.Question {
	if source == nil {
		return nil
	}
	return &model.Question{
		QuestionId: source.QuestionID,
		Content:    source.Content,
		Score:      int32(source.Score),
		Answers:    helper.MapList(source.Answers, mapAnswerFromDomain),
	}
}

func mapAnswerFromDomain(source *quiz.Answer) *model.Answer {
	if source == nil {
		return nil
	}
	return &model.Answer{
		Title:   source.Title,
		Content: source.Content,
	}
}

func MapScoreFromDomain(source *quiz.Score) *model.Score {
	if source == nil {
		return nil
	}
	rs := &model.Score{
		UserId:   source.UserID,
		Score:    int32(source.Score),
		UpdateAt: timestamppb.New(source.UpdatedAt),
	}
	if source.User != nil {
		rs.UserName = source.User.UserName
	}
	return rs
}

func MapUserFromDomain(source *quiz.User) *model.User {
	if source == nil {
		return nil
	}
	return &model.User{
		UserId:   source.UserID,
		UserName: source.UserName,
	}
}
