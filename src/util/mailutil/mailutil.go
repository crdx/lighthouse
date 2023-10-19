package mailutil

import (
	"fmt"
	"net/smtp"
	"os"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/logger"
)

type Func func(string, smtp.Auth, string, []string, []byte) error

// Send sends an email.
func Send(subject string, body string) error {
	if !env.Production {
		logger.Get().Info("mail sent to stderr")
		fmt.Fprintf(os.Stderr, buildBody(subject, body))
		return nil
	}

	defer logger.With("subject", subject).Info("mail sent to configured mailserver")
	return SendFunc(smtp.SendMail, subject, body)
}

// SendFunc sends an email using the supplied Func.
func SendFunc(send Func, subject string, body string) error {
	return send(
		env.SMTPHost+":"+env.SMTPPort,
		smtp.PlainAuth("", env.SMTPUser, env.SMTPPass, env.SMTPHost),
		env.NotificationFromAddress,
		[]string{env.NotificationToAddress},
		[]byte(buildBody(subject, body)),
	)
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
