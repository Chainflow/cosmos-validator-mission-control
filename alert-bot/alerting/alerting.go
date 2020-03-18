package alerting

// Telegram alert interface
type Telegram interface {
	Send(msgText, botToken string, alerthnatID int64) error
}

type telegramAlert struct{}

// NewTelegramAlerter returns an instance of telegramAlert
func NewTelegramAlerter() *telegramAlert {
	return &telegramAlert{}
}

// Email to send mail alert
type Email interface {
	Send(msg, token, toEmail string) error
}

type emailAlert struct {
}

// NewEmailAlerter returns an instance for emailAlert
func NewEmailAlerter() *emailAlert {
	return &emailAlert{}
}
