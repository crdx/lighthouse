package mysql

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
)

func Open(dataSourceName string) (*sql.DB, error) {
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
