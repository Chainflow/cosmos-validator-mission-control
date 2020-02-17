package targets

import (
	"chainflow-vitwit/alerting"
	"chainflow-vitwit/config"
	"log"
)

func SendTelegramAlert(msg string, cfg *config.Config) error {
	if err := alerting.NewTelegramAlerter().Send("Gaiad is not running", cfg.Telegram.BotToken, cfg.Telegram.ChatId); err != nil {
		log.Printf("failed to send tg alert: %v", err)
		return err
	}
	return nil
}
