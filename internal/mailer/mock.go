package mailer

import "fmt"

type MockMailer struct {
	builder *EmailBuilder
}

func NewMockMailer() Mailer {
	return &MockMailer{
		builder: NewBuilder(),
	}
}

func (m *MockMailer) SendConfirmationEmail(_ string, email ConfirmationEmail) error {
	m.builder.BuildConfirmationEmail(email)
	fmt.Println("email sent")

	return nil
}

func (m *MockMailer) SendNotificationEmail(_ string, email NotificationEmail) error {
	m.builder.BuildNotificationEmail(email)
	fmt.Println("email sent")

	return nil
}

func (m *MockMailer) SendConfirmationSuccessEmail(_ string, email ConfirmationSuccessEmail) error {
	m.builder.BuildConfirmationSuccessEmail(email)
	fmt.Println("email sent")

	return nil
}
