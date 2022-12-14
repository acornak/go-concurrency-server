package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
)

const version = "1.0.0"

var logger *zap.SugaredLogger

type config struct {
	port int
	env  string
}

type application struct {
	config  config
	logger  *zap.SugaredLogger
	version string
}

func init() {
	logger = zap.NewExample().Sugar()
}

// serving application
func (app *application) serve(routes http.Handler) error {
	// initialize http server
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           routes,
		IdleTimeout:       60 * time.Second,
		ReadTimeout:       60 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      20 * time.Second,
	}

	// start serving
	return srv.ListenAndServe()
}

func main() {
	// initialize zap logger
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Fatal("failed to initialize zap logger: ", err)
		}
	}()

	// parse env variable `port`
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		logger.Error("failed to get port from env vars: ", zap.Error(err))
		return
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
	}

	// serve application
	app.logger.Info("starting server in ", app.config.env, " mode on port ", app.config.port)

	// run serving in another thread (for testing purposes)
	quit := make(chan os.Signal, 1)
	go func() {
		if err := app.serve(app.routes()); err != nil {
			app.logger.Error("unable to start the application: ", zap.Error(err))
			return
		}
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	<-quit
}
