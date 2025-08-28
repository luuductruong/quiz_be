package job

// Event names - Convention: UPPERCASE with '_' and domain prefixes (e.g. EVENT_[DOMAIN]_[ACTION])
// Values use lowercase with dots to maintain message broker compatibility
const (
	EVENT_QUIZ_CREATED = "quiz.quiz_created"
	EVENT_QUIZ_UPDATED = "quiz.quiz_updated"
	EVENT_QUIZ_DELETED = "quiz.quiz_deleted"

	EVENT_QUESTION_CREATED = "quiz.question_created"
	EVENT_QUESTION_UPDATED = "quiz.question_updated"
	EVENT_QUESTION_DELETED = "quiz.question_deleted"
)

// Topic names - Convention: UPPERCASE with '_' (e.g. TOPIC_[NAME])
// Values should match the topic/exchange names in configuration (lowercase with dots)
const (
	TOPIC_LEADERBOARD = "leaderboard.events"
	TOPIC_SEARCH      = "search.indexing"
)
