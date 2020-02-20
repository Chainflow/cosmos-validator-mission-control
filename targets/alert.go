package targets

import (
	"chainflow-vitwit/alerting"
	"chainflow-vitwit/config"
	"log"
)

func SendTelegramAlert(msg string, cfg *config.Config) error {
	if err := alerting.NewTelegramAlerter().Send(msg, cfg.Telegram.BotToken, cfg.Telegram.ChatId); err != nil {
		log.Printf("failed to send tg alert: %v", err)
		return err
	}
	return nil
}

func SendEmailAlert(msg string, cfg *config.Config) error {
	if err := alerting.NewEmailAlerter().Send(msg, cfg.SendGrid.Token, cfg.SendGrid.ToEmail); err != nil {
		log.Printf("failed to send email alert: %v", err)
		return err
	}
	return nil
}
