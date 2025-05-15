package assets

import "embed"

const (
	MailTemplatesDir            = "templates/mail"
	TemplateConfirmation        = "confirmation.html"
	TemplateNotification        = "notification.html"
	TemplateConfirmationSuccess = "confirmation_success.html"
)

//go:embed migrations/*.sql
var Migrations embed.FS

//go:embed templates/mail/*.html
var MailTemplates embed.FS
