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
	Port string
}

func NewSMTP(pass, from, host, port string) *SMTP {
	return &SMTP{Pass: pass, Host: host, From: from, Port: port}
}

func (s *SMTP) Push(to string, msg string) error {
	auth := smtp.PlainAuth("", "tejiriaustin123@gmail.com", "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", "smtp.gmail.com:587")

	err := smtp.SendMail(defaultAddr, auth, s.From, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	log.Println("mail sent successfully")
	return nil
}
