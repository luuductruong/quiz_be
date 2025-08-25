// Package pubsub provides interfaces and types for implementing publish-subscribe messaging patterns.
// It supports various message brokers and includes functionality for message routing and interception.
package pubsub

import (
	"encoding/json"
)

const (
	RabbitPubSubKind = "rabbitMQ"
	GooglePubSubKind = "googleCloud"
	AWSPubSubKind    = "amazonAWS"
)

// Message represents a data structure containing a type identifier (Kind) and associated binary content (Payload).
type Message struct {
	// ID represents the unique identifier of the message, used to distinguish it from other messages.
	ID string `json:"id"`

	Kind    string      `json:"kind"`
	Name    MessageName `json:"name"`
	Payload []byte      `json:"payload"`
}

func NewMessage(kind string, payload []byte) *Message {
	return &Message{
		Kind:    kind,
		Payload: payload,
	}
}

func (m Message) ScanPayload(payload interface{}) error {
	return json.Unmarshal(m.Payload, payload)
}

// PubSub defines an interface for a publish-subscribe mechanism with methods to access publisher and subscriber instances.
type PubSub interface {
	// Publisher method returns an associated Publisher instance for publishing messages.
	Publisher() Publisher
	// Subscriber method returns a Subscriber instance for subscribing to messages.
	Subscriber() Subscriber
	// Close method is used to close and clean up resources used by the PubSub implementation.
	Close()
}
