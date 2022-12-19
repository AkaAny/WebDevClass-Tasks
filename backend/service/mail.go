package service

import (
	"gopkg.in/gomail.v2"
	"webdevclass-tasks/config"
)

type MailClient struct {
	dialer *gomail.Dialer
}

func NewFromMailConfig(mailConfig *config.MailConfig) *MailClient {
	var dialer = gomail.NewDialer(mailConfig.ServerHost, mailConfig.ServerPort,
		mailConfig.UserName, mailConfig.Password)
	return &MailClient{dialer: dialer}
}

func (c *MailClient) SendMail(messages ...*gomail.Message) error {
	return c.dialer.DialAndSend(messages...)
}
