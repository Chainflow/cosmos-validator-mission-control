package alerting

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// Send allows bot to send a telegram alert to the configured chatID
func (t telegramAlert) Send(msgText, botToken string, chatID int64) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	bot.Debug = true
	msg := tgbotapi.NewMessage(chatID, "")
	msg.Text = msgText

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
