package alerting

// Telegram is an interface definition for telegram related actions
// like Send telegram alert
type Telegram interface {
	Send(msgText, botToken string, alerthnatID int64) error
}

type telegramAlert struct{}

// NewTelegramAlerter returns a new instance for telegramAlert
func NewTelegramAlerter() *telegramAlert {
	return &telegramAlert{}
}

// Email is an interface definition for email actions
// like Send email alert
type Email interface {
	Send(msg, token, toEmail string) error
}

type emailAlert struct {}

//NewEmailAlerter returns emailAlert
func NewEmailAlerter() *emailAlert {
	return &emailAlert{}
}
