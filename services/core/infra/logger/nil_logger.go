package logger

import "github.com/quiz_be/services/core/context"

type nilLogger struct {
}

func init() {
	Default = NewNilLogger()
}

func NewNilLogger() Logger {
	return &nilLogger{}
}

func (n *nilLogger) Fatal(i ...interface{}) {}

func (n *nilLogger) Panic(i ...interface{}) {}

func (n *nilLogger) Error(i ...interface{}) {}

func (n *nilLogger) Debug(i ...interface{}) {}

func (n *nilLogger) ErrorCtx(context context.Context, err error, i ...interface{}) {}

func (n *nilLogger) DebugCtx(context context.Context, i ...interface{}) {}
