package pubsub

import (
	"context"

	ccontext "github.com/quiz_be/services/core/context"
)

var mqContextKey = "messages"

const MessageBatchMessageName MessageName = "MQ_MESSAGE_BATCH"

type ContextMessages struct {
	Messages []*Message
}

func batchHandler(route SubscriptionRouter) SubscriptionHandler {
	return func(ctx context.Context, msg *Message) error {
		var messages = new(ContextMessages)
		err := msg.ScanPayload(messages)
		if err != nil {
			return err
		}

		for _, message := range messages.Messages {
			fn := route(message)
			//TODO: check remove
			appCtx := ccontext.FromContext(ctx)
			err := fn(appCtx, message)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// MessageFromContext get the message batch from context. REQUIRES Middleware to have run.
func MessageFromContext(ctx context.Context) *ContextMessages {
	raw, _ := ctx.Value(mqContextKey).(*ContextMessages)
	return raw
}

func addMessagesToContext(ctx context.Context, messages ...*Message) {
	payload, _ := ctx.Value(mqContextKey).(*ContextMessages)
	if payload != nil {
		payload.Messages = append(payload.Messages, messages...)
	}
}
