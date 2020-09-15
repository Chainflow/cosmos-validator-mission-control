package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type (
	// Telegram bot config details
	TelegramBotConfig struct {
		BotToken string `mapstructure:"tg_bot_token"`
		ChatID   int64  `mapstructure:"tg_chat_id"`
	}

	// EmailConfig
	EmailConfig struct {
		SendGridAPIToken    string `mapstructure:"sendgrid_token"`
		ReceiverMailAddress string `mapstructure:"email_address"`
	}

	// Config defines all the app configurations
	Config struct {
		ValOperatorAddress  string            `mapstructure:"val_operator_addr"`
		ValidatorHexAddress string            `mapstructure:"validator_hex_addr"`
		LCDEndpoint         string            `mapstructure:"lcd_endpoint"`
		Telegram            TelegramBotConfig `mapstructure:"telegram"`
		SendGrid            EmailConfig       `mapstructure:"sendgrid"`
		RPCEndpoint         string            `mapstructure:"rpc_endpoint"`
		ExternalRPC         string            `mapstructure:"external_rpc"`
		AlertTime1          string            `mapstructure:"alert_time1"`
		AlertTime2          string            `mapstructure:"alert_time2"`
		ValidatorName       string            `mapstructure:"validator_name"`
	}
)

// ReadConfigFromFile to read config details using viper
func ReadConfigFromFile() (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("./config/")
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config.toml: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshaling config.toml to application config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("error occurred in config validation: %v", err)
	}

	return &cfg, nil
}

// Validate config struct
func (c *Config) Validate(e ...string) error {
	v := validator.New()
	if len(e) == 0 {
		return v.Struct(c)
	}
	return v.StructExcept(c, e...)
}
