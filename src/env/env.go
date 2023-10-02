package env

import (
	"fmt"
	"os"
)

const (
	ModeDevelopment = "development"
	ModeProduction  = "production"

	LogTypeAll    = "all"
	LogTypeDisk   = "disk"
	LogTypeStdout = "stdout"
	LogTypeNone   = "none"

	AuthTypeBasic = "basic"
	AuthTypeNone  = "none"

	MailTypeSMTP = "smtp"
	MailTypeNone = "none"
)

var (
	Debug   = len(os.Getenv("LIGHTHOUSE_DEBUG")) > 0
	Verbose = len(os.Getenv("LIGHTHOUSE_VERBOSE")) > 0

	Production = os.Getenv("MODE") == ModeProduction
	BindHost   = os.Getenv("HOST")
	BindPort   = os.Getenv("PORT")

	LogType = os.Getenv("LOG_TYPE")

	AuthType = os.Getenv("AUTH_TYPE")
	AuthUser = os.Getenv("AUTH_USER")
	AuthPass = os.Getenv("AUTH_PASS")

	DatabaseName     = os.Getenv("DB_NAME")
	DatabaseUser     = os.Getenv("DB_USER")
	DatabasePass     = os.Getenv("DB_PASS")
	DatabaseSocket   = os.Getenv("DB_SOCK")
	DatabaseHost     = os.Getenv("DB_HOST")
	DatabaseCharSet  = os.Getenv("DB_CHARSET")
	DatabaseTimeZone = os.Getenv("DB_TZ")

	MailType = os.Getenv("MAIL_TYPE")
	SMTPHost = os.Getenv("SMTP_HOST")
	SMTPPort = os.Getenv("SMTP_PORT")
	SMTPUser = os.Getenv("SMTP_USER")
	SMTPPass = os.Getenv("SMTP_PASS")

	MailFrom = os.Getenv("MAIL_FROM")
	MailTo   = os.Getenv("MAIL_TO")

	MACVendorsAPIKey = os.Getenv("MACVENDORS_API_KEY")

	LocalTimeZone = os.Getenv("LOCAL_TZ")
)

func Check() {
	// In development no port means use a random port, but this will never be correct for production.
	if Production && BindPort == "" {
		panic("running in production (MODE=production) but no port set")
	}

	if AuthType != AuthTypeNone {
		require(AuthUser, "auth user", "AUTH_USER")
		require(AuthPass, "auth pass", "AUTH_PASS")
	}

	if DatabaseSocket == "" && DatabaseHost == "" {
		panic("database socket (DB_SOCK) or host (DB_HOST) not set")
	}

	require(DatabaseName, "database name", "DB_NAME")
	require(DatabaseUser, "database user", "DB_USER")

	if MailType == MailTypeSMTP {
		require(SMTPHost, "SMTP host", "SMTP_HOST")
		require(SMTPPort, "SMTP port", "SMTP_PORT")
		require(SMTPUser, "SMTP user", "SMTP_USER")
		require(SMTPPass, "SMTP pass", "SMTP_PASS")
		require(MailFrom, "mail from", "MAIL_FROM")
		require(MailTo, "mail to", "MAIL_TO")
	}

	require(LogType, "log type", "LOG_TYPE")
	require(LocalTimeZone, "local timezone", "LOCAL_TZ")
}

func require(v, name, envvar string) {
	if v == "" {
		panic(fmt.Sprintf("%s (%s) not set", name, envvar))
	}
}
