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
	swagger "github.com/swaggo/http-swagger/v2"

	"github.com/user/golang-api-rest/internal/config"
	"github.com/user/golang-api-rest/internal/controllers"
	"github.com/user/golang-api-rest/internal/database"
	"github.com/user/golang-api-rest/internal/middleware"
	"github.com/user/golang-api-rest/internal/repositories"
	"github.com/user/golang-api-rest/internal/services"

	_ "github.com/user/golang-api-rest/docs"
)

// @title Go REST API
// @version 1.0
// @description API de prova de conceito com autenticação JWT, MySQL e envio de email.
// @termsOfService http://swagger.io/terms/

// @contact.name Suporte
// @contact.email suporte@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Insira o token JWT no formato: Bearer {token}
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Erro ao executar migração: %v", err)
	}

	userRepo := repositories.NewMySQLRepository(db)
	emailService := services.NewEmailService(cfg)
	authService := services.NewAuthService(userRepo, cfg, emailService)
	userService := services.NewUserService(userRepo)
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	r.Get("/swagger/*", swagger.Handler(swagger.URL("http://localhost:8080/swagger/doc.json")))

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
		r.Post("/refresh", authController.Refresh)
		r.Get("/verify-email", authController.Verify)
		r.Post("/verify-email", authController.Verify)
		r.Post("/verify-token", authController.VerifyJWT)
	})

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Use(middleware.RequireAuth(cfg))

		r.Get("/", userController.GetAll)
		r.Get("/{id}", userController.GetByID)
		r.Put("/{id}", userController.Update)

		r.With(middleware.RequireRole("admin")).Delete("/{id}", userController.Delete)
	})

	//Fallbacks
	//404
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"endpoint não encontrado"}`))
	})
	//500
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"método não permitido"}`))
	})
	//503
	r.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		log.Printf("Servidor rodando na porta :%s", cfg.AppPort)
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
