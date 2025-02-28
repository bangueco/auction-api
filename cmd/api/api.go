package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bangueco/auction-api/internal/config"
	"github.com/bangueco/auction-api/internal/handlers"
	"github.com/bangueco/auction-api/internal/handlers/helper"
	"github.com/bangueco/auction-api/internal/lib"
	"github.com/bangueco/auction-api/internal/middleware"
	"github.com/bangueco/auction-api/internal/repositories"
	"github.com/bangueco/auction-api/internal/services"
	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	// This loads the .env file in the root of the project allowing us to access them via os.Getenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load config
	// This loads the config from the .env file, without godotenv.Load() this will return empty values
	cfg := config.Load()

	// Initialize database connection
	dbpool := lib.InitDBConnection()

	// Initialize dependencies (handlers, services, repositories)
	itemRepository := repositories.NewItemRepository(dbpool)
	itemService := services.NewItemService(itemRepository)
	itemHandler := handlers.NewItemHandler(itemService)

	userRepository := repositories.NewUserRepository(dbpool)
	userService := services.NewUserService(userRepository)

	authHandler := handlers.NewAuthHandler(userService)

	// Initialize router
	r := chi.NewRouter()
	r.Use(chimiddle.Logger)
	r.Use(chimiddle.Recoverer)
	r.Use(chimiddle.StripSlashes)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		helper.WriteResponseMessage(w, "404 Route Not found", http.StatusNotFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		helper.WriteResponseMessage(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/api/items", func(r chi.Router) {
		r.Use(middleware.AuthGuard)
		r.Get("/", itemHandler.GetItems)
		r.Get("/{id}", itemHandler.GetItemByID)
		r.Post("/", itemHandler.CreateItem)
		r.Put("/{id}", itemHandler.UpdateItem)
		r.Delete("/{id}", itemHandler.DeleteItem)
	})

	// Start server
	log.Printf("Server started on port %s", cfg.PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.PORT), r))
}
