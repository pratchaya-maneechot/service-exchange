package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"environment" validate:"required,oneof=development staging production"`
	Version     string         `mapstructure:"version" validate:"required"`
	Server      ServerConfig   `mapstructure:"server" validate:"required"`
	Database    DatabaseConfig `mapstructure:"database" validate:"required"`
	Logging     LoggingConfig  `mapstructure:"logging" validate:"required"`
	Metrics     MetricsConfig  `mapstructure:"metrics" validate:"required"`
	Security    SecurityConfig `mapstructure:"security" validate:"required"`
}

type ServerConfig struct {
	Address               string        `mapstructure:"address" validate:"required"`
	Port                  int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	ReadTimeout           time.Duration `mapstructure:"read_timeout" validate:"gte=0"`
	WriteTimeout          time.Duration `mapstructure:"write_timeout" validate:"gte=0"`
	ShutdownTimeout       time.Duration `mapstructure:"shutdown_timeout" validate:"gte=0"`
	MaxConnections        int           `mapstructure:"max_connections" validate:"gte=0"`
	MaxConcurrentStreams  uint32        `mapstructure:"max_concurrent_streams" validate:"gte=0"`
	EnableReflection      bool          `mapstructure:"enable_reflection"`
	EnableHealthCheck     bool          `mapstructure:"enable_health_check"`
	ConnectionTimeout     time.Duration `mapstructure:"connection_timeout" validate:"gte=0"`
	KeepaliveTime         time.Duration `mapstructure:"keepalive_time" validate:"gte=0"`
	KeepaliveTimeout      time.Duration `mapstructure:"keepalive_timeout" validate:"gte=0"`
	MaxConnectionIdle     time.Duration `mapstructure:"max_connection_idle" validate:"gte=0"`
	MaxConnectionAge      time.Duration `mapstructure:"max_connection_age" validate:"gte=0"`
	MaxConnectionAgeGrace time.Duration `mapstructure:"max_connection_age_grace" validate:"gte=0"`
	MaxRecvMsgSize        int           `mapstructure:"max_recv_msg_size" validate:"gte=0"`
	MaxSendMsgSize        int           `mapstructure:"max_send_msg_size" validate:"gte=0"`
}

type DatabaseConfig struct {
	URL               string        `mapstructure:"url" validate:"required,url"`
	MaxOpenConns      int           `mapstructure:"max_open_conns" validate:"gte=0"`
	MaxIdleConns      int           `mapstructure:"max_idle_conns" validate:"gte=0"`
	ConnMaxLifetime   time.Duration `mapstructure:"conn_max_lifetime" validate:"gte=0"`
	ConnMaxIdleTime   time.Duration `mapstructure:"conn_max_idle_time" validate:"gte=0"`
	HealthCheckPeriod time.Duration `mapstructure:"health_check_period" validate:"gte=0"`
	ConnectionTimeout time.Duration `mapstructure:"connection_timeout" validate:"gte=0"`
	QueryTimeout      time.Duration `mapstructure:"query_timeout" validate:"gte=0"`
}

type LoggingConfig struct {
	Level      string `mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Format     string `mapstructure:"format" validate:"required,oneof=console json"`
	Output     string `mapstructure:"output" validate:"required,oneof=stdout stderr"`
	Structured bool   `mapstructure:"structured"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled" validate:"required"`
	Address string `mapstructure:"address" validate:"required_if=Enabled true"`
	Path    string `mapstructure:"path" validate:"required_if=Enabled true"`
}

type SecurityConfig struct {
	EnableTLS      bool     `mapstructure:"enable_tls"`
	TLSCertFile    string   `mapstructure:"tls_cert_file" validate:"required_if=EnableTLS true"`
	TLSKeyFile     string   `mapstructure:"tls_key_file" validate:"required_if=EnableTLS true"`
	TrustedProxies []string `mapstructure:"trusted_proxies" validate:"omitempty"`
	RateLimitRPS   int      `mapstructure:"rate_limit_rps" validate:"gte=0"`
	RateLimitBurst int      `mapstructure:"rate_limit_burst" validate:"gte=0"`
	MaxRequestSize int64    `mapstructure:"max_request_size" validate:"gte=0"`
	EnableCORS     bool     `mapstructure:"enable_cors"`
	AllowedOrigins []string `mapstructure:"allowed_origins" validate:"required_if=EnableCORS true"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, continuing with defaults and env vars")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("./apps/tasks/internal/config")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	validate := validator.New()

	if err := validate.Struct(c); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var errs []string
			for _, e := range validationErrors {
				errs = append(errs, fmt.Sprintf("field %s: %s", e.Namespace(), e.Tag()))
			}
			return fmt.Errorf("validation errors: %s", strings.Join(errs, ", "))
		}
		return fmt.Errorf("validation failed: %w", err)
	}

	if c.Environment == "production" {
		if c.Server.EnableReflection {
			return fmt.Errorf("server.enable_reflection must be false in production")
		}
		if !c.Security.EnableTLS {
			return fmt.Errorf("security.enable_tls must be true in production")
		}
	}

	if c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		return fmt.Errorf("database.max_idle_conns must not exceed database.max_open_conns")
	}

	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
