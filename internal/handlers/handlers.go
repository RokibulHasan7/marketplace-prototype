package handlers

import (
	"github.com/RokibulHasan7/marketplace-prototype/internal/catalog"
	"github.com/RokibulHasan7/marketplace-prototype/internal/deployments"
	"github.com/RokibulHasan7/marketplace-prototype/internal/users"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux) {
	// User routes
	r.Route("/api/users", func(r chi.Router) {
		r.Post("/", users.CreateUser) // Create a new user
		r.Get("/", users.ListUsers)   // List all users
	})

	// Application catalog routes
	r.Route("/api/apps", func(r chi.Router) {
		r.Post("/", catalog.AddApplication)          // Add a new application
		r.Get("/", catalog.ListApplications)         // List all applications
		r.Get("/{id}", catalog.GetApplication)       // Get app details
		r.Put("/{id}", catalog.UpdateApplication)    // Update app
		r.Delete("/{id}", catalog.DeleteApplication) // Delete app
	})

	// Deployment routes
	r.Route("/api/deployments", func(r chi.Router) {
		r.Post("/", deployments.DeployApplication) // Deploy an application

		r.Get("/{id}", deployments.GetDeployment)       // Get deployment details
		r.Delete("/{id}", deployments.DeleteDeployment) // Remove deployment
	})
}
