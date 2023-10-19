package mailutil

import (
	"fmt"
	"net/smtp"
	"os"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/logger"
)

type SendFunc func(string, smtp.Auth, string, []string, []byte) error

func SendNotification(subject string, body string) error {
	if !env.Production {
		logger.Get().Info("notification sent to stderr")
		fmt.Fprintf(os.Stderr, buildBody(subject, body))
		return nil
	}

	defer logger.With("subject", subject).Info("notification sent")
	return SendNotificationFunc(smtp.SendMail, subject, body)
}

func buildBody(subject string, body string) string {
	return fmt.Sprintf(
		"From: %s\nTo: %s\nSubject: %s\n\n%s",
		env.NotificationFromHeader,
		env.NotificationToHeader,
		subject,
		body,
	)
}

// SendNotificationFunc sends an email using the supplied SendFunc.
func SendNotificationFunc(send SendFunc, subject string, body string) error {
	return send(
		env.SMTPHost+":"+env.SMTPPort,
		smtp.PlainAuth("", env.SMTPUser, env.SMTPPass, env.SMTPHost),
		env.NotificationFromAddress,
		[]string{env.NotificationToAddress},
		[]byte(buildBody(subject, body)),
	)
}
