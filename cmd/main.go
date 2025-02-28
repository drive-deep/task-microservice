package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/drive-deep/task-microservice/routes"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	port := 8080
	fmt.Printf("🚀 Server running on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), router))
}
