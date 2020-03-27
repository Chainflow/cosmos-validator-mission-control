package config

import (
	"os"

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

	//Scraper time interval
	Scraper struct {
		Rate          string `mapstructure:"rate"`
		Port          string `mapstructure:"port"`
		ValidatorRate string `mapstructure:"validator_rate"`
	}

	//InfluxDB details
	InfluxDB struct {
		Port     string `mapstructure:"port"`
		Database string `mapstructure:"database"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}

	//Config
	Config struct {
		NodeURL               string   `mapstructure:"node_url"`
		OperatorAddress       string   `mapstructure:"operator_addr"`
		AccountAddress        string   `mapstructure:"account_addr"`
		ValidatorAddress      string   `mapstructure:"validator_addr"`
		LCDEndpoint           string   `mapstructure:"lcd_endpoint"`
		VotingPowerThreshold  int64    `mapstructure:"voting_power_threshold"`
		NumPeersThreshold     int64    `mapstructure:"num_peers_threshold"`
		Scraper               Scraper  `mapstructure:"scraper"`
		Telegram              Telegram `mapstructure:"telegram"`
		SendGrid              SendGrid `mapstructure:"sendgrid"`
		InfluxDB              InfluxDB `mapstructure:"influxdb"`
		RPCEndpoint           string   `mapstructure:"rpc_endpoint"`
		ExternalRPC           string   `mapstructure:"external_rpc"`
		MissedBlocksThreshold int64    `mapstructure:"missed_blocks_threshold"`
		AlertTime1            string   `mapstructure:"alert_time1"`
		AlertTime2            string   `mapstructure:"alert_time2"`
	}
)

//ReadFromEnv to read env details
func ReadFromEnv() *Config {
	return &Config{
		NodeURL:         getEnv("NODE_URL", ""),
		OperatorAddress: getEnv("OPERATOR_ADDR", ""),
		AccountAddress:  getEnv("ACCOUNT_ADDR", ""),
		LCDEndpoint:     getEnv("LCD_ENDPOINT", ""),
		RPCEndpoint:     getEnv("EXTERNAL_RPC", ""),
	}
}

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

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else if defaultVal == "" {
		log.Fatalf("environment variable %s cannot have a nil value", key)
	}
	return defaultVal
}

//Validate config struct
func (c *Config) Validate(e ...string) error {
	v := validator.New()
	if len(e) == 0 {
		return v.Struct(c)
	}
	return v.StructExcept(c, e...)
}
