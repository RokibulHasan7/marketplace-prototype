package catalog

import (
	"encoding/json"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/models"
	"net/http"
	"strconv"
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
	// Get query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	publisherName := r.URL.Query().Get("publisher")       // Filter by publisher name
	deploymentType := r.URL.Query().Get("deploymentType") // Filter by deployment type (k8s/vm)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Query builder
	query := database.DB.Preload("Publisher").Model(&models.Application{})

	// Apply filters
	if publisherName != "" {
		query = query.Joins("JOIN users ON users.id = applications.publisher_id").
			Where("users.name = ?", publisherName)
	}
	if deploymentType != "" {
		query = query.Where("type = ?", deploymentType)
	}

	var total int64
	query.Count(&total)

	var applications []models.Application
	result := query.Limit(limit).Offset(offset).Find(&applications)
	if result.Error != nil {
		http.Error(w, "Failed to fetch applications", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"data":        applications,
		"page":        page,
		"limit":       limit,
		"total_items": total,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
