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
	auth := smtp.PlainAuth("", cfg.SmtpUsername, cfg.SmtpPassword, cfg.SmtpHost)

	if cfg.SmtpType == "test" {
		auth = nil
	}

	smtpAddr := fmt.Sprintf("%s:%s", cfg.SmtpHost, cfg.SmtpPort)
	subject := "Welcome to Our Service!"
	body := fmt.Sprintf("Hi %s, \n\n Thank you for registering! We're excited to have you. \n\nBest,\nAzizbek", event.Name)
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", event.Email, subject, body)

	err := smtp.SendMail(smtpAddr, auth, cfg.EmailFrom, []string{event.Email}, []byte(msg))

	if err != nil {
		log.Printf("Failed to send email to %s: %v", event.Email, err)
		return err
	}

	log.Printf("Successfully sent welcome email to %s", event.Email)
	return nil
}
