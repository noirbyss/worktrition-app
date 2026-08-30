package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultEnvironment            = "local"
	defaultHTTPPort               = 8080
	defaultAllowedOrigin          = "http://localhost:5173"
	defaultUpstreamRequestTimeout = 15 * time.Second
	defaultUserServiceAddr        = "localhost:50051"
	defaultNutritionServiceAddr   = "localhost:50052"
	defaultWorkoutServiceAddr     = "localhost:50054"
	defaultAIServiceAddr          = "localhost:50056"
	defaultRefreshTokenCookieName = "refresh_token"
)

type Config struct {
	Environment            string
	HTTPPort               int
	AllowedOrigin          string
	UpstreamRequestTimeout time.Duration
	UserServiceAddr        string
	NutritionServiceAddr   string
	WorkoutServiceAddr     string
	AIServiceAddr          string
	UserJWTSecret          string
	RefreshTokenCookieName string
}

func Load() (*Config, error) {
	httpPort, err := getInt("GATEWAY_HTTP_PORT", defaultHTTPPort)
	if err != nil {
		return nil, err
	}

	upstreamRequestTimeout, err := getDuration("UPSTREAM_REQUEST_TIMEOUT", defaultUpstreamRequestTimeout)
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
		UpstreamRequestTimeout: upstreamRequestTimeout,
		UserServiceAddr:        getString("USER_SERVICE_ADDR", defaultUserServiceAddr),
		NutritionServiceAddr:   getString("NUTRITION_SERVICE_ADDR", defaultNutritionServiceAddr),
		WorkoutServiceAddr:     getString("WORKOUT_SERVICE_ADDR", defaultWorkoutServiceAddr),
		AIServiceAddr:          getString("AI_SERVICE_ADDR", defaultAIServiceAddr),
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

func getDuration(name string, defaultValue time.Duration) (time.Duration, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return defaultValue, nil
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", name, err)
	}
	if value <= 0 {
		return 0, fmt.Errorf("%s must be positive", name)
	}

	return value, nil
}
