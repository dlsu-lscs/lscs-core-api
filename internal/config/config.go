package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port               int
	GoEnv              string
	ServerIdleTimeout  time.Duration
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration

	// Database
	DBHost     string
	DBPort     string
	DBDatabase string
	DBUsername string
	DBPassword string

	// Authentication (API Keys)
	JWTSecret         string
	GoogleClientID    string
	JWTDevExpiryDays  int
	JWTProdExpiryDays int

	// OAuth (Web UI Sessions)
	GoogleClientSecret string
	OAuthRedirectURL   string

	// Session (Web UI)
	SessionSecret           string
	SessionDuration         int // seconds, default 24 hours
	SessionRememberDuration int // seconds, default 30 days

	// CORS
	AllowedOrigins []string

	// Logging
	LogLevel string

	// S3/Garage Storage
	S3Endpoint        string
	S3Bucket          string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3Region          string
}

var cfg *Config

// Load reads environment variables and returns a validated Config.
// It should be called once at application startup.
func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		// Server
		Port:               getEnvInt("PORT", 8080),
		GoEnv:              getEnv("GO_ENV", "development"),
		ServerIdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", time.Minute),
		ServerReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		ServerWriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),

		// Database
		DBHost:     getEnvRequired("DB_HOST"),
		DBPort:     getEnvRequired("DB_PORT"),
		DBDatabase: getEnvRequired("DB_DATABASE"),
		DBUsername: getEnvRequired("DB_USERNAME"),
		DBPassword: getEnvRequired("DB_PASSWORD"),

		// Authentication (API Keys)
		JWTSecret:         getEnvRequired("JWT_SECRET"),
		GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
		JWTDevExpiryDays:  getEnvInt("JWT_DEV_EXPIRY_DAYS", 30),
		JWTProdExpiryDays: getEnvInt("JWT_PROD_EXPIRY_DAYS", 365),

		// OAuth (Web UI Sessions)
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		OAuthRedirectURL:   getEnv("OAUTH_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),

		// Session (Web UI)
		SessionSecret:           getEnv("SESSION_SECRET", ""),
		SessionDuration:         getEnvInt("SESSION_DURATION", 86400),            // 24 hours
		SessionRememberDuration: getEnvInt("SESSION_REMEMBER_DURATION", 2592000), // 30 days

		// CORS
		AllowedOrigins: getEnvList("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),

		// Logging
		LogLevel: getEnv("LOG_LEVEL", "info"),

		// S3/Garage Storage
		S3Endpoint:        getEnv("S3_ENDPOINT", ""),
		S3Bucket:          getEnv("S3_BUCKET", "lscs-core"),
		S3AccessKeyID:     getEnv("S3_ACCESS_KEY", ""),
		S3SecretAccessKey: getEnv("S3_SECRET_KEY", ""),
		S3Region:          getEnv("S3_REGION", "garage"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Get returns the current config. Panics if Load() hasn't been called.
func Get() *Config {
	if cfg == nil {
		panic("config.Load() must be called before config.Get()")
	}
	return cfg
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.GoEnv == "development" || c.GoEnv == "dev"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.GoEnv == "production" || c.GoEnv == "prod"
}

// FrontendURL returns the first allowed origin (used for redirects after OAuth)
func (c *Config) FrontendURL() string {
	if len(c.AllowedOrigins) > 0 {
		return c.AllowedOrigins[0]
	}
	return "http://localhost:3000"
}

// DSN returns the MySQL connection string
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUsername, c.DBPassword, c.DBHost, c.DBPort, c.DBDatabase)
}

func (c *Config) validate() error {
	var missing []string

	if c.DBHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.DBPort == "" {
		missing = append(missing, "DB_PORT")
	}
	if c.DBDatabase == "" {
		missing = append(missing, "DB_DATABASE")
	}
	if c.DBUsername == "" {
		missing = append(missing, "DB_USERNAME")
	}
	if c.DBPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	// validate log level
	validLevels := map[string]bool{
		"trace": true, "debug": true, "info": true,
		"warn": true, "error": true, "fatal": true, "panic": true,
	}
	if !validLevels[strings.ToLower(c.LogLevel)] {
		return fmt.Errorf("invalid LOG_LEVEL: %s (must be one of: trace, debug, info, warn, error, fatal, panic)", c.LogLevel)
	}

	return nil
}

// helper functions for reading environment variables

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvRequired(key string) string {
	return os.Getenv(key)
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvList(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
