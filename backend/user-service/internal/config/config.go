package config

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `env:"ENV" env-default:"local"`
	GRPC        GRPCConfig
	Database    DatabaseConfig
	Log         LogConfig
	JWT         JWTConfig
}

type GRPCConfig struct {
	Port int `env:"USER_GRPC_PORT" env-default:"50051"`
}

func (c GRPCConfig) Address() string {
	return net.JoinHostPort("", strconv.Itoa(c.Port))
}

type DatabaseConfig struct {
	Host     string `env:"USER_DB_HOST" env-required:"true"`
	Port     int    `env:"USER_DB_PORT" env-default:"5432"`
	User     string `env:"USER_DB_USER" env-required:"true"`
	Password string `env:"USER_DB_PASSWORD" env-required:"true"`
	Name     string `env:"USER_DB_NAME" env-required:"true"`
	SSLMode  string `env:"USER_DB_SSL_MODE" env-default:"disable"`
}

func (c DatabaseConfig) DSN() string {
	dsn := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:   c.Name,
	}

	query := dsn.Query()
	query.Set("sslmode", c.SSLMode)
	dsn.RawQuery = query.Encode()

	return dsn.String()
}

type LogConfig struct {
	Level string `env:"USER_LOG_LEVEL" env-default:"info"`
}

type JWTConfig struct {
	Secret string        `env:"USER_JWT_SECRET" env-required:"true"`
	TTL    time.Duration `env:"USER_JWT_TTL" env-default:"24h"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}

	return cfg, nil
}
