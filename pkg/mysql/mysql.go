package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"crdx.org/lighthouse/pkg/logger"
	"github.com/go-sql-driver/mysql"
)

const (
	maximumConnectionLifetime = 3 * time.Minute
	maximumConnectionIdleTime = 3 * time.Minute
	maximumOpenConnections    = 5
	maximumIdleConnections    = 2
)

var (
	poolsMutex sync.Mutex
	pools      []*sql.DB

	driverLoggerOnce sync.Once
)

type driverLogger struct{}

func (self driverLogger) Print(values ...any) {
	logger.Get().Warn("mysql driver", "message", strings.TrimSpace(fmt.Sprint(values...)))
}

func Open(dataSourceName string) (*sql.DB, error) {
	driverLoggerOnce.Do(func() {
		_ = mysql.SetLogger(driverLogger{})
	})

	pool, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	pool.SetMaxOpenConns(maximumOpenConnections)
	pool.SetMaxIdleConns(maximumIdleConnections)
	pool.SetConnMaxLifetime(maximumConnectionLifetime)
	pool.SetConnMaxIdleTime(maximumConnectionIdleTime)

	poolsMutex.Lock()
	defer poolsMutex.Unlock()
	pools = append(pools, pool)

	return pool, nil
}

func CloseAll() {
	poolsMutex.Lock()
	defer poolsMutex.Unlock()

	for _, pool := range pools {
		_ = pool.Close()
	}
	pools = nil
}
