package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config  config
	logger  *zap.SugaredLogger
	version string
	timeout int
}

// serving application
func (app *application) serve() error {
	// initialize http server
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	return srv.ListenAndServe()
}

func main() {
	// initialize zap logger
	loggerInit, _ := zap.NewProduction()
	defer loggerInit.Sync()
	logger := loggerInit.Sugar()

	// parse env variable `port`
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		logger.Fatal("failed to get port from env vars: ", zap.Error(err))
	}

	// parse env variable `timeout`
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		logger.Fatal("failed to get port from env vars: ", zap.Error(err))
	}

	// setup application config
	cfg := config{
		port: port,
		env:  os.Getenv("ENV"),
	}

	// initialize application
	app := &application{
		config:  cfg,
		logger:  logger,
		version: version,
		timeout: timeout,
	}

	// serve application
	if err := app.serve(); err != nil {
		app.logger.Fatal("unable to start the application: ", err)
	}
}
