package main

import (
	"fmt"
	"github.com/serhiirubets/rubeticket/config"
	_ "github.com/serhiirubets/rubeticket/docs"
	"github.com/serhiirubets/rubeticket/internal/accounts"
	"github.com/serhiirubets/rubeticket/internal/auth"
	"github.com/serhiirubets/rubeticket/internal/users"
	"github.com/serhiirubets/rubeticket/pkg/db"
	"github.com/serhiirubets/rubeticket/pkg/log"
	"github.com/serhiirubets/rubeticket/pkg/middleware"
	"github.com/swaggo/http-swagger"
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

	accounts.NewAccountHandler(router, &accounts.AccountHandlerDeps{
		Logger:         logger,
		UserRepository: usersRepository,
		Config:         conf,
	})

	// Middlewares
	middlewares := middleware.Chain(middleware.CORS)

	router.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:7777/swagger/doc.json"),
	))

	return middlewares(router)
}

// @title Concert booking API
// @version 1.0
// @description This is a Concert booking application
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /v1
// @host localhost:777
func main() {
	router := App()

	server := http.Server{
		Addr:    ":7777",
		Handler: router,
	}

	fmt.Println("Server is listening on port 7777")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n ", err)
		return
	}
}
