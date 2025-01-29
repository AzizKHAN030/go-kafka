package email

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/azizkhan030/go-kafka/notification/config"
)

type UserRegisterEvent struct {
	ID    int
	Name  string
	Email string
}

func SendWelcomeEmail(cfg config.Config, event UserRegisterEvent) error {
	subject := "Welcome to Our Service!"
	body := fmt.Sprintf("Hi %s, \n\n Thank you for registering! We're excited to have you. \n\nBest,\nAzizbek", event.Name)

	err := sendEmail(cfg, cfg.EmailFrom, subject, body)

	if err != nil {
		log.Printf("Failed to send email to %s: %v", event.Email, err)
		return err
	}

	log.Printf("Successfully sent welcome email to %s", event.Email)
	return nil
}

func SendLoginNotification(cfg config.Config, event struct {
	ID    int
	Name  string
	Email string
}) error {
	subject := "New Login Detected!"
	body := fmt.Sprintf("Hi %s,\n\nWe detected a new login to your account.\n\nIf this was you, no action is needed. If you didn't log in, please reset your password immediately.\n\nBest,\nSecurity Team", event.Name)

	return sendEmail(cfg, event.Email, subject, body)
}

func sendEmail(cfg config.Config, recipient, subject, body string) error {
	smtpAddr := fmt.Sprintf("%s:%s", cfg.SmtpHost, cfg.SmtpPort)
	auth := smtp.PlainAuth("", cfg.SmtpUsername, cfg.SmtpPassword, cfg.SmtpHost)

	if cfg.SmtpType == "test" {
		auth = nil
	}

	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", recipient, subject, body)

	err := smtp.SendMail(smtpAddr, auth, cfg.EmailFrom, []string{recipient}, []byte(msg))
	if err != nil {
		log.Printf("Failed to send email to %s: %v", recipient, err)
		return err
	}

	log.Printf("Successfully sent email to %s", recipient)
	return nil
}
