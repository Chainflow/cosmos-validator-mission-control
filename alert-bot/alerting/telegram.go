package alerting

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (t telegramAlert) Send(msgText, botToken string, chatId int64) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}
	bot.Debug = true
	msg := tgbotapi.NewMessage(chatId, "")
	msg.Text = msgText
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
