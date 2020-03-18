package alerting

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Send to send mail alert
func (e emailAlert) Send(msg, token, toEmail string) error {
	from := mail.NewEmail("ICF Tool", "icf@vitwit.com")
	subject := msg
	to := mail.NewEmail("IFC Tool", toEmail)
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
