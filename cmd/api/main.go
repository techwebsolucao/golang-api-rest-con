package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	r.Get("/api/v1/users", userController.GetAll)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		log.Println("Servidor rodando na porta :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Servidor encerrado")
}
