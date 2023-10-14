package mailutil

import (
	"fmt"
	"net/smtp"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/logger"
)

type SendFunc func(string, smtp.Auth, string, []string, []byte) error

func SendNotification(subject string, body string) error {
	if !env.Production {
		logger.With("subject", subject, "body", body).Info("notification NOT sent due to non-production environment")
		return nil
	}

	defer logger.With("subject", subject).Info("notification sent")
	return SendNotificationFunc(smtp.SendMail, subject, body)
}

// SendNotificationFunc sends an email using the supplied SendFunc.
func SendNotificationFunc(send SendFunc, subject string, body string) error {
	return send(
		env.SMTPHost+":"+env.SMTPPort,
		smtp.PlainAuth("", env.SMTPUser, env.SMTPPass, env.SMTPHost),
		env.NotificationFromAddress,
		[]string{env.NotificationToAddress},
		[]byte(fmt.Sprintf(
			"From: %s\nTo: %s\nSubject: %s\n\n%s",
			env.NotificationFromHeader,
			env.NotificationToAddress,
			subject,
			body,
		)),
	)
}
