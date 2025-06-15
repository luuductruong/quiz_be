package application

import (
	"github.com/quiz_be/services/gateway/application/http/ws"
	"google.golang.org/grpc"

	baseHttp "github.com/quiz_be/services/core/http"
	"github.com/quiz_be/services/gateway/application/http/quiz"
)

type HttpHandler interface {
	Quiz() quiz.HttpHandler
}

type httpHandler struct {
	quizHandler quiz.HttpHandler
}

func (h *httpHandler) Quiz() quiz.HttpHandler {
	return h.quizHandler
}

func NewHttpHandler(quizConn *grpc.ClientConn, hub *ws.WsHub) HttpHandler {
	baseHd := baseHttp.NewBaseHandler()
	return &httpHandler{
		quizHandler: quiz.NewHttpHandler(baseHd, quizConn, hub),
	}
}
