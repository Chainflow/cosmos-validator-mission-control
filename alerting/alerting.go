package alerting

type Telegram interface {
	Send(msgText, botToken string, chatId int64) error
}

type telegramAlert struct{}

func NewTelegramAlerter() *telegramAlert {
	return &telegramAlert{}
}

type emailAlert struct {
}
