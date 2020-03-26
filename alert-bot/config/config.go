package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type (
	//Telegram bot details struct
	Telegram struct {
		BotToken string `mapstructure:"bot_token"`
		ChatID   int64  `mapstructure:"chat_id"`
	}

	//SendGrid tokens
	SendGrid struct {
		Token   string `mapstructure:"token"`
		ToEmail string `mapstructure:"to_email"`
	}

	//Config
	Config struct {
		OperatorAddress  string   `mapstructure:"operator_addr"`
		AccountAddress   string   `mapstructure:"account_addr"`
		ValidatorAddress string   `mapstructure:"validator_addr"`
		LCDEndpoint      string   `mapstructure:"lcd_endpoint"`
		Telegram         Telegram `mapstructure:"telegram"`
		SendGrid         SendGrid `mapstructure:"sendgrid"`
		RPCEndpoint      string   `mapstructure:"rpc_endpoint"`
		ExternalRPC      string   `mapstructure:"external_rpc"`
		AlertTime1       string   `mapstructure:"alert_time1"`
		AlertTime2       string   `mapstructure:"alert_time2"`
	}
)

//ReadFromFile to read config details using viper
func ReadFromFile() (*Config, error) {
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
