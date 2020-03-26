package server

import (
	"chainflow-vitwit/alert-bot/alerting"
	"chainflow-vitwit/alert-bot/config"
	"log"
)

// SendTelegramAlert sends the alert to telegram chat
func SendTelegramAlert(msg string, cfg *config.Config) error {
	if err := alerting.NewTelegramAlerter().Send(msg, cfg.Telegram.BotToken, cfg.Telegram.ChatID); err != nil {
		log.Printf("failed to send telegram alert: %v", err)
		return err
	}
	return nil
}

// SendEmailAlert to send mail
func SendEmailAlert(msg string, cfg *config.Config) error {
	if err := alerting.NewEmailAlerter().Send(msg, cfg.SendGrid.Token, cfg.SendGrid.ToEmail); err != nil {
		log.Printf("failed to send email alert: %v", err)
		return err
	}
	return nil
}
