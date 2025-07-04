package i18n

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kataras/i18n"
	"google.golang.org/protobuf/types/known/structpb"

	common "github.com/quiz_be/services/core/application/common/dto"
	"github.com/quiz_be/services/core/context"
	"github.com/quiz_be/services/core/middleware"
)

type LocKey string

func (key LocKey) String() string {
	return string(key)
}

type Config struct {
	Locale string `config:"locale"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Locale, validation.Required),
	)
}

const DefaultLocale = middleware.DefaultLocale

var defaultConfig = &Config{Locale: DefaultLocale}

var I18n *i18n.I18n

func Init(cfg *Config) {
	if cfg == nil {
		cfg = defaultConfig
	}
	SetDefault(cfg.Locale)
	I18n = i18n.Default
}

func SetDefault(langCode string) {
	if langCode == "" {
		langCode = DefaultLocale
	}
	i18n.SetDefaultLanguage(langCode)
}

func Lt(ctx context.Context, locKey LocKey, args ...interface{}) *common.LocalizedText {
	if I18n == nil {
		return &common.LocalizedText{
			Text: locKey.String(),
			Key:  locKey.String(),
		}
	}
	return lt(Tr(ctx, locKey, args...), locKey, args...)
}

func Tr(ctx context.Context, locKey LocKey, args ...interface{}) string {
	if I18n == nil {
		return ""
	}
	return I18n.Tr(ctx.GetLocale(), locKey.String(), args...)
}

func lt(text string, locKey LocKey, args ...interface{}) *common.LocalizedText {
	anyArgs := make([]*structpb.Value, len(args))

	for i, arg := range args {
		val, err := structpb.NewValue(arg)
		if err != nil {
			return nil
		}
		anyArgs[i] = val
	}

	return &common.LocalizedText{
		Text: text,
		Key:  locKey.String(),
		Args: anyArgs,
	}
}
