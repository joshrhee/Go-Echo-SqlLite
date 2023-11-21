package main

import (
	"context"
	"errors"
	"fmt"
	"identity-coding-test/client"
	"identity-coding-test/internal/config"
	"identity-coding-test/internal/database"
	"identity-coding-test/internal/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	clocklib "github.com/benbjohnson/clock"
	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := run(); err != nil {
		logrus.WithError(err).Panic()
	}
}

func run() error {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)

	setting := config.NewSetting()

	clock := clocklib.New()

	// Init DB
	db, err := client.NewSQLite3(&client.SQLite3Config{
		DataSourceName:  setting.SQLite3DataSourceName,
		MaxIdleConn:     setting.SQLite3MaxIdleConn,
		MaxOpenConn:     setting.SQLite3MaxOpenConn,
		MaxConnLifetime: setting.SQLite3MaxConnLifetime,
		MaxConnIdleTime: setting.SQLite3MaxConnIdleTime,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			logrus.WithError(err).Error("close db")
		}
	}()

	// Init Tables
	if err := database.CreateTables(db); err != nil {
		return err
	}

	// Init Config
	cfg, err := config.New(setting, clock, db)
	if err != nil {
		return err
	}

	// Init EchoServer
	echoServer, err := server.NewEchoServer(cfg)
	if err != nil {
		return err
	}

	go func() {
		logrus.Infof("starting echo server, localhost:%d", setting.HTTPPort)
		if err := echoServer.Start(fmt.Sprintf(":%d", setting.HTTPPort)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.WithError(err).Panic("shutting down echo server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(setting.GracefulShutdownTimeoutSec)*time.Second)
	defer cancel()

	logrus.Info("closing echo server")
	if err := echoServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
