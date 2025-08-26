package subscriber

import (
	"context"
	appCtx "github.com/quiz_be/services/core/context"
	"github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/infra/logger"
	"github.com/quiz_be/services/core/infra/pubsub"
	"gorm.io/gorm"
)

type Handler interface {
	pubsub.BaseSubsHandler
}

type SubsHandlerParam struct {
	Service quiz.Service
	DB      *gorm.DB
	Logger  logger.Logger
}

func NewHandler(param *SubsHandlerParam) Handler {
	return &handler{
		domain: param.Service,
		db:     param.DB,
		logger: param.Logger,
	}
}

type handler struct {
	domain quiz.Service
	db     *gorm.DB
	logger logger.Logger
}

func (h *handler) RouteSetup() map[pubsub.MessageName]pubsub.SubscriptionHandler {
	return map[pubsub.MessageName]pubsub.SubscriptionHandler{
		quiz.PushJobGetQuizDetail: h.HandleGetQuizDetail,
	}
}

func (h *handler) HandleGetQuizDetail(ctx context.Context, msg *pubsub.Message) error {
	q := quiz.Quiz{}
	appContext := appCtx.FromContext(ctx).WithDbTx(h.db)
	err := msg.ScanPayload(&q)
	if err != nil {
		return err
	}
	h.logger.DebugCtx(appContext, "HandleGetQuizDetail: ", q.QuizID)
	_, err = h.domain.GetQuizDetail(appContext, q.QuizID)
	if err != nil {
		h.logger.ErrorCtx(appContext, err, "GetQuizDetail failed")
		return nil
	}
	return nil
}
