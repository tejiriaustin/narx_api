package messaging

import (
	"log"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type (
	MailjetClient struct {
		from   Author
		client *mailjet.Client
	}
	Options func(m *MailjetClient)
)

var _ Messaging = (*MailjetClient)(nil)

func NewMailjetClient(publicKey, privateKey string, opts ...Options) *MailjetClient {
	m := &MailjetClient{
		from:   Author{},
		client: mailjet.NewMailjetClient(publicKey, privateKey),
	}

	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *MailjetClient) Push(to string, msg string) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: m.from.email,
				Name:  m.from.name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
					Name:  "",
				},
			},
			TextPart: msg,
			HTMLPart: msg,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := m.client.SendMailV31(&messages)
	if err != nil {
		log.Println("Mailjet send failed: ", err)
		return err
	}
	return nil
}

func WithSenderEmail(email string) Options {
	return func(m *MailjetClient) {
		m.from.email = email
	}
}
