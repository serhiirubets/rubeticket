package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serhiirubets/rubeticket/app"
	"github.com/serhiirubets/rubeticket/config"
	_ "github.com/serhiirubets/rubeticket/docs"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
)

// @title Concert booking API
// @version 1.0
// @description This is a Concert booking application
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /v1
// @host localhost:777
func main() {
	conf := config.LoadConfig()
	logger := log.NewLogrusLogger(conf.LogLevel)
	router, initErr := app.InitApp(conf, logger)

	if initErr != nil {
		logger.Error("Server error: %v\n ", initErr)
		os.Exit(1)
	}

	port := conf.App.Port

	app := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Server is listening on port ", port)
		if err := app.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start", "error", err.Error())
			os.Exit(1)
		}
	}()

	<-stop
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "error", err.Error())
		os.Exit(1)
	}

	logger.Info("Server stopped gracefully")
}
