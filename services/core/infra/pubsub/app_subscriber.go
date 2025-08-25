package pubsub

import (
	"context"
	"github.com/quiz_be/services/core/infra/db"
	"github.com/quiz_be/services/core/infra/logger"
)

type BaseSubsHandler interface {
	RouteSetup() map[MessageName]SubscriptionHandler
}

type AppSubscriber interface {
	Subscriber
	RegisterEventSubscriber(subs map[MessageName]SubscriptionHandler)
}

func NewAppSubscriber(sb Subscriber, subscription string, db db.SQL, messageBus ...MessageBus) AppSubscriber {
	subs := sb.Subscription(subscription)
	appSubs := &appSubscriber{
		Subscriber:   sb,
		db:           db,
		subscription: subs,
	}

	if len(messageBus) > 0 {
		appSubs.messageBus = messageBus[0]
	}
	return appSubs
}

type appSubscriber struct {
	Subscriber
	db           db.SQL
	subscription Subscription
	messageBus   MessageBus
}

func (a *appSubscriber) RegisterEventSubscriber(subs map[MessageName]SubscriptionHandler) {
	var route SubscriptionRouter
	route = func(msg *Message) SubscriptionHandler {
		handler := subs[msg.Name]
		if handler == nil {
			//ignore unknown message
			handler = func(ctx context.Context, message *Message) error {
				logger.Default.Debug("receive message, but not handle: ", message.Name)
				return nil
			}
		}
		return handler
	}

	//set handler for message batch
	if subs[MessageBatchMessageName] == nil {
		subs[MessageBatchMessageName] = batchHandler(route)
	}
	// run with order top to bottom
	a.Use(
	// add more
	)
	a.Subscribe(a.subscription, RoutingHandler(route))
}
