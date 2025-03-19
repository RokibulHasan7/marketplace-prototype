package billing

import (
	"fmt"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/models"
	"log"
	"time"
)

func StartBillingUpdater() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		<-ticker.C
		updateBillingRecords()
	}
}

func updateBillingRecords() {
	log.Println("🔄 Updating billing records...")

	// Fetch active billing records (where EndTime is NULL)
	var records []models.BillingRecord
	if err := database.DB.Where("end_time IS NULL").Find(&records).Error; err != nil {
		log.Println("❌ Failed to fetch billing records:", err)
		return
	}

	for _, record := range records {
		elapsedHours := time.Since(record.StartTime).Hours()
		newAmount := elapsedHours * record.HourlyRate

		// Update Billing Record
		if err := database.DB.Model(&record).Updates(models.BillingRecord{
			Amount:    newAmount,
			UpdatedAt: time.Now(),
		}).Error; err != nil {
			log.Println("❌ Failed to update billing:", err)
		} else {
			fmt.Printf("💰 Billing updated: %s → $%.2f\n", record.DeploymentID, newAmount)
		}
	}
}
