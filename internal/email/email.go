package email

import (
	"fmt"

	gomail "github.com/wneessen/go-mail"
)

// Sender defines the interface for sending emails.
type Sender interface {
	SendInvoice(to, subject, body string, pdfPath string) error
}

// SMTPConfig holds SMTP configuration.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// SMTPSender sends emails via SMTP.
type SMTPSender struct {
	Config SMTPConfig
}

// NewSMTPSender creates a new SMTP sender.
func NewSMTPSender(cfg SMTPConfig) *SMTPSender {
	return &SMTPSender{Config: cfg}
}

// SendInvoice sends an invoice email with a PDF attachment.
func (s *SMTPSender) SendInvoice(to, subject, body string, pdfPath string) error {
	m := gomail.NewMsg()
	if err := m.From(s.Config.From); err != nil {
		return fmt.Errorf("set from: %w", err)
	}
	if err := m.To(to); err != nil {
		return fmt.Errorf("set to: %w", err)
	}
	m.Subject(subject)
	m.SetBodyString(gomail.TypeTextPlain, body)

	if pdfPath != "" {
		m.AttachFile(pdfPath)
	}

	c, err := gomail.NewClient(s.Config.Host,
		gomail.WithPort(s.Config.Port),
		gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
		gomail.WithUsername(s.Config.Username),
		gomail.WithPassword(s.Config.Password),
	)
	if err != nil {
		return fmt.Errorf("create mail client: %w", err)
	}
	if err := c.DialAndSend(m); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}
