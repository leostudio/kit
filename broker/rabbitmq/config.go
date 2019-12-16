package rabbitmq

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	URL       string        `env:"RABBIT_URL" envDefault:"amqp://guest:guest@127.0.0.1:5672"`
	Exchange  string        `env:"RABBIT_EXCHANGE,required"`
	DLX       string        `env:"RABBIT_DLX,required"`
	ALT       string        `env:"RABBIT_ALT,required"`
	NackDelay time.Duration `env:"RABBIT_NACK_DELAY" envDefault:"5s"`
}

func LoadConfig() (*Config, error) {
	out := new(Config)
	if err := env.Parse(out); err != nil {
		return nil, err
	}
	return out, nil
}
