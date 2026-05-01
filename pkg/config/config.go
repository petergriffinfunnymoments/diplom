package config

import (
	"os"
	"strconv"
	"time"
)

// Config - основная структура конфигурации приложения
type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	Logger     LoggerConfig
	Auth       AuthConfig
	AntiFraud  AntiFraudConfig
	Tokenizer  TokenizerConfig
	Adapters   AdaptersConfig
	Notification NotificationConfig
	Retry      RetryConfig
}

// AppConfig - конфигурация приложения
type AppConfig struct {
	Name           string
	Environment    string // development, staging, production
	Version        string
	Host           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig - конфигурация базы данных
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig - конфигурация Redis (для кэша и сессий)
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

// LoggerConfig - конфигурация логирования
type LoggerConfig struct {
	Level      string // debug, info, warn, error
	Format     string // json, console
	OutputPath string
	FilePath   string
	MaxSize    int // MB
	MaxBackups int
	MaxAge     int // days
}

// AuthConfig - конфигурация аутентификации
type AuthConfig struct {
	APIKeyName        string
	APIKeyHeader      string
	MasterAPIKey      string
	TokenExpiration   time.Duration
	RefreshEnabled    bool
}

// AntiFraudConfig - конфигурация антифрода
type AntiFraudConfig struct {
	Enabled            bool
	RiskThreshold      float64
	AutoBlockThreshold float64
	ReviewThreshold    float64
	MaxAmountWithoutCheck int64
	Rules              []FraudRule
}

// FraudRule - правило антифрода
type FraudRule struct {
	Name        string
	Description string
	Enabled     bool
	Parameters  map[string]interface{}
}

// TokenizerConfig - конфигурация токенизации
type TokenizerConfig struct {
	Algorithm       string // AES-256-GCM, RSA-OAEP
	KeyRotationDays int
	VaultEnabled    bool
	VaultURL        string
	VaultToken      string
}

// AdaptersConfig - конфигурация адаптеров платежных систем
type AdaptersConfig struct {
	DefaultAdapter string
	Adapters       map[string]AdapterConfig
}

// AdapterConfig - конфигурация конкретного адаптера
type AdapterConfig struct {
	Enabled      bool
	BaseURL      string
	APIKey       string
	SecretKey    string
	MerchantID   string
	Timeout      time.Duration
	RetryCount   int
	RetryDelay   time.Duration
	SandboxMode  bool
}

// NotificationConfig - конфигурация уведомлений
type NotificationConfig struct {
	Enabled       bool
	WebhookTimeout time.Duration
	EmailEnabled   bool
	SMSEnabled     bool
	RetryAttempts  int
	RetryDelay     time.Duration
}

// RetryConfig - общая конфигурация повторных попыток
type RetryConfig struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	return &Config{
		App: AppConfig{
			Name:            getEnv("APP_NAME", "payment-gateway"),
			Environment:     getEnv("APP_ENV", "development"),
			Version:         getEnv("APP_VERSION", "1.0.0"),
			Host:            getEnv("APP_HOST", "0.0.0.0"),
			Port:            getEnvAsInt("APP_PORT", 8080),
			ReadTimeout:     getEnvAsDuration("APP_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:    getEnvAsDuration("APP_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout: getEnvAsDuration("APP_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvAsInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "payment_gateway"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnvAsInt("REDIS_PORT", 6379),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
		},
		Logger: LoggerConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			OutputPath: getEnv("LOG_OUTPUT", "stdout"),
			FilePath:   getEnv("LOG_FILE_PATH", "/var/log/payment-gateway"),
			MaxSize:    getEnvAsInt("LOG_MAX_SIZE", 100),
			MaxBackups: getEnvAsInt("LOG_MAX_BACKUPS", 3),
			MaxAge:     getEnvAsInt("LOG_MAX_AGE", 30),
		},
		Auth: AuthConfig{
			APIKeyName:      getEnv("AUTH_API_KEY_NAME", "X-API-Key"),
			APIKeyHeader:    getEnv("AUTH_API_KEY_HEADER", "X-API-Key"),
			MasterAPIKey:    getEnv("AUTH_MASTER_API_KEY", "master-key-change-in-production"),
			TokenExpiration: getEnvAsDuration("AUTH_TOKEN_EXPIRATION", 24*time.Hour),
			RefreshEnabled:  getEnvAsBool("AUTH_REFRESH_ENABLED", true),
		},
		AntiFraud: AntiFraudConfig{
			Enabled:             getEnvAsBool("ANTIFRAUD_ENABLED", true),
			RiskThreshold:       getEnvAsFloat("ANTIFRAUD_RISK_THRESHOLD", 0.7),
			AutoBlockThreshold:  getEnvAsFloat("ANTIFRAUD_AUTO_BLOCK_THRESHOLD", 0.9),
			ReviewThreshold:     getEnvAsFloat("ANTIFRAUD_REVIEW_THRESHOLD", 0.5),
			MaxAmountWithoutCheck: getEnvAsInt64("ANTIFRAUD_MAX_AMOUNT_WITHOUT_CHECK", 100000), // 1000 RUB в копейках
			Rules:               []FraudRule{}, // Можно загрузить из конфига
		},
		Tokenizer: TokenizerConfig{
			Algorithm:       getEnv("TOKENIZER_ALGORITHM", "AES-256-GCM"),
			KeyRotationDays: getEnvAsInt("TOKENIZER_KEY_ROTATION_DAYS", 90),
			VaultEnabled:    getEnvAsBool("TOKENIZER_VAULT_ENABLED", false),
			VaultURL:        getEnv("TOKENIZER_VAULT_URL", ""),
			VaultToken:      getEnv("TOKENIZER_VAULT_TOKEN", ""),
		},
		Adapters: AdaptersConfig{
			DefaultAdapter: getEnv("ADAPTER_DEFAULT", "sbp"),
			Adapters: map[string]AdapterConfig{
				"sbp": {
					Enabled:     getEnvAsBool("ADAPTER_SBP_ENABLED", true),
					BaseURL:     getEnv("ADAPTER_SBP_URL", "https://sandbox.sbp.ru/api/v1"),
					APIKey:      getEnv("ADAPTER_SBP_API_KEY", ""),
					MerchantID:  getEnv("ADAPTER_SBP_MERCHANT_ID", ""),
					Timeout:     getEnvAsDuration("ADAPTER_SBP_TIMEOUT", 30*time.Second),
					RetryCount:  getEnvAsInt("ADAPTER_SBP_RETRY_COUNT", 3),
					RetryDelay:  getEnvAsDuration("ADAPTER_SBP_RETRY_DELAY", 1*time.Second),
					SandboxMode: getEnvAsBool("ADAPTER_SBP_SANDBOX", true),
				},
				"digital_ruble": {
					Enabled:     getEnvAsBool("ADAPTER_DIGITAL_RUBLE_ENABLED", true),
					BaseURL:     getEnv("ADAPTER_DIGITAL_RUBLE_URL", "https://sandbox.cbr.ru/digital-ruble/api/v1"),
					APIKey:      getEnv("ADAPTER_DIGITAL_RUBLE_API_KEY", ""),
					MerchantID:  getEnv("ADAPTER_DIGITAL_RUBLE_MERCHANT_ID", ""),
					Timeout:     getEnvAsDuration("ADAPTER_DIGITAL_RUBLE_TIMEOUT", 30*time.Second),
					RetryCount:  getEnvAsInt("ADAPTER_DIGITAL_RUBLE_RETRY_COUNT", 3),
					RetryDelay:  getEnvAsDuration("ADAPTER_DIGITAL_RUBLE_RETRY_DELAY", 1*time.Second),
					SandboxMode: getEnvAsBool("ADAPTER_DIGITAL_RUBLE_SANDBOX", true),
				},
				"card_acquirer": {
					Enabled:     getEnvAsBool("ADAPTER_CARD_ACQUIRER_ENABLED", true),
					BaseURL:     getEnv("ADAPTER_CARD_ACQUIRER_URL", "https://sandbox.acquirer.ru/api/v1"),
					APIKey:      getEnv("ADAPTER_CARD_ACQUIRER_API_KEY", ""),
					SecretKey:   getEnv("ADAPTER_CARD_ACQUIRER_SECRET_KEY", ""),
					MerchantID:  getEnv("ADAPTER_CARD_ACQUIRER_MERCHANT_ID", ""),
					Timeout:     getEnvAsDuration("ADAPTER_CARD_ACQUIRER_TIMEOUT", 30*time.Second),
					RetryCount:  getEnvAsInt("ADAPTER_CARD_ACQUIRER_RETRY_COUNT", 3),
					RetryDelay:  getEnvAsDuration("ADAPTER_CARD_ACQUIRER_RETRY_DELAY", 1*time.Second),
					SandboxMode: getEnvAsBool("ADAPTER_CARD_ACQUIRER_SANDBOX", true),
				},
			},
		},
		Notification: NotificationConfig{
			Enabled:        getEnvAsBool("NOTIFICATION_ENABLED", true),
			WebhookTimeout: getEnvAsDuration("NOTIFICATION_WEBHOOK_TIMEOUT", 10*time.Second),
			EmailEnabled:   getEnvAsBool("NOTIFICATION_EMAIL_ENABLED", false),
			SMSEnabled:     getEnvAsBool("NOTIFICATION_SMS_ENABLED", false),
			RetryAttempts:  getEnvAsInt("NOTIFICATION_RETRY_ATTEMPTS", 3),
			RetryDelay:     getEnvAsDuration("NOTIFICATION_RETRY_DELAY", 5*time.Second),
		},
		Retry: RetryConfig{
			MaxRetries:   getEnvAsInt("RETRY_MAX_RETRIES", 3),
			InitialDelay: getEnvAsDuration("RETRY_INITIAL_DELAY", 1*time.Second),
			MaxDelay:     getEnvAsDuration("RETRY_MAX_DELAY", 30*time.Second),
			Multiplier:   getEnvAsFloat("RETRY_MULTIPLIER", 2.0),
		},
	}
}

// Вспомогательные функции для получения значений из переменных окружения

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if durationValue, err := time.ParseDuration(value); err == nil {
			return durationValue
		}
	}
	return defaultValue
}
