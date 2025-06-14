package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Log      LogConfig      `mapstructure:"log"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	RSAKeys  RSAKeysConfig  `mapstructure:"rsa_keys"`
	Services ServicesConfig `mapstructure:"services"`
}

type ServerConfig struct {
	HTTPPort                string        `mapstructure:"http_port"`
	GRPCPort                string        `mapstructure:"grpc_port"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout_sec"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type PostgresConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"db_name"`
	SSLMode      string `mapstructure:"ssl_mode"`
	PoolMaxConns int    `mapstructure:"pool_max_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	AccessTokenSecret      string        `mapstructure:"access_token_secret"` // Used for HS256 if applicable, or general secret
	RefreshTokenSecret     string        `mapstructure:"refresh_token_secret"`// Used for HS256 if applicable, or general secret
	AccessTokenExpiryMin   time.Duration `mapstructure:"access_token_expiry_min"`
	RefreshTokenExpiryDays time.Duration `mapstructure:"refresh_token_expiry_days"`
	Issuer                 string        `mapstructure:"issuer"`
	Audience               string        `mapstructure:"audience"`
	PrivateKeyPath         string        `mapstructure:"private_key_path"` // For RS256
	PublicKeyPath          string        `mapstructure:"public_key_path"`  // For RS256
}

type RSAKeysConfig struct {
	PrivateKeyPath string `mapstructure:"private_path"`
	PublicKeyPath  string `mapstructure:"public_path"`
}

type KafkaConfig struct {
	Brokers     string `mapstructure:"brokers"`
	TopicPrefix string `mapstructure:"topic_prefix"`
}

type ServicesConfig struct {
	NotificationServiceTopic string `mapstructure:"notification_service_topic"`
}


var AppConfig Config

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)     // path to look for the config file in
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs") // For tests

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file: %s. Using defaults or env vars.", err)
		// Allow running without a config file if all settings are provided by env vars
		// For critical settings not found, return error or handle appropriately
	}

	// Unmarshal the config into AppConfig
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return nil, err
	}

    // Convert durations
    AppConfig.Server.GracefulShutdownTimeout = AppConfig.Server.GracefulShutdownTimeout * time.Second
    AppConfig.JWT.AccessTokenExpiryMin = AppConfig.JWT.AccessTokenExpiryMin * time.Minute
    AppConfig.JWT.RefreshTokenExpiryDays = AppConfig.JWT.RefreshTokenExpiryDays * 24 * time.Hour


	// Update JWT paths if they are in the old JWTConfig location for backward compatibility
    // The new rsa_keys struct is preferred
    if AppConfig.JWT.PrivateKeyPath != "" && AppConfig.RSAKeys.PrivateKeyPath == "" {
        AppConfig.RSAKeys.PrivateKeyPath = AppConfig.JWT.PrivateKeyPath
    }
    if AppConfig.JWT.PublicKeyPath != "" && AppConfig.RSAKeys.PublicKeyPath == "" {
        AppConfig.RSAKeys.PublicKeyPath = AppConfig.JWT.PublicKeyPath
    }


	log.Println("Configuration loaded successfully.")
	return &AppConfig, nil
}
