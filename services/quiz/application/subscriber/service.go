package subscriber

import (
	"context"
	appContext "github.com/quiz_be/services/core/context"
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
	return map[pubsub.MessageName]pubsub.SubscriptionHandler{}
}

func (h *handler) Test(ctx context.Context, msg *pubsub.Message) error {
	appCtx := appContext.FromContext(ctx)
	type Test struct {
		ProductID int64
	}
	var test Test
	err := msg.ScanPayload(&test)
	if err != nil {
		return err
	}
	if _, err = h.domain.GetQuizDetail(appCtx, "10"); err != nil {
		return err
	}
	return nil
}
