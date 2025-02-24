package main

import (
	"fmt"
	"github.com/serhiirubets/rubeticket/config"
	"github.com/serhiirubets/rubeticket/internal/auth"
	"github.com/serhiirubets/rubeticket/internal/users"
	"github.com/serhiirubets/rubeticket/pkg/db"
	"github.com/serhiirubets/rubeticket/pkg/log"
	"github.com/serhiirubets/rubeticket/pkg/middleware"
	"net/http"
)

func App() http.Handler {
	conf := config.LoadConfig()
	dbInstance := db.NewDb(conf)
	logger := log.NewLogrusLogger(conf.LogLevel)
	router := http.NewServeMux()

	// Repositories
	usersRepository := users.NewUserRepository(dbInstance)

	// Services
	authService := auth.NewAuthService(usersRepository)

	// Handlers
	users.NewUserHandler(router, &users.UserHandlerDeps{
		UserRepository: usersRepository,
		Logger:         logger,
	})

	auth.NewAuthHandler(router, &auth.AuthHandlerDeps{
		Config:      conf,
		Logger:      logger,
		AuthService: authService,
	})

	// Middlewares
	middlewares := middleware.Chain(middleware.CORS)
	return middlewares(router)
}

func main() {
	router := App()

	server := http.Server{
		Addr:    ":7777",
		Handler: router,
	}

	fmt.Println("Server is listening on port 8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n ", err)
		return
	}
}
