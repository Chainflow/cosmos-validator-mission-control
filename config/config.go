package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
	"os"
)

type (
	Config struct {
		NodeURL         string `json:"node_url" toml:"node_url"`
		OperatorAddress string `json:"operator_addr" toml:"operator_addr"`
		AccountAddress  string `json:"account_addr" toml:"account_addr"`
		LCDEndpoint     string `json:"lcd_endpoint" toml:"lcd_endpoint"`
	}
)

func ReadFromEnv() *Config {
	return &Config{
		NodeURL:         getEnv("NODE_URL", ""),
		OperatorAddress: getEnv("OPERATOR_ADDR", ""),
		AccountAddress:  getEnv("ACCOUNT_ADDR", ""),
		LCDEndpoint:     getEnv("LCD_ENDPOINT", ""),
	}
}

func ReadFromTomlFile() (*Config, error) {
	viper.AddConfigPath("./../")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error while reading config.toml: %v", err)
	}

	var cfg Config
	if err = viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshaling config.toml to application config: %v", err)
	}

	if err = cfg.Validate(); err != nil {
		log.Fatalf("error occurred in config validation: %v", err)
	}

	return &cfg, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else if defaultVal == "" {
		log.Fatalf("environment variable %s cannot have a nil value", key)
	}
	return defaultVal
}

func (c *Config) Validate(e ...string) error {
	v := validator.New()
	if len(e) == 0 {
		return v.Struct(c)
	}
	return v.StructExcept(c, e...)
}
