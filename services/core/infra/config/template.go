package config

import (
	"fmt"

	"github.com/quiz_be/services/core/infra/pubsub/rabbitmq"
)

type Client struct {
	Host string `config:"host"`
	Port string `config:"port"`
	SSL  bool   `config:"ssl"`
}

func (c *Client) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type PubSubConfig struct {
	Kind         string   `config:"kind"` // rabbitMQ, google, aws
	Topic        string   `config:"topic"`
	Subscription string   `config:"subscription"`
	Subscribes   []string `config:"subscribes"`
	// add more message queue here
	RbConfig *rabbitmq.RbConfig `config:"rabbitmq"`
}
