package pubsub

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/quiz_be/services/core/context"
)

type MessageName string

func (n MessageName) String() string {
	return string(n)
}

type MessageQueueJobPayload struct {
	Topic   string   `json:"topic"`
	Message *Message `json:"message"`
}

type MessageBus interface {
	Publish(ctx context.Context, name MessageName, payload interface{}) error
	Submit(ctx context.Context) error
}

func newRawMessageJobPayload(ctx context.Context, topic string, name MessageName, rawPayload []byte) (*MessageQueueJobPayload, error) {
	messageId := uuid.NewString()
	return &MessageQueueJobPayload{
		Topic: topic,
		Message: &Message{
			ID:      messageId,
			Kind:    "",
			Payload: rawPayload,
		},
	}, nil
}

func newMessageJobPayload(ctx context.Context, topic string, name MessageName, payload interface{}) (*MessageQueueJobPayload, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return newRawMessageJobPayload(ctx, topic, name, jsonPayload)
}

type batchStagedMessageBus struct {
	topic string
}

func NewBatchStagedMessageBus(topic string) MessageBus {
	return &batchStagedMessageBus{
		topic: topic,
	}
}

func (bus *batchStagedMessageBus) Publish(ctx context.Context, name MessageName, payload interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (bus *batchStagedMessageBus) Submit(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
