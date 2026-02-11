// config for the backend app

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// structs - config

// parent struct for all the config
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

// load function for loading .env file

func Load() (*Config, error) {

	// loading .env file

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found , using environment variables")
	}

	// getting jwt secret access exipry duration (time limit it will stay relevant)

	accessExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_EXPIRY: %w", err)
	}

	// refreshing the jwt access expiry

	refreshExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRY: %w", err)
	}

	// returning the overall config

	return &Config{

			// server config

			Server: ServerConfig{
				Port: getEnv("PORT", "8080"),
				Env:  getEnv("ENV", "development"),
			},

			// db config

			Database: DatabaseConfig{
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnv("DB_PORT", "5432"),
				User:     getEnv("DB_USER", "postgres"),
				Password: getEnv("DB_PASSWORD", "postgres"),
				DBName:   getEnv("DB_NAME", "invoicego"),
				SSLMode:  getEnv("DB_SSLMODE", "disable"),
			},

			JWT: JWTConfig{
				Secret:        getEnv("JWT_SECRET", "secret-key-production"),
				AccessExpiry:  accessExpiry,
				RefreshExpiry: refreshExpiry,
			},

			CORS: CORSConfig{
				AllowedOrigins: []string{
					getEnv("ALLOWED_ORIGINS", "https://localhost:3000"),
				},
			},
		},
		nil
}

// some functions  to use

// string connection function returns resulting string

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// function to get env

// key = secrect key , defaultValue - time duration

func getEnv(key, defaultValue string) string {

	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

// getting the env as integer

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")

	// converting string key to integer

	if value, err := strconv.Atoi(valueStr); err != nil {
		return value
	}

	return defaultValue
}
