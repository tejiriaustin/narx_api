package messaging

import (
	"log"
	"net/smtp"
)

var (
	defaultAddr           = "smtp.gmail.com:587"
	_           Messaging = (*SMTP)(nil)
)

type SMTP struct {
	From string
	To   string
	Pass string
	Host string
}

func NewSMTP(pass, from, host string) *SMTP {
	if host == "" {
		host = defaultAddr
	}
	return &SMTP{Pass: pass, Host: host, From: from}
}

func (s *SMTP) Push(to string, msg string) error {
	auth := smtp.PlainAuth("", s.From, s.From, "smtp.gmail.com")
	err := smtp.SendMail(s.Host, auth, s.From, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Println("mail sent successfully")
	return nil
}
