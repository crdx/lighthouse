package mailutil

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"crdx.org/lighthouse/env"
	"crdx.org/lighthouse/logger"
)

type Func func(string, smtp.Auth, string, []string, []byte) error

type Config struct {
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

var packageConfig *Config

func Init(config *Config) {
	packageConfig = config
}

// Send sends an email if SMTP is enabled.
func Send(subject string, body string) error {
	if packageConfig == nil {
		panic("no mail configuration")
	}

	if packageConfig.Enabled == nil || !packageConfig.Enabled() {
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
		packageConfig.Host()+":"+packageConfig.Port(),
		smtp.PlainAuth("", packageConfig.User(), packageConfig.Pass(), packageConfig.Host()),
		packageConfig.FromAddress(),
		[]string{packageConfig.ToAddress()},
		[]byte(buildBody(subject, body)),
	)
}

func buildBody(subject string, body string) string {
	return fmt.Sprintf(
		"From: %s\nTo: %s\nSubject: %s\n\n%s",
		packageConfig.FromHeader(),
		packageConfig.ToHeader(),
		subject,
		body,
	)
}
