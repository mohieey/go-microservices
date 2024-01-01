package main

import (
	"fmt"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain     string
	Host       string
	Port       int
	Username   string
	Password   string
	Encryption string
	From       string
	FromName   string
}

type Message struct {
	From       string
	FromName   string
	To         string
	Subject    string
	Attchments []string
	Text       string
}

func (m *Mail) Send(msg Message) error {
	if msg.From == "" {
		msg.From = m.From
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	server := m.getServer()
	client, err := server.Connect()
	if err != nil {
		return fmt.Errorf("error connecting to server %w", err)
	}

	email := m.buildEmail(msg)

	err = email.Send(client)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

func (m *Mail) getEncryption(encryption string) mail.Encryption {
	switch encryption {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

func (m *Mail) getServer() mail.SMTPServer {
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	return *server
}

func (m *Mail) buildEmail(msg Message) *mail.Email {
	email := mail.NewMSG()
	email.SetFrom(m.From)
	email.AddTo(msg.To)
	email.SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, msg.Text)

	for _, attachment := range msg.Attchments {
		email.AddAttachment(attachment)
	}
	return email
}
