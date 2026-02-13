// main file - starting of the app code here

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Suthar345Piyush/invoicego/internal/config"
	"github.com/Suthar345Piyush/invoicego/internal/database"
	"github.com/Suthar345Piyush/invoicego/internal/handler"
	"github.com/Suthar345Piyush/invoicego/internal/middleware"
	"github.com/Suthar345Piyush/invoicego/internal/service"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// loading the configuration to start the app

func main() {

	//loading configuration

	cfg, err := config.Load()

	//if config failed to laod

	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// connecting to DB

	db, err := database.New(cfg.Database.ConnectionString())

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer db.Close()

	// initializing the auth and user service

	userService := service.NewUserService(db)
	authService := service.NewAuthService(userService, &cfg.JWT)

	// initializing the auth and user handlers

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// setting router using chi framework
	//NewRouter returns a mux object which implements router interface

	r := chi.NewRouter()

	// defining global middlewares

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logging)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	// health checks of the server

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// api routes for auth - login and register

	r.Route("/api/v1", func(r chi.Router) {

		// setting some public routes

		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		// protected routes

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(cfg.JWT.Secret))

			r.Get("/users/me", userHandler.GetMe)
		})

	})

	// at final starting the server and logging in terminal

	addr := fmt.Sprintf("%s", cfg.Server.Port)
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	fmt.Printf("Environment: %s\n", cfg.Server.Env)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
