package main

import (
	"github.com/RokibulHasan7/marketplace-prototype/internal/handlers"
	"github.com/RokibulHasan7/marketplace-prototype/internal/queue"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	database.ConnectDatabase()
	// Start Redis Queue Consumer in Background
	go queue.StartConsumer()

	r := chi.NewRouter()
	handlers.RegisterRoutes(r)

	log.Fatal(http.ListenAndServe(":3000", r))
}
