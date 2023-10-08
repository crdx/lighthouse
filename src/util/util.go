package util

// A temporary home for all methods that don't yet have a better place to go.

import (
	"fmt"
	"net/smtp"
	"runtime"
	"strings"

	"crdx.org/lighthouse/env"
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

func SendMail(subject string, message string) error {
	return SendMailFunc(smtp.SendMail, subject, message)
}

// SendMail sends an email via the supplied MailSender.
func SendMailFunc(send SendFunc, subject string, message string) error {
	return send(
		env.SMTPHost+":"+env.SMTPPort,
		smtp.PlainAuth("", env.SMTPUser, env.SMTPPass, env.SMTPHost),
		env.MailFrom,
		[]string{env.MailTo},
		[]byte(fmt.Sprintf(
			"To: %s\nSubject: %s\n\n%s",
			env.MailTo,
			subject,
			message,
		)),
	)
}
