package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/slbmax/ses-weather-app/assets"
)

// ensure that the builder will be initialized with all required templates
func init() {
	b := NewBuilder()
	b.BuildConfirmationEmail(ConfirmationEmail{})
}

type EmailBuilder struct {
	confirmationTemplate        *template.Template
	notificationTemplate        *template.Template
	confirmationSuccessTemplate *template.Template
}

func (b *EmailBuilder) initialized() bool {
	return b.confirmationTemplate != nil
}

func NewBuilder() *EmailBuilder {
	templates, err := template.
		New("").
		Funcs(template.FuncMap{"title": capitalize}).
		ParseFS(assets.MailTemplates, assets.MailTemplatesDir+"/*.html")
	if err != nil {
		panic(err)
	}

	builder := &EmailBuilder{}
	for _, tmpl := range templates.Templates() {
		switch tmpl.Name() {
		case assets.TemplateConfirmation:
			builder.confirmationTemplate = tmpl
		case assets.TemplateNotification:
			builder.notificationTemplate = tmpl
		case assets.TemplateConfirmationSuccess:
			builder.confirmationSuccessTemplate = tmpl
		default:
			continue
		}
	}

	if !builder.initialized() {
		panic("builder was not initialized with all required templates")
	}

	return builder
}

func (b *EmailBuilder) BuildConfirmationEmail(message ConfirmationEmail) []byte {
	return b.build(b.confirmationTemplate, message)
}

func (b *EmailBuilder) BuildNotificationEmail(message NotificationEmail) []byte {
	return b.build(b.notificationTemplate, message)
}

func (b *EmailBuilder) BuildConfirmationSuccessEmail(message ConfirmationSuccessEmail) []byte {
	return b.build(b.confirmationSuccessTemplate, message)
}

func (b *EmailBuilder) build(tmpl *template.Template, msg any) []byte {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, msg); err != nil {
		panic(fmt.Errorf("error executing template: %w", err))
	}

	return buf.Bytes()
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
