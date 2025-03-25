package app

import (
	"net/http"

	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/app/accounts"
	"github.com/serhiirubets/rubeticket/internal/app/admin/bands"
	"github.com/serhiirubets/rubeticket/internal/app/admin/concerts"
	"github.com/serhiirubets/rubeticket/internal/app/admin/venues"
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
)

func InitApp(conf *config.Config, logger log.ILogger) (http.Handler, error) {
	router := http.NewServeMux()
	v1Router := http.NewServeMux()
	v1AdminRouter := http.NewServeMux()
	dbInstance := db.NewDb(conf)
	storage := filestorage.NewLocalStorage("uploads")

	// Middlewares
	middlewares := middleware.Chain(middleware.CORS)

	// Open routes
	openRoutes := []string{
		"/auth/login",
		"/auth/register",
		"/uploads/{fileName}",
	}
	authMiddleware := middleware.NewAuthMiddleware(conf, logger, openRoutes, "/api/v1")
	authMiddlewareAdmin := middleware.NewAuthMiddleware(conf, logger, nil, "/admin/v1")

	// Repositories
	usersRepository := users.NewUserRepository(dbInstance)
	fileRepository := file.NewRepository(dbInstance)
	venueRepository := venues.NewVenueRepository(dbInstance)
	bandRepository := bands.NewBandRepository(dbInstance)
	concertRepository := concerts.NewConcertRepository(dbInstance)

	// Services
	authService := auth.NewAuthService(usersRepository)
	venueService := venues.NewVenueService(venueRepository)
	bandService := bands.NewBandService(bandRepository)
	concertService := concerts.NewConcertService(concertRepository, venueRepository, bandRepository)

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

	// Public handlers
	auth.NewAuthHandler(v1Router, &auth.AuthHandlerDeps{
		Config:      conf,
		Logger:      logger,
		AuthService: authService,
	})

	// Private handlers
	users.NewUserHandler(v1Router, &users.UserHandlerDeps{
		UserRepository: usersRepository,
		Logger:         logger,
	})

	accounts.NewAccountHandler(v1Router, &accounts.AccountHandlerDeps{
		Logger:         logger,
		UserRepository: usersRepository,
		Config:         conf,
		FileUploader:   fileUploader,
	})

	uploads.NewUploadsHandler(v1Router, &uploads.HandlerDeps{
		Logger:       logger,
		Config:       conf,
		FileUploader: fileUploader,
	})

	// Admin handlers
	venues.NewVenueHandler(v1AdminRouter, &venues.VenueHandlerDeps{
		Config:         conf,
		Logger:         logger,
		Service:        venueService,
		UserRepository: usersRepository,
	})

	bands.NewBandHandler(v1AdminRouter, &bands.BandHandlerDeps{
		Config:         conf,
		Logger:         logger,
		Service:        bandService,
		UserRepository: usersRepository,
	})

	concerts.NewConcertHandler(v1AdminRouter, &concerts.ConcertHandlerDeps{
		Config:         conf,
		Logger:         logger,
		Service:        concertService,
		UserRepository: usersRepository,
	})

	// Apply middleware
	v1RouterWithAuth := authMiddleware.Auth(v1Router)
	adminRouterWithAuth := authMiddlewareAdmin.Auth(v1AdminRouter)
	adminRouterWithAuthAndAdmin := authMiddlewareAdmin.AdminOnly(adminRouterWithAuth)

	// Swagger
	router.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:7777/swagger/doc.json"),
	))

	// Setup routes
	router.Handle("/admin/v1/", http.StripPrefix("/admin/v1", adminRouterWithAuthAndAdmin))
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", v1RouterWithAuth))
	router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	return middlewares(router), nil
}
