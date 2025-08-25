package pubsub

import (
	"context"
)

// Subscription represents a messaging endpoint where messages can be received.
// It provides a String method to get the subscription's string representation.
type Subscription interface {
	String() string
}

// SubscriptionHandler defines a function type that processes a Message within a given Context and returns an error if any.
type SubscriptionHandler func(context.Context, *Message) error

// SubscriptionRouter defines a function type that processes a Message using a given Context and a SubscriptionHandler, returning an error.
type SubscriptionRouter = func(msg *Message) SubscriptionHandler

// SubscriptionInterceptor defines a function type that intercepts a message processing flow in a subscription system.
// It receives a context, a message, and a handler, allowing for pre-processing, post-processing, or bypassing the handler.
// The function must return an error if the interception fails or processing cannot proceed.
type SubscriptionInterceptor func(ctx context.Context, msg *Message, handler SubscriptionHandler) error

// RoutingHandler returns a SubscriptionHandler that executes the handler derived from the given SubscriptionRouter for a message.
func RoutingHandler(route SubscriptionRouter) SubscriptionHandler {
	return func(ctx context.Context, msg *Message) error {
		fn := route(msg)
		return fn(ctx, msg)
	}
}

// Subscriber defines an interface for subscribing to and receiving messages from subscriptions.
// It provides methods for managing subscriptions and message processing.
type Subscriber interface {
	// Subscription returns a Subscription instance for the given ID.
	Subscription(id string) Subscription
	// Unsubscribe removes the given subscription from active subscriptions.
	Unsubscribe(Subscription)
	// Subscribe registers a handler for processing messages from the specified subscription.
	Subscribe(subscription Subscription, handler SubscriptionHandler)
	// Use adds middleware interceptors to the message processing pipeline.
	Use(middleware ...SubscriptionInterceptor)
	// StartReceiving begins receiving messages from all active subscriptions.
	StartReceiving() error
	// StopReceiving stops receiving messages from all subscriptions.
	StopReceiving()
}
