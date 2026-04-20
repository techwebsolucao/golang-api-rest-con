package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/user/golang-api-rest/internal/controllers"
	"github.com/user/golang-api-rest/internal/middleware"
	"github.com/user/golang-api-rest/internal/repositories"
	"github.com/user/golang-api-rest/internal/services"
)

func main() {
	userRepo := repositories.NewMemoryRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	r.Get("/users", userController.GetAll)

	log.Println("🚀 Servidor rodando na porta :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
