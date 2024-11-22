package config

import "github.com/spf13/viper"

type config struct {
	HttpServerConfig    httpServerConfig    `yaml:"http_server"`
	LoggerConfig        loggerConfig        `yaml:"logger"`
	PostgresConfig      postgresConfig      `yaml:"postgres"`
	KafkaConsumerConfig kafkaConsumerConfig `yaml:"kafka_consumer"`
	MemoryCacheConfig   memoryCacheConfig   `yaml:"memory_cache"`
}

type kafkaConsumerConfig struct {
	Topic    string   `yaml:"topic"`
	MaxBytes int      `yaml:"max_bytes"`
	Brokers  []string `yaml:"brokers"`
}

type loggerConfig struct {
	Level    string `yaml:"level"`
	Path     string `yaml:"path"`
	MaxBytes int    `yaml:"max_bytes"`
}

type httpServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type postgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	sslmode  string `yaml:"sslmode"`
}

type memoryCacheConfig struct {
	Capacity int `yaml:"capacity"`
}

var (
	Global *config
)

func InitConfig(configPath string) error {
	v := viper.New()
	setDefaults(v)
	if err := readConfig(v, configPath); err != nil {
		return err
	}

	if err := v.Unmarshal(&Global); err != nil {
		return err
	}
	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("http_server.port", 8080)
	v.SetDefault("http_server.host", "localhost")

	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.path", "./log/log.log")
	v.SetDefault("logger.max_bytes", 10e6)

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", "5432")
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "postgres")
	v.SetDefault("postgres.dbname", "postgres")
	v.SetDefault("postgres.sslmode", "disable")

	v.SetDefault("kafka_consumer.topic", "orders")
	v.SetDefault("kafka_consumer.max_bytes", 10e6)
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
