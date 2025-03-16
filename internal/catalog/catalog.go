package catalog

import (
	"encoding/json"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/models"
	"net/http"
)

// AddApplication API (only for publishers)
func AddApplication(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string                `json:"name"`
		Description string                `json:"description"`
		PublisherID uint                  `json:"publisher_id"`
		Deployment  models.DeploymentSpec `json:"deployment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate Deployment Type
	if req.Deployment.Type != "k8s" && req.Deployment.Type != "vm" {
		http.Error(w, "Invalid deployment type", http.StatusBadRequest)
		return
	}

	app := models.Application{
		Name:        req.Name,
		Description: req.Description,
		PublisherID: req.PublisherID,
		Deployment:  req.Deployment,
	}

	if err := database.DB.Create(&app).Error; err != nil {
		http.Error(w, "Failed to add application", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(app)
}

func ListApplications(w http.ResponseWriter, r *http.Request) {
	var apps []models.Application
	database.DB.Preload("Publisher").Find(&apps)
	json.NewEncoder(w).Encode(apps)
}
