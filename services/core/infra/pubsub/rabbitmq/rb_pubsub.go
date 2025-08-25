package rabbitmq

import (
	"github.com/quiz_be/services/core/infra/logger"
	"github.com/quiz_be/services/core/infra/pubsub"
)

type rbPs struct {
	publisher  pubsub.Publisher
	subscriber pubsub.Subscriber
	service    *Service
}

// Publisher returns publisher object
func (p *rbPs) Publisher() pubsub.Publisher {
	return p.publisher
}

// Subscriber subscribe, and syncing on a subscription
func (p *rbPs) Subscriber() pubsub.Subscriber {
	return p.subscriber
}

// Close pubsub client
func (p *rbPs) Close() {
	p.service.Close()
}

func NewRbPubSub(log logger.Logger, config *RbConfig) pubsub.PubSub {
	svc := NewService(config)
	err := svc.Connect()
	if err != nil {
		panic(err)
	}
	return &rbPs{
		service:    svc,
		publisher:  newRbPublisher(log, svc),
		subscriber: newRbSubscriber(log, svc),
	}
}
