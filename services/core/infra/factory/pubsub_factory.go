package factory

import (
	"errors"
	"github.com/quiz_be/services/core/infra/config"
	"github.com/quiz_be/services/core/infra/logger"
	"github.com/quiz_be/services/core/infra/pubsub"
	"github.com/quiz_be/services/core/infra/pubsub/rabbitmq"
)

func NewPubSub(logger logger.Logger, psConfig *config.PubSubConfig) (pubsub.PubSub, error) {
	if psConfig == nil {
		return nil, errors.New("PubSubConfig is nil")
	}
	switch psConfig.Kind {
	case pubsub.RabbitPubSubKind:
		psConfig.RbConfig.Topic = psConfig.Topic
		psConfig.RbConfig.Subscribes = psConfig.Subscribes
		return rabbitmq.NewRbPubSub(logger, psConfig.RbConfig), nil
	case pubsub.GooglePubSubKind:
		return nil, errors.New("not implemented")
	case pubsub.AWSPubSubKind:
		return nil, errors.New("not implemented")
	default:
		return nil, errors.New("PubSubConfig.Kind is invalid")
	}
}
