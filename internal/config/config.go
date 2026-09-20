package config

import (
	"errors"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

const (
	ENVProduction  = "production"
	ENVStaging     = "staging"
	ENVDevelopment = "development"
)

type Config struct {
	ENV     string `envconfig:"ENV"`
	AppName string `envconfig:"APP_NAME"`
	PORT    string `envconfig:"PORT"`

	LogLevel string `envconfig:"LOG_LEVEL"`

	JWTSecret string `envconfig:"JWT_SECRET"`
	GRPCPort  string `envconfig:"GRPC_PORT" default:"50051"`

	Database DatabaseConfig
	CORS     CORS
}

type DatabaseConfig struct {
	DBUser  string `envconfig:"DB_USER"`
	DBPass  string `envconfig:"DB_PASS"`
	DBHost  string `envconfig:"DB_HOST"`
	DBPort  string `envconfig:"DB_PORT"`
	DBName  string `envconfig:"DB_NAME"`
	SSLMode string `envconfig:"DB_SSL_MODE"`
}

type CORS struct {
	AllowHosts []string `envconfig:"ALLOW_HOSTS"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() (*Config, error) {
	var (
		err error
		cfg Config
	)
	once.Do(func() {
		err = envconfig.Process("", &cfg)
		if err == nil {
			instance = &cfg
		}
	})

	if err != nil {
		return nil, err
	}

	if instance == nil {
		return nil, errors.New("Config is nil")
	}

	return instance, nil
}
