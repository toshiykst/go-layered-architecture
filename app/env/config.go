package env

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBHost     string `envconfig:"MYSQL_HOST"`
	DBName     string `envconfig:"MYSQL_DATABASE"`
	DBUser     string `envconfig:"MYSQL_USER"`
	DBPassword string `envconfig:"MYSQL_PASSWORD"`
	DBDebug    bool   `envconfig:"MYSQL_DEBUG"`
}

func NewConfig() (*Config, error) {
	var v Config
	if err := envconfig.Process("", &v); err != nil {
		return nil, err
	}
	return &v, nil
}
