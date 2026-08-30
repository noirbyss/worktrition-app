package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	defaultEnvironment            = "local"
	defaultHTTPPort               = 8080
	defaultAllowedOrigin          = "http://localhost:5173"
	defaultUserServiceAddr        = "localhost:50051"
	defaultNutritionServiceAddr   = "localhost:50052"
	defaultRefreshTokenCookieName = "refresh_token"
)

type Config struct {
	Environment            string
	HTTPPort               int
	AllowedOrigin          string
	UserServiceAddr        string
	NutritionServiceAddr   string
	UserJWTSecret          string
	RefreshTokenCookieName string
}

func Load() (*Config, error) {
	httpPort, err := getInt("GATEWAY_HTTP_PORT", defaultHTTPPort)
	if err != nil {
		return nil, err
	}

	userJWTSecret := strings.TrimSpace(os.Getenv("USER_JWT_SECRET"))
	if userJWTSecret == "" {
		return nil, fmt.Errorf("USER_JWT_SECRET is required")
	}

	return &Config{
		Environment:            getString("ENV", defaultEnvironment),
		HTTPPort:               httpPort,
		AllowedOrigin:          getString("ALLOWED_ORIGIN", defaultAllowedOrigin),
		UserServiceAddr:        getString("USER_SERVICE_ADDR", defaultUserServiceAddr),
		NutritionServiceAddr:   getString("NUTRITION_SERVICE_ADDR", defaultNutritionServiceAddr),
		UserJWTSecret:          userJWTSecret,
		RefreshTokenCookieName: getString("REFRESH_TOKEN_COOKIE_NAME", defaultRefreshTokenCookieName),
	}, nil
}

func (c *Config) HTTPAddress() string {
	return net.JoinHostPort("", strconv.Itoa(c.HTTPPort))
}

func (c *Config) SecureCookies() bool {
	env := strings.ToLower(strings.TrimSpace(c.Environment))
	return env == "prod" || env == "production"
}

func getString(name, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return defaultValue
	}

	return value
}

func getInt(name string, defaultValue int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", name, err)
	}
	if value <= 0 {
		return 0, fmt.Errorf("%s must be positive", name)
	}

	return value, nil
}
