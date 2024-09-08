package mail

import (
	"bank-authentication-system/pkg/config"

	"github.com/mailgun/mailgun-go"
	"github.com/sirupsen/logrus"
)

const (
	sender = "sabaasahban@gmail.com"
)

type Mailgun struct {
	APIKEY string
	Client *mailgun.MailgunImpl
}

func NewConnection(cfg config.MailGun) *Mailgun {
	return &Mailgun{
		APIKEY: cfg.APIKEY,
		Client: mailgun.NewMailgun(cfg.Domain, cfg.APIKEY),
	}
}

func (m *Mailgun) Send(content, subject, receiver string) error {
	message := m.Client.NewMessage(sender, subject, content, receiver)

	_, _, err := m.Client.Send(message)

	if err != nil {
		logrus.Errorf("mailgun client failed with error %s", err)

		return err
	}

	return nil
}
