package cmd

import (
	"github.com/quiz_be/services/core/i18n"
	cfn "github.com/quiz_be/services/core/infra/config"
	"github.com/quiz_be/services/core/infra/logger"
)

type Config struct {
	Http   *cfn.Client    `config:"http"`
	Logger *logger.Config `config:"logger"`
	I18n   *i18n.Config   `config:"i18n"`
	// microservices
	Quiz *cfn.Client `config:"quiz"`
}
