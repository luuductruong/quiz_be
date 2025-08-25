package rabbitmq

import (
	"encoding/json"
	"github.com/quiz_be/services/core/infra/logger"
	"github.com/quiz_be/services/core/infra/pubsub"
)

// rb rbPublisher
func newRbPublisher(logger logger.Logger, rbService *Service) pubsub.Publisher {
	return &rbPublisher{
		logger: logger,
		svc:    rbService,
	}
}

type rbPublisher struct {
	logger logger.Logger
	svc    *Service
}

func (publisher *rbPublisher) Publish(topic pubsub.Topic, msg *pubsub.Message) error {
	publisher.logger.Debug("publish message with topic: ", topic.String(), " and message: ", msg)
	messageByte, err := json.Marshal(msg)
	if err != nil {
		publisher.logger.Error("Error marshaling message: ", err)
		return err
	}
	err = publisher.svc.Publish(messageByte, topic.String())
	if err != nil {
		publisher.logger.Error("Error publishing message: ", err)
		return err
	}
	publisher.logger.Debug("message published successfully")
	return nil
}

func (publisher *rbPublisher) PublishRaw(topic pubsub.Topic, rawMsg []byte) error {
	publisher.logger.Debug("publish raw message with topic: ", topic.String(), " and message: ", string(rawMsg))
	err := publisher.svc.Publish(rawMsg, topic.String())
	if err != nil {
		publisher.logger.Error("Error publishing message: ", err)
		return err
	}
	publisher.logger.Debug("message published successfully")
	return nil
}

func (publisher *rbPublisher) Topic(id string) pubsub.Topic {
	return rbTopic(id)
}

// just a type which will implement interface Topic from message.go
type rbTopic string

func (a rbTopic) String() string {
	return string(a)
}
