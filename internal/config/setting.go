package config

import (
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Setting struct {
	HTTPPort                   int
	GracefulShutdownTimeoutSec int

	SQLite3DataSourceName  string
	SQLite3MaxOpenConn     int
	SQLite3MaxIdleConn     int
	SQLite3MaxConnLifetime time.Duration
	SQLite3MaxConnIdleTime time.Duration
}

func getEnv(key, defaultValue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		if defaultValue != "" {
			value = defaultValue
		} else {
			logrus.Fatal("missing required environment variable: " + key)
		}
	}
	return
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func mustParseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		panic(err)
	}
	return b
}

func mustParseFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

func NewSetting() Setting {
	return Setting{
		HTTPPort:                   mustAtoi(getEnv("HTTP_PORT", "3000")),
		GracefulShutdownTimeoutSec: mustAtoi(getEnv("GRACEFUL_SHUTDOWN_TIMEOUT_SEC", "10")),

		SQLite3DataSourceName:  getEnv("SQLITE3_DATA_SOURCE_NAME", "./identity.db"),
		SQLite3MaxOpenConn:     mustAtoi(getEnv("SQLITE3_MAX_OPEN_CONN", "20")),
		SQLite3MaxIdleConn:     mustAtoi(getEnv("SQLITE3_MAX_IDLE_CONN", "10")),
		SQLite3MaxConnLifetime: mustParseDuration(getEnv("SQLITE3_MAX_CONN_LIFETIME", "600s")),
		SQLite3MaxConnIdleTime: mustParseDuration(getEnv("SQLITE3L_MAX_CONN_IDLE_TIME", "5s")),
	}
}
