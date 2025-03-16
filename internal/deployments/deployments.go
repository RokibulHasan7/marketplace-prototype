package deployments

import (
	"encoding/json"
	"fmt"
	"github.com/RokibulHasan7/marketplace-prototype/internal/queue"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/models"
	"net/http"
)

// DeployApplication API (only for consumers)
func DeployApplication(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ConsumerID    uint `json:"consumer_id"`
		ApplicationID uint `json:"application_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch application details
	var app models.Application
	if err := database.DB.First(&app, req.ApplicationID).Error; err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	// Initialize Deployment
	deployment := models.Deployment{
		ConsumerID:     req.ConsumerID,
		ApplicationID:  req.ApplicationID,
		DeploymentType: app.Deployment.Type,
	}

	// Store Deployment Record (Initial Status)
	if err := database.DB.Create(&deployment).Error; err != nil {
		http.Error(w, "Failed to save deployment record", http.StatusInternalServerError)
		return
	}

	// Push to Redis Queue for Asynchronous Processing
	err := queue.PushToQueue(queue.InstallRequest{
		DeploymentID: fmt.Sprintf("%d", deployment.ID),
		ConsumerID:   fmt.Sprintf("%d", req.ConsumerID),
		Application:  app.Name,
		DeployType:   app.Deployment.Type,
		RepoURL:      app.Deployment.RepoURL,
		ChartName:    app.Deployment.ChartName,
	})
	if err != nil {
		http.Error(w, "Failed to queue deployment", http.StatusInternalServerError)
		return
	}

	// Return Deployment ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Deployment request queued",
		"deploymentID": deployment.ID,
	})
}
