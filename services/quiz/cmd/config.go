package cmd

import (
	"github.com/quiz_be/services/core/infra/config"
	"github.com/quiz_be/services/core/infra/db"
	"github.com/quiz_be/services/core/infra/logger"
)

type Config struct {
	Grpc     *config.Client `config:"grpc"`
	Logger   *logger.Config `config:"logger"`
	DataBase *db.Config     `config:"database"`
}
