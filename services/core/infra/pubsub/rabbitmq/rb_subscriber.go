// Package rabbitmq implements RabbitMQ-specific publisher/subscriber functionality
// for message queue operations within the application.
package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/quiz_be/services/core/infra/logger"
	"github.com/quiz_be/services/core/infra/pubsub"
)

// newRbSubscriber creates a new RabbitMQ subscriber instance with the specified logger
// and RabbitMQ service. It returns an implementation of the pubsub.Subscriber interface.
func newRbSubscriber(logger logger.Logger, rbService *Service) pubsub.Subscriber {
	return &rbSubscriber{
		logger:    logger,
		svc:       rbService,
		observers: map[string]pubsub.SubscriptionHandler{},
	}
}

// rbSubscriber implements the pubsub.Subscriber interface for RabbitMQ.
// It manages subscriptions, message handling, and lifecycle of message consumption.
type rbSubscriber struct {
	logger       logger.Logger                         // Logger for operation logging
	svc          *Service                              // RabbitMQ service instance
	observers    map[string]pubsub.SubscriptionHandler // Map of subscription handlers
	interceptors []pubsub.SubscriptionInterceptor      // Slice of middleware interceptors
	stopSync     chan bool                             // Channel for stopping message consumption
}

// rbSubscriptionMsg represents a message received from RabbitMQ subscription
type rbSubscriptionMsg struct {
	data []byte // Raw message data
}

// Scan unmarshals the message data into the provided interface
func (rbMsg *rbSubscriptionMsg) Scan(i interface{}) error {
	return json.Unmarshal(rbMsg.data, &i)
}

// rbSubscription represents a RabbitMQ subscription identifier
type rbSubscription string

// String returns the string representation of the subscription
func (a rbSubscription) String() string {
	return string(a)
}

func (subscriber *rbSubscriber) Subscription(id string) pubsub.Subscription {
	return rbSubscription(id)
}

func (subscriber *rbSubscriber) Unsubscribe(subscription pubsub.Subscription) {
	return
}

// Subscribe registers a handler for the specified subscription
func (subscriber *rbSubscriber) Subscribe(subscription pubsub.Subscription, handler pubsub.SubscriptionHandler) {
	queueName := subscription.String()
	subscriber.logger.Debug("subscribe to queue: ", queueName)
	subscriber.observers[queueName] = handler
}

// Use adds middleware interceptors to the subscription pipeline
func (subscriber *rbSubscriber) Use(middleware ...pubsub.SubscriptionInterceptor) {
	if subscriber.interceptors == nil {
		subscriber.interceptors = []pubsub.SubscriptionInterceptor{}
	}
	subscriber.interceptors = append(subscriber.interceptors, middleware...)
}

// StartReceiving begins consuming messages from RabbitMQ for all registered subscriptions.
// It processes messages through the middleware chain and handlers.
func (subscriber *rbSubscriber) StartReceiving() error {
	// Initialize the stop synchronization channel if it hasn't been created yet
	if subscriber.stopSync == nil {
		subscriber.stopSync = make(chan bool, 1)
	}
	// Start a goroutine for each subscription handler to process messages asynchronously
	for _, handler := range subscriber.observers {
		go func() {
			subscriber.svc.ListenAndServe(func(msg []byte) error {
				// Parse the raw message bytes into a structured message format and validate
				var rbMsg = &rbSubscriptionMsg{data: msg}
				msgPayload := new(pubsub.Message)
				err := rbMsg.Scan(msgPayload)
				if err != nil {
					subscriber.logger.Error(fmt.Sprintf("can not parse subscription message err: %v, msgBody %+v", err, string(rbMsg.data)))
					return err
				}
				subscriber.logger.Debug("receive message: ", msgPayload)

				// Execute the message through the interceptor chain before passing to handler
				ctx := context.Background()
				chainInterceptor := pubsub.ChainSubscriptionInterceptor(subscriber.interceptors...)
				err = chainInterceptor(ctx, msgPayload, handler)
				if err != nil {
					subscriber.logger.Error(fmt.Sprintf("can not handle subscription message err: %v, msgBody %+v", err, string(rbMsg.data)))
					return err
				}
				return nil
			})
		}()
	}
	return nil
}

// StopReceiving halts message consumption for all subscriptions by closing
// the stop channel and cleaning up resources
func (subscriber *rbSubscriber) StopReceiving() {
	subscriber.logger.Debug("get stop to receive message")
	if subscriber.stopSync == nil {
		subscriber.logger.Debug("stop sync is nil")
		return
	}
	subscriber.logger.Debug("stop to receive messages")
	// close chan to stop all receiving subscription
	close(subscriber.stopSync)
	subscriber.stopSync = nil
}
