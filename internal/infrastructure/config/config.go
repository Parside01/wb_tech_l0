package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HttpServerConfig    *HttpServerConfig    `yaml:"http_server" mapstructure:"http_server"`
	LoggerConfig        *LoggerConfig        `yaml:"logger" mapstructure:"logger"`
	PostgresConfig      *PostgresConfig      `yaml:"postgres" mapstructure:"postgres"`
	KafkaConsumerConfig *KafkaConsumerConfig `yaml:"kafka_consumer" mapstructure:"kafka_consumer"`
	MemoryCacheConfig   *MemoryCacheConfig   `yaml:"memory_cache" mapstructure:"memory_cache"`
}

type KafkaConsumerConfig struct {
	Topic    string   `yaml:"topic" mapstructure:"topic"`
	MaxBytes int      `yaml:"max_bytes" mapstructure:"max_bytes"`
	Brokers  []string `yaml:"brokers" mapstructure:"brokers"`
}

type LoggerConfig struct {
	Level    string `yaml:"level" mapstructure:"level"`
	Path     string `yaml:"path" mapstructure:"path"`
	MaxBytes int    `yaml:"max_bytes" mapstructure:"max_bytes"`
}

type HttpServerConfig struct {
	Port string `yaml:"port" mapstructure:"port"`
	Host string `yaml:"host" mapstructure:"host"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"dbname" mapstructure:"dbname"`
	SSLmode  string `yaml:"sslmode" mapstructure:"sslmode"`
}

type MemoryCacheConfig struct {
	Capacity int `yaml:"capacity" mapstructure:"capacity"`
}

var (
	Global *Config
)

func InitConfig(configPath string) error {
	v := viper.New()
	setDefaults(v)
	Global = new(Config)
	if err := readConfig(v, configPath); err != nil {
		return err
	}
	if err := v.Unmarshal(Global); err != nil {
		return err
	}
	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("http_server.port", 8080)
	v.SetDefault("http_server.host", "localhost")

	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.path", "./log/log.log")
	v.SetDefault("logger.max_bytes", 0)

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", "5432")
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "postgres")
	v.SetDefault("postgres.dbname", "postgres")
	v.SetDefault("postgres.sslmode", "disable")

	v.SetDefault("kafka_consumer.topic", "orders")
	v.SetDefault("kafka_consumer.max_bytes", 0)
	v.SetDefault("kafka_consumer.brokers", []string{"localhost:9092"})

	v.SetDefault("memory_cache.capacity", 10000)
}

func readConfig(v *viper.Viper, configPath string) error {
	v.AddConfigPath(configPath)
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
