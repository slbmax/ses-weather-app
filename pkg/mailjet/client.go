package mailjet

import (
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type From struct {
	Email string
	Name  string
}

type Client struct {
	mailjet *mailjet.Client
	from    From
}

func NewClient(mailjetKey, mailjetSecret string, from From) *Client {
	return &Client{
		mailjet: mailjet.NewMailjetClient(mailjetKey, mailjetSecret),
		from:    from,
	}
}

func (c *Client) Send(to, subject string, bodyHtml string) error {
	msgInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: c.from.Email,
				Name:  c.from.Name,
			},
			To: &mailjet.RecipientsV31{
				{
					Email: to,
				},
			},
			Subject:  subject,
			HTMLPart: bodyHtml,
		},
	}

	msg := &mailjet.MessagesV31{Info: msgInfo}
	if _, err := c.mailjet.SendMailV31(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
