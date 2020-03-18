package alerting

// Send Telegram alert interface
type Telegram interface {
	Send(msgText, botToken string, alerthnatID int64) error
}

type telegramAlert struct{}

// Telegram alerter
func NewTelegramAlerter() *telegramAlert {
	return &telegramAlert{}
}

// Send mail alert
type Email interface {
	Send(msg, token, toEmail string) error
}

type emailAlert struct {
}

// Mail alerter
func NewEmailAlerter() *emailAlert {
	return &emailAlert{}
}
