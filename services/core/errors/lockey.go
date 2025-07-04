package errors

import "github.com/quiz_be/services/core/i18n"

type LocKey = i18n.LocKey

const (
	LocKeyInternalServerError  LocKey = "InternalServerError"
	LocKeyInvalidArgumentError LocKey = "InvalidArgumentError"
	LocKeyNotFoundError        LocKey = "NotFoundError"
)
