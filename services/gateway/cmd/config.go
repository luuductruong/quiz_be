package cmd

import (
	cfn "github.com/quiz_be/services/core/infra/config"
	"github.com/quiz_be/services/core/infra/logger"
)

type Config struct {
	Http   *cfn.Client    `config:"http"`
	Logger *logger.Config `config:"logger"`
	// microservices
	Quiz *cfn.Client `config:"quiz"`
}
