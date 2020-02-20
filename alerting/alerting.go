package alerting

type Telegram interface {
	Send(msgText, botToken string, chatId int64) error
}

type telegramAlert struct{}

func NewTelegramAlerter() *telegramAlert {
	return &telegramAlert{}
}

type Email interface {
	Send(msg, token, toEmail string) error
}

type emailAlert struct {
}

func NewEmailAlerter() *emailAlert {
	return &emailAlert{}
}
