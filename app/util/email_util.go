package util

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type EmailUtil interface {
	SendEmail(to []string, cc []string, subject, message string) error
}

func NewEmailUtil() EmailUtil {
	return &AuthUtilImpl{}
}

type AuthUtilImpl struct{}

func (a AuthUtilImpl) SendEmail(to []string, cc []string, subject, message string) error {
	body := "From: " + os.Getenv("MAIL_SENDER") + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", os.Getenv("MAIL_EMAIL"), os.Getenv("MAIL_PASSWORD"), os.Getenv("MAIL_SMTP_HOST"))
	smtpAddr := fmt.Sprintf("%s:%s", os.Getenv("MAIL_SMTP_HOST"), os.Getenv("MAIL_SMTP_PORT"))
	err := smtp.SendMail(smtpAddr, auth, os.Getenv("MAIL_EMAIL"), append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil

}
