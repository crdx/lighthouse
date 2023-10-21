package env

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	ModeDevelopment = "development"
	ModeProduction  = "production"

	LogTypeAll    = "all"
	LogTypeDisk   = "disk"
	LogTypeStderr = "stderr"
	LogTypeNone   = "none"

	AuthTypeBasic = "basic"
	AuthTypeNone  = "none"

	MailTypeSMTP = "smtp"
	MailTypeNone = "none"

	NotificationTypeMail = "mail"
	NotificationTypeNone = "none"
)

var (
	Debug = len(os.Getenv("LIGHTHOUSE_DEBUG")) > 0

	Production = os.Getenv("MODE") == ModeProduction
	BindHost   = os.Getenv("HOST")
	BindPort   = os.Getenv("PORT")

	LogType = os.Getenv("LOG_TYPE")
	LogPath = os.Getenv("LOG_PATH")

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

	NotificationType        = os.Getenv("NOTIFICATION_TYPE")
	NotificationFromHeader  = os.Getenv("NOTIFICATION_FROM_HEADER")
	NotificationFromAddress = os.Getenv("NOTIFICATION_FROM_ADDRESS")
	NotificationToHeader    = os.Getenv("NOTIFICATION_TO_HEADER")
	NotificationToAddress   = os.Getenv("NOTIFICATION_TO_ADDRESS")

	EnableLiveReload = os.Getenv("LIVE_RELOAD") != ""
)

func Check() {
	// In development no port means use a random port, but this will never be correct for production.
	if Production && BindPort == "" {
		panic("running in production (MODE=production) but no port set")
	}

	require("HOST")

	requireIn("AUTH_TYPE", []string{"none", "basic"}, true)

	if AuthType == AuthTypeBasic {
		require("AUTH_USER")
		require("AUTH_PASS")
	}

	if DatabaseSocket == "" && DatabaseHost == "" {
		panic("database socket (DB_SOCK) or host (DB_HOST) not set")
	}

	require("DB_NAME")
	require("DB_USER")

	requireIn("MAIL_TYPE", []string{"smtp", "none"}, true)

	if MailType == MailTypeSMTP {
		require("SMTP_HOST")
		require("SMTP_PORT")
		require("SMTP_USER")
		require("SMTP_PASS")
	}

	requireIn("LOG_TYPE", []string{"all", "disk", "stderr", "none"}, false)

	if LogType == LogTypeAll || LogType == LogTypeDisk {
		require("LOG_PATH")
	}

	requireIn("NOTIFICATION_TYPE", []string{"mail", "none"}, false)

	if NotificationType == NotificationTypeMail {
		require("NOTIFICATION_FROM_HEADER")
		require("NOTIFICATION_FROM_ADDRESS")
		require("NOTIFICATION_TO_HEADER")
		require("NOTIFICATION_TO_ADDRESS")
	}
}

func require(name string) {
	if os.Getenv(name) == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", name))
	}
}

func requireIn(name string, values []string, canBeEmpty bool) {
	if !canBeEmpty {
		require(name)
	}

	value := os.Getenv(name)

	if canBeEmpty && value == "" {
		return
	}

	if !slices.Contains(values, value) {
		s := ""
		if canBeEmpty {
			s = `, or the empty string ("")`
		}

		panic(fmt.Sprintf(
			`required environment variable %s contains an invalid value (must be one of: "%s"%s)`,
			name,
			strings.Join(values, `", "`),
			s,
		))
	}
}
