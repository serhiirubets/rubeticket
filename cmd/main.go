package main

import (
	"context"
	"github.com/serhiirubets/rubeticket/config"
	_ "github.com/serhiirubets/rubeticket/docs"
	"github.com/serhiirubets/rubeticket/internal/app/accounts"
	"github.com/serhiirubets/rubeticket/internal/app/auth"
	"github.com/serhiirubets/rubeticket/internal/app/file"
	"github.com/serhiirubets/rubeticket/internal/app/fileuploader"
	"github.com/serhiirubets/rubeticket/internal/app/uploads"
	"github.com/serhiirubets/rubeticket/internal/app/users"
	"github.com/serhiirubets/rubeticket/internal/pkg/db"
	"github.com/serhiirubets/rubeticket/internal/pkg/filestorage"
	"github.com/serhiirubets/rubeticket/internal/pkg/log"
	"github.com/serhiirubets/rubeticket/internal/pkg/middleware"
	"github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	Router http.Handler
	Logger log.ILogger
	Server *http.Server
}

func InitApp(conf *config.Config, logger log.ILogger) (http.Handler, error) {
	router := http.NewServeMux()
	v1Router := http.NewServeMux()
	dbInstance := db.NewDb(conf)
	storage := filestorage.NewLocalStorage("uploads")

	// Middlewares
	middlewares := middleware.Chain(middleware.CORS)

	// Repositories
	usersRepository := users.NewUserRepository(dbInstance)
	fileRepository := file.NewRepository(dbInstance)

	// Services
	authService := auth.NewAuthService(usersRepository)

	// Utils
	fileUploader := fileuploader.NewFileUploader(&fileuploader.Deps{
		Logger:         logger,
		DB:             dbInstance,
		MaxSizeMB:      10,
		AllowedTypes:   []string{"image/"},
		Storage:        storage,
		FileRepository: fileRepository,
	})

	// Handlers
	users.NewUserHandler(v1Router, &users.UserHandlerDeps{
		UserRepository: usersRepository,
		Logger:         logger,
	})

	auth.NewAuthHandler(v1Router, &auth.AuthHandlerDeps{
		Config:      conf,
		Logger:      logger,
		AuthService: authService,
	})

	accounts.NewAccountHandler(v1Router, &accounts.AccountHandlerDeps{
		Logger:         logger,
		UserRepository: usersRepository,
		Config:         conf,
		FileUploader:   fileUploader,
	})

	uploads.NewUploadsHandler(router, &uploads.HandlerDeps{
		Logger:       logger,
		Config:       conf,
		FileUploader: fileUploader,
	})

	router.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:7777/swagger/doc.json"),
	))

	router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	router.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	return middlewares(router), nil
}

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
	router, initErr := InitApp(conf, logger)

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
		logger.Info("Server is listening on port", "port", port)
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
