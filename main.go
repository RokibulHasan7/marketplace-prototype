package main

import (
	"github.com/RokibulHasan7/marketplace-prototype/internal/apis"
	"github.com/RokibulHasan7/marketplace-prototype/internal/queue"
	"github.com/RokibulHasan7/marketplace-prototype/internal/services/billing"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	database.ConnectDatabase()

	// Start Redis Queue Consumers in Background
	go queue.StartInstallerConsumer()
	go queue.StartUninstallerConsumer()

	// Start billing background job which will update billing data on hourly basis
	go billing.StartBillingUpdater()

	r := chi.NewRouter()
	apis.RegisterRoutes(r)

	log.Fatal(http.ListenAndServe(":3000", r))
}
