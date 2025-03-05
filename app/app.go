package app

import (
	"github.com/serhiirubets/rubeticket/config"
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
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

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

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Router))

	return middlewares(router), nil
}
