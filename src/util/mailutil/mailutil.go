package mailutil

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/logger"
	"crdx.org/lighthouse/setting"
)

type Func func(string, smtp.Auth, string, []string, []byte) error

// Send sends an email if SMTP is enabled.
func Send(subject string, body string) error {
	if !setting.GetBool(setting.EnableSMTP) {
		return nil
	}

	if !env.Production {
		logger.Get().Info("mail sent to stderr")
		fmt.Fprintln(os.Stderr, strings.TrimSpace(buildBody(subject, body)))
		return nil
	}

	defer logger.With("subject", subject).Info("mail sent to configured smtp server")
	return sendFunc(smtp.SendMail, subject, body)
}

func sendFunc(send Func, subject string, body string) error {
	return send(
		setting.Get(setting.SMTPHost)+":"+setting.Get(setting.SMTPPort),
		smtp.PlainAuth("", setting.Get(setting.SMTPUser), setting.Get(setting.SMTPPass), setting.Get(setting.SMTPHost)),
		setting.Get(setting.NotificationFromAddress),
		[]string{setting.Get(setting.NotificationToAddress)},
		[]byte(buildBody(subject, body)),
	)
}

func buildBody(subject string, body string) string {
	return fmt.Sprintf(
		"From: %s\nTo: %s\nSubject: %s\n\n%s",
		setting.Get(setting.NotificationFromHeader),
		setting.Get(setting.NotificationToHeader),
		subject,
		body,
	)
}
