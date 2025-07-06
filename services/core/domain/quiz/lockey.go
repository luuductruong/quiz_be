package quiz

import "github.com/quiz_be/services/core/i18n"

type LocKey = i18n.LocKey

// Quiz error loc key
const (
	LocKeyQuizNotFound            = "quiz_not_found"
	LocKeyQuizQuestionHaveHistory = "quiz_question_have_history"
	LocKeyQuizTitleEmpty          = "quiz_title_empty"
)

// Question error loc key
const (
	LocKeyQuestionContentEmpty    = "question_content_empty"
	LocKeyInvalidScore            = "invalid_score"
	LocKeyQuestionNotFound        = "question_not_found"
	LocKeyInputEmptyListQuestion  = "input_empty_list_question"
	LocKeyInputDuplicatedQuestion = "input_duplicated_question"
)

// Answer error loc key
const (
	LocKeyEmptyListAnswer       = "empty_list_answer"
	LocKeyEmptyCorrectAnswer    = "empty_correct_answer"
	LocKeyInvalidInputAnswer    = "invalid_input_answer"
	LocKeyDuplicatedAnswerTitle = "duplicated_answer_title"
	LocKeyNotFoundCorrectAnswer = "not_found_correct_answer"
)

// User error loc key
const (
	LocKeyUserNotFound = "user_not_found"
)
