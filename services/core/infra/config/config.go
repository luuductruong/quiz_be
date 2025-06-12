package config

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {
	viper.AddConfigPath(os.Getenv("CNF_DIR"))
	viper.AddConfigPath(".")
	viper.SetConfigName("app")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
}

func Unmarshal(i interface{}) error {
	return viper.GetViper().Unmarshal(i, func(config *mapstructure.DecoderConfig) {
		config.TagName = "config"
	})
}

func LoadConfig(i interface{}) error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return Unmarshal(i)
}
