package util

// A temporary home for all methods that don't yet have a better place to go.

import (
	"fmt"
	"net/smtp"
	"runtime"
	"strings"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/logger"
)

// PrintStackTrace prints out the current stack trace, stripping out stripDepth levels from the
// top.
func PrintStackTrace(stripDepth int) {
	b := make([]byte, 1024*10)
	length := runtime.Stack(b, false)
	s := string(b[:length])
	lines := strings.Split(s, "\n")
	header := lines[0]
	lines = lines[1+stripDepth*2:]
	fmt.Println(header)
	fmt.Println(strings.Join(lines, "\n"))
}

func Pluralise(count int, unit string) string {
	if count == 1 {
		return unit
	}
	return unit + "s"
}

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
		env.NotificationSender,
		[]string{env.NotificationRecipient},
		[]byte(fmt.Sprintf(
			"From: %s <%s>\nTo: %s\nSubject: %s\n\n%s",
			env.NotificationSenderName,
			env.NotificationSender,
			env.NotificationRecipient,
			subject,
			body,
		)),
	)
}
