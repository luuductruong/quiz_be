package pubsub

import "context"

func ChainSubscriptionInterceptor(interceptors ...SubscriptionInterceptor) SubscriptionInterceptor {
	n := len(interceptors)

	return func(ctx context.Context, info *Message, handler SubscriptionHandler) error {
		chainer := func(currentInter SubscriptionInterceptor, currentHandler SubscriptionHandler) SubscriptionHandler {
			return func(currentCtx context.Context, currentInfo *Message) error {
				return currentInter(currentCtx, currentInfo, currentHandler)
			}
		}

		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(interceptors[i], chainedHandler)
		}

		return chainedHandler(ctx, info)
	}
}
