package main

import (
	"flag"
	"log"
	"test/config"
	"test/internal/http"
	router "test/internal/http/route"
	initialization "test/internal/init"
	rabbitmq "test/internal/rabbitMq"
	pgStorage "test/internal/repository/pgStorage"
	"test/internal/repository/redis"
	authservice "test/internal/usecases/authService"
	"test/internal/usecases/taskService"

	_ "test/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
)

// @title My API
// @version 1.0
// @description This is a sample server.

// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Printf("Warning: error loading .env: %v (using defaults)", err)
		cfg, _ = config.Load("")
	}

	addr := flag.String("addr", cfg.AppAddr, "address for http server")
	flag.Parse()

	objectRepo, err := pgStorage.NewPostgresStorsge(cfg.PostgresConnStr())
	if err != nil {
		log.Fatalf("failed creating postgres connection %v", err)
	}

	redisStorage, err := redis.NewRadisConnectiron()
	if err != nil {
		log.Fatalf("failed creating redis connection %v", err)
	}

	objectServiceAuth := authservice.NewObject(objectRepo, redisStorage)
	objectSender, err := rabbitmq.NewRabbitMQConnection(cfg.RabbitMQURL(), "que")
	if err != nil {
		log.Fatalf("failed creating RabbitMq connection %v", err)
	}
	objectService := taskService.NewObject(objectRepo, objectSender)
	objectHandlers := http.NewHandler(objectService, objectServiceAuth)

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	router.WithObjectHandlers(r, objectHandlers)

	log.Printf("Starting server on %s, my homie", *addr)
	if err := initialization.CreateAndRunServer(r, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
