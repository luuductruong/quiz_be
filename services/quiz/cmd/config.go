package cmd

import (
	"github.com/quiz_be/services/core/i18n"
	"github.com/quiz_be/services/core/infra/config"
	"github.com/quiz_be/services/core/infra/db"
	"github.com/quiz_be/services/core/infra/logger"
)

type Config struct {
	Grpc     *config.Client       `config:"grpc"`
	Logger   *logger.Config       `config:"logger"`
	I18n     *i18n.Config         `config:"i18n"`
	DataBase *db.Config           `config:"database"`
	PubSub   *config.PubSubConfig `config:"pubsub"`

	// client
	Job *config.Client `config:"job"`
}
