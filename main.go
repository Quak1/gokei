package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Quak1/gokei/internal/auth"
	"github.com/Quak1/gokei/internal/config"
	"github.com/Quak1/gokei/internal/database"
	"github.com/Quak1/gokei/internal/handlers"
	"github.com/Quak1/gokei/internal/services"
)

func main() {
	cfg := config.Load()

	db := database.NewDBConnection(cfg.Database.Url)
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	jwtService := auth.NewJWTService([]byte(cfg.Server.TokenSecret))
	userService := services.NewUserService(db.Queries, jwtService)
	userHandler := handlers.NewUserHandler(*userService)

	mux.HandleFunc("POST /api/register", userHandler.Register)
	mux.HandleFunc("POST /api/login", userHandler.TokenLogin)
	mux.Handle("GET /api/echo", jwtService.AuthMiddleware(http.HandlerFunc(userHandler.EchoUsername)))

	server := http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}
	log.Printf("Serving on port: %s\n", cfg.Server.Port)
	log.Fatal(server.ListenAndServe())
}
