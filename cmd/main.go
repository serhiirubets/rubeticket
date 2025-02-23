package main

import (
	"fmt"
	"net/http"
)

func App() *http.ServeMux {
	//conf := configs.LoadConfig()
	//database := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories

	return router
}

func main() {
	router := App()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server is listening on port 8080")
	server.ListenAndServe()
}
