package subscriber

import (
	"context"
	"fmt"
	"github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/infra/pubsub"
)

type Handler interface {
	pubsub.BaseSubsHandler
}

func NewHandler(domain quiz.Service) Handler {
	return &handler{
		domain: domain,
	}
}

type handler struct {
	domain quiz.Service
}

func (h *handler) RouteSetup() map[pubsub.MessageName]pubsub.SubscriptionHandler {
	return map[pubsub.MessageName]pubsub.SubscriptionHandler{
		"test": h.Test,
	}
}

func (h *handler) Test(ctx context.Context, msg *pubsub.Message) error {
	q := quiz.Quiz{}
	err := msg.ScanPayload(&q)
	if err != nil {
		return err
	}
	fmt.Println("TESTTTTTTTTT: ", q.QuizID)
	return nil
}
