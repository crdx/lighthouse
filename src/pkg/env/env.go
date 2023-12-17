package env

import (
	"errors"
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
)

var env map[string]string

var (
	mode = func() string { return env["MODE"] }

	Debug      = func() bool { return truthy(env["LIGHTHOUSE_DEBUG"]) }
	Production = func() bool { return env["MODE"] == ModeProduction }

	Host = func() string { return env["HOST"] }
	Port = func() string { return env["PORT"] }

	LogType = func() string { return env["LOG_TYPE"] }
	LogPath = func() string { return env["LOG_PATH"] }

	DatabaseName     = func() string { return env["DB_NAME"] }
	DatabaseUser     = func() string { return env["DB_USER"] }
	DatabasePass     = func() string { return env["DB_PASS"] }
	DatabaseSocket   = func() string { return env["DB_SOCK"] }
	DatabaseHost     = func() string { return env["DB_HOST"] }
	DatabaseCharset  = func() string { return env["DB_CHARSET"] }
	DatabaseTimezone = func() string { return env["DB_TZ"] }

	DefaultRootPass = func() string { return or("DEFAULT_ROOT_PASS", "root") }
	DefaultAnonPass = func() string { return or("DEFAULT_ANON_PASS", "anon") }

	LiveReload      = func() bool { return truthy(env["LIVE_RELOAD"]) }
	DisableAuth     = func() bool { return truthy(env["DISABLE_AUTH"]) }
	DisableServices = func() bool { return truthy(env["DISABLE_SERVICES"]) }

	TrustedProxies = func() string { return env["TRUSTED_PROXIES"] }
)

func Init() {
	if env == nil {
		env = map[string]string{}
	}

	for _, v := range os.Environ() {
		name, value, ok := strings.Cut(v, "=")
		if !ok {
			continue
		}

		env[name] = value
	}
}

func InitFrom(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	env = parse(string(b))
	Init()

	return nil
}

func Validate() error {
	var errs []error
	e := func(err error) {
		errs = append(errs, err)
	}

	e(require(mode, "MODE"))

	if Production() && Port() == "" {
		// In development no port means use a random port, but this will never be correct for production.
		return fmt.Errorf("running in production but no port set")
	}

	e(require(Host, "HOST"))

	if DatabaseSocket() == "" && DatabaseHost() == "" {
		e(fmt.Errorf("DB_SOCK or DB_HOST required"))
	}

	e(require(DatabaseName, "DB_NAME"))
	e(require(DatabaseUser, "DB_USER"))
	e(require(DatabaseCharset, "DB_CHARSET"))
	e(require(DatabaseTimezone, "DB_TZ"))

	e(requireIn(LogType, "LOG_TYPE", []string{"all", "disk", "stderr", "none"}, false))

	if LogType() == LogTypeAll || LogType() == LogTypeDisk {
		e(require(LogPath, "LOG_PATH"))
	}

	return errors.Join(errs...)
}

func parse(s string) map[string]string {
	m := map[string]string{}

	for _, line := range strings.Split(s, "\n") {
		line := strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		name, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		if len(value) > 0 {
			n := len(value) - 1
			if value[0] == '"' && value[n] == '"' {
				value = value[1:n]
			}
		}

		m[name] = value
	}

	return m
}

func require(f func() string, name string) error {
	if f() == "" {
		return fmt.Errorf("%s required", name)
	}
	return nil
}

func requireIn(f func() string, name string, values []string, canBeEmpty bool) error {
	if !canBeEmpty {
		if err := require(f, name); err != nil {
			return err
		}
	}

	value := f()

	if canBeEmpty && value == "" {
		return nil
	}

	if !slices.Contains(values, value) {
		s := ""
		if canBeEmpty {
			s = `, or the empty string ("")`
		}

		return fmt.Errorf(
			`%s contains an invalid value (must be one of: "%s"%s)`,
			name,
			strings.Join(values, `", "`),
			s,
		)
	}

	return nil
}

func or(name string, default_ string) string {
	value := env[name]

	if value != "" {
		return value
	} else {
		return default_
	}
}

func truthy(s string) bool {
	return slices.Contains([]string{"true", "1", "yes"}, s)
}
