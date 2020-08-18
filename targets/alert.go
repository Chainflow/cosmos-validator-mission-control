package targets

import (
	"cosmos-validator-mission-control/alerting"
	"cosmos-validator-mission-control/config"
	"log"
	"strings"
)

// SendTelegramAlert sends the alert to telegram account
//by checking user's choice
func SendTelegramAlert(msg string, cfg *config.Config) error {
	if strings.ToUpper(cfg.EnableTelegramAlerts) == "YES" {
		if err := alerting.NewTelegramAlerter().Send(msg, cfg.Telegram.BotToken, cfg.Telegram.ChatID); err != nil {
			log.Printf("failed to send tg alert: %v", err)
			return err
		}
	}
	return nil
}

// SendEmailAlert sends alert to email account
//by checking user's choice
func SendEmailAlert(msg string, cfg *config.Config) error {
	if strings.ToUpper(cfg.EnableEmailAlerts) == "YES" {
		if err := alerting.NewEmailAlerter().Send(msg, cfg.SendGrid.Token, cfg.SendGrid.EmailAddress); err != nil {
			log.Printf("failed to send email alert: %v", err)
			return err
		}
	}
	return nil
}

// SendEmergencyEmailAlert sends alert pager duty account
func SendEmergencyEmailAlert(msg string, cfg *config.Config) error {
	if strings.ToUpper(cfg.EnableEmailAlerts) == "YES" {
		if err := alerting.NewEmailAlerter().Send(msg, cfg.SendGrid.Token, cfg.PagerdutyEmail); err != nil {
			log.Printf("failed to send email alert to pager duty: %v", err)
			return err
		}
	}
	return nil
}
