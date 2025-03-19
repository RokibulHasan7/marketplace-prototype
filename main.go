package main

import (
	"github.com/RokibulHasan7/marketplace-prototype/internal/billing"
	"github.com/RokibulHasan7/marketplace-prototype/internal/handlers"
	"github.com/RokibulHasan7/marketplace-prototype/internal/queue"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	database.ConnectDatabase()

	// Start Redis Queue Consumers in Background
	go queue.StartCreateConsumer()
	go queue.StartDeleteConsumer()

	// Start billing background job which will update billing data on hourly basis
	go billing.StartBillingUpdater()

	r := chi.NewRouter()
	handlers.RegisterRoutes(r)

	log.Fatal(http.ListenAndServe(":3000", r))
}
