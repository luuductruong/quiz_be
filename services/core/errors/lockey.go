package errors

import "github.com/quiz_be/services/core/i18n"

type LocKey = i18n.LocKey

const (
	LocKeyInternalServerError  LocKey = "internal_server_error"
	LocKeyInvalidArgumentError LocKey = "invalid_argument_error"
	LocKeyNotFoundError        LocKey = "not_found_error"
	LocKeyFailedPreCondition   LocKey = "failed_precondition"
)
