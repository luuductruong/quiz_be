package helper

import "fmt"

const (
	QuizIDCode     = "quiz_id_"
	QuestionIDCode = "q"
)

func GenQuizID(countQuizzes int) string {
	return fmt.Sprintf("%s%02d", QuizIDCode, countQuizzes+1)
}

func GenQuestionID(countQuestions int) string {
	return fmt.Sprintf("%s%d", QuestionIDCode, countQuestions+1)
}
