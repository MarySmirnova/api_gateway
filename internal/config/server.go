package config

import "time"

type Server struct {
	Listen       string        `env:"GATEWAY_LISTEN" envDefault:":8080"`
	ReadTimeout  time.Duration `env:"GATEWAY_READ_TIMEOUT" envDefault:"30s"`
	WriteTimeout time.Duration `env:"GATEWAY_WRITE_TIMEOUT" envDefault:"30s"`

	NewsAddress     string `env:"NEWS_ADDRESS" envDefault:"localhost:8081"`
	CommentsAddress string `env:"COMMENTS_ADDRESS" envDefault:"localhost:8082"`
	ModerateAddress string `env:"MODERATE_ADDRESS" envDefault:"localhost:8083"`
}
