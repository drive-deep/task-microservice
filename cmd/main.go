package main

import (
	"log"
	"net/http"

	"github.com/drive-deep/task-microservice/cache"
	"github.com/drive-deep/task-microservice/config"
	"github.com/drive-deep/task-microservice/database"
	"github.com/drive-deep/task-microservice/message_queue"
	"github.com/drive-deep/task-microservice/repositories"
	"github.com/drive-deep/task-microservice/routes"
	"github.com/drive-deep/task-microservice/services"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	postgres, err := database.NewPostgresDB().Connect()
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	redis, err := cache.NewRedisCache(20).Connect()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	repo := repositories.NewTaskRepository(postgres)

	// Initialize the Kafka message queue
	kafkaMessageQueue, err := message_queue.NewKafkaMessageQueue(services.NewTaskService(repo, redis)).Connect([]string{cfg.Kafka.Broker}, cfg.Kafka.GroupID)
	if err != nil {
		log.Fatalf("failed to connect to Kafka: %v", err)
	}
	go kafkaMessageQueue.StartConsuming(cfg.Kafka.Topics)
	defer kafkaMessageQueue.Close()

	services := services.NewTaskService(repo, redis)

	mux := mux.NewRouter()
	routes.RegisterRoutes(mux, *services)

	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", mux)

}
