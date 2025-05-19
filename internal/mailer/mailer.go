package mailer

import (
	"fmt"

	"github.com/slbmax/ses-weather-app/pkg/mailjet"
)

const (
	EmailSubjectConfirmation        = "Weather App - Confirm your email"
	EmailSubjectNotification        = "Weather App - Weather Notification"
	EmailSubjectConfirmationSuccess = "Weather App - Confirmation Success"
)

type Mailer interface {
	SendConfirmationEmail(to string, email ConfirmationEmail) error
	SendNotificationEmail(to string, email NotificationEmail) error
	SendConfirmationSuccessEmail(to string, message ConfirmationSuccessEmail) error
}

type mailer struct {
	builder *EmailBuilder
	client  *mailjet.Client
}

func NewMailer(client *mailjet.Client) Mailer {
	return &mailer{
		builder: NewBuilder(),
		client:  client,
	}
}

func (m *mailer) sendEmail(to, subject string, body []byte) error {
	if err := m.client.Send(to, subject, string(body)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (m *mailer) SendConfirmationEmail(to string, email ConfirmationEmail) error {
	return m.sendEmail(to, EmailSubjectConfirmation, m.builder.BuildConfirmationEmail(email))
}

func (m *mailer) SendNotificationEmail(to string, email NotificationEmail) error {
	return m.sendEmail(to, EmailSubjectNotification, m.builder.BuildNotificationEmail(email))
}

func (m *mailer) SendConfirmationSuccessEmail(to string, email ConfirmationSuccessEmail) error {
	return m.sendEmail(to, EmailSubjectConfirmationSuccess, m.builder.BuildConfirmationSuccessEmail(email))
}
