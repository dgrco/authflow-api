package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrco/autoflow/internal/config"
	"github.com/dgrco/autoflow/internal/database"
	"github.com/dgrco/autoflow/internal/handler"
	"github.com/dgrco/autoflow/internal/infra/repo"
	"github.com/dgrco/autoflow/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load environment
	cfg := config.Load()

	// Connect to database
	pool, err := database.Connect(cfg.DatabaseUrl);
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Database connected successfully")

	pgRepo := repo.NewPgRepository(pool)

	// Auth service/handler 
	authService := service.NewAuthService(pgRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService, cfg.IsSecureMode())
 
	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Handlers
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/logout", authHandler.Logout)
	})

	// Listen
	log.Printf("Server started on port %s", cfg.ApiPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.ApiPort), r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
