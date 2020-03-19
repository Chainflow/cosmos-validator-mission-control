package alerting

//Telegram to send telegram alert interface
type Telegram interface {
	Send(msgText, botToken string, chatID int64) error
}

type telegramAlert struct{}

//NewTelegramAlerter returns telegram alerter
func NewTelegramAlerter() *telegramAlert {
	return &telegramAlert{}
}

//Email to send mail alert
type Email interface {
	Send(msg, token, toEmail string) error
}

type emailAlert struct {
}

//NewEmailAlerter returns email alert
func NewEmailAlerter() *emailAlert {
	return &emailAlert{}
}
