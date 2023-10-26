package mailutil

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"crdx.org/lighthouse/logger"
)

type Func func(string, smtp.Auth, string, []string, []byte) error

type Config struct {
	SendToStdErr bool

	Enabled     func() bool
	Host        func() string
	Port        func() string
	User        func() string
	Pass        func() string
	FromAddress func() string
	ToAddress   func() string
	FromHeader  func() string
	ToHeader    func() string
}

var pkgConfig *Config

func Init(config *Config) {
	pkgConfig = config
}

// Send sends an email if SMTP is enabled.
func Send(subject string, body string) error {
	if pkgConfig == nil {
		return fmt.Errorf("no mail configuration")
	}

	if pkgConfig.Enabled == nil || !pkgConfig.Enabled() {
		return nil
	}

	if pkgConfig.SendToStdErr {
		logger.Get().Info("mail sent to stderr")
		fmt.Fprintln(os.Stderr, strings.TrimSpace(buildBody(subject, body)))
		return nil
	}

	defer logger.With("subject", subject).Info("mail sent to configured smtp server")
	return sendFunc(smtp.SendMail, subject, body)
}

func sendFunc(send Func, subject string, body string) error {
	return send(
		pkgConfig.Host()+":"+pkgConfig.Port(),
		smtp.PlainAuth("", pkgConfig.User(), pkgConfig.Pass(), pkgConfig.Host()),
		pkgConfig.FromAddress(),
		[]string{pkgConfig.ToAddress()},
		[]byte(buildBody(subject, body)),
	)
}

func buildBody(subject string, body string) string {
	return fmt.Sprintf(
		"From: %s\nTo: %s\nSubject: %s\n\n%s",
		pkgConfig.FromHeader(),
		pkgConfig.ToHeader(),
		subject,
		body,
	)
}
