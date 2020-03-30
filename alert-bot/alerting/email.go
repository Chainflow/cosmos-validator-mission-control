package alerting

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Send email alert
func (e emailAlert) Send(msg, token, toEmail string) error {
	from := mail.NewEmail("The Validator Voting Bot", "emailbot@chainflow.io")
	subject := msg
	to := mail.NewEmail("to", toEmail)
	plainTextContent := msg
	htmlContent := msg
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(token)
	_, err := client.Send(message)
	if err != nil {
		return err
	}
	return nil
}
