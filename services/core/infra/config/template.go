package config

import "fmt"

type Client struct {
	Host string `config:"host"`
	Port string `config:"port"`
	SSL  bool   `config:"ssl"`
}

func (c *Client) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
