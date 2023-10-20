package mailutil

import (
	"fmt"
	"net/smtp"
	"os"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/logger"
)

type Func func(string, smtp.Auth, string, []string, []byte) error

// Send sends an email if SMTP is enabled.
func Send(subject string, body string) error {
	if env.MailType != env.MailTypeSMTP {
		return nil
	}

	if !env.Production {
		logger.Get().Info("mail sent to stderr")
		fmt.Fprintln(os.Stderr, strings.TrimSpace(buildBody(subject, body)))
		return nil
	}

	defer logger.With("subject", subject).Info("mail sent to configured mailserver")
	return sendFunc(smtp.SendMail, subject, body)
}

func sendFunc(send Func, subject string, body string) error {
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
