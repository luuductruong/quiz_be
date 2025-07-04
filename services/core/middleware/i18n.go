package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const localeContextKeyName = "i18n-locale"
const acceptLanguageContextKeyName = "accept-language"

var i18nLocaleContextKey = localeContextKeyName

func GrpcLocaleMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		locale := ""
		acceptLanguage := md.Get(localeContextKeyName)
		if len(acceptLanguage) > 0 {
			locale = acceptLanguage[0]
		}

		if locale == "" {
			acceptLanguage = md.Get(acceptLanguageContextKeyName)
			if len(acceptLanguage) > 0 {
				locale = acceptLanguage[0]
			}
		}

		ctx = context.WithValue(ctx, i18nLocaleContextKey, locale)
		ctx = metadata.AppendToOutgoingContext(ctx, localeContextKeyName, locale)
		return handler(ctx, req)
	}
}

func LocaleFromContext(ctx context.Context) string {
	raw, _ := ctx.Value(i18nLocaleContextKey).(string)
	return raw
}

func LocaleToContext(ctx context.Context, locale string) context.Context {
	if currentLocale, ok := ctx.Value(i18nLocaleContextKey).(string); !ok || currentLocale != locale {
		ctx = context.WithValue(ctx, i18nLocaleContextKey, locale)
	}
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md[localeContextKeyName] = []string{locale}
		return metadata.NewOutgoingContext(ctx, md)
	} else {
		return metadata.AppendToOutgoingContext(ctx, localeContextKeyName, locale)
	}
}
