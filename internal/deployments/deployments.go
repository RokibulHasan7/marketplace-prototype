package deployments

import (
	"encoding/json"
	"fmt"
	"github.com/RokibulHasan7/marketplace-prototype/internal/queue"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
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

	// Check if the consumer exists
	var publisher models.User
	if err := database.DB.First(&publisher, req.ConsumerID).Error; err != nil {
		http.Error(w, "Publisher not found", http.StatusBadRequest)
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
		DeploymentID:  fmt.Sprintf("%d", deployment.ID),
		ConsumerID:    fmt.Sprintf("%d", req.ConsumerID),
		ApplicationID: fmt.Sprintf("%d", req.ApplicationID),
		Application:   app.Name,
		DeployType:    app.Deployment.Type,
		RepoURL:       app.Deployment.RepoURL,
		ChartName:     app.Deployment.ChartName,
		Inputs:        app.Inputs,
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

func GetDeployment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // Get deployment ID from URL

	var deployment models.Deployment
	if err := database.DB.
		Preload("Application").
		Preload("Consumer").
		First(&deployment, id).Error; err != nil {
		http.Error(w, "Deployment not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(deployment)
}

func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var deployment models.Deployment
	if err := database.DB.First(&deployment, id).Error; err != nil {
		http.Error(w, "Deployment not found", http.StatusNotFound)
		return
	}

	// Queue a delete message
	queue.PushToDeleteQueue(queue.DeleteRequest{
		DeploymentID:   id,
		DeploymentType: deployment.DeploymentType,
		ClusterName:    deployment.ClusterName,
		VMName:         deployment.VMName,
	})

	// Fetch Billing Record
	var billing models.BillingRecord
	if err := database.DB.Where("deployment_id = ?", id).First(&billing).Error; err != nil {
		http.Error(w, "Billing record not found", http.StatusNotFound)
		return
	}

	// Calculate final amount
	endTime := time.Now()
	totalHours := endTime.Sub(billing.StartTime).Hours()
	totalCost := totalHours * billing.HourlyRate

	// Update Billing Record with final amount
	billing.EndTime = &endTime
	billing.Amount = totalCost
	if err := database.DB.Save(&billing).Error; err != nil {
		http.Error(w, "Failed to save billing record", http.StatusInternalServerError)
		return
	}

	// Delete deployment from DB
	if err := database.DB.Delete(&deployment).Error; err != nil {
		http.Error(w, "Failed to delete deployment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
