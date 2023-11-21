package config

import (
	"database/sql"

	clocklib "github.com/benbjohnson/clock"
)

type Config interface {
	Setting() Setting
	Clock() clocklib.Clock
	ReadWriteDB() *sql.DB
}

type DefaultConfig struct {
	setting     Setting
	clock       clocklib.Clock
	readWriteDB *sql.DB
}

func New(
	setting Setting,
	clock clocklib.Clock,
	readWriteDB *sql.DB,
) (Config, error) {
	return &DefaultConfig{
		setting:     setting,
		clock:       clock,
		readWriteDB: readWriteDB,
	}, nil
}

func (c *DefaultConfig) Setting() Setting {
	return c.setting
}

func (c *DefaultConfig) Clock() clocklib.Clock {
	return c.clock
}

func (c *DefaultConfig) ReadWriteDB() *sql.DB {
	return c.readWriteDB
}
