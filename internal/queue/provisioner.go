package queue

import (
	"errors"
	"fmt"
	"github.com/RokibulHasan7/marketplace-prototype/internal/helm"
	"github.com/RokibulHasan7/marketplace-prototype/internal/kubernetes"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/models"
	"log"
	"os/exec"
	"strconv"
	"time"
)

func provisionApplication(installReq InstallRequest) error {
	var deployment models.Deployment
	if err := database.DB.First(&deployment, installReq.DeploymentID).Error; err != nil {
		log.Println("❌ Deployment not found:", err)
		return err
	}

	// Process Based on Deployment Type
	switch installReq.DeployType {
	case "k8s":
		clusterName := fmt.Sprintf("kind-cluster-%s-%s-%d", installReq.ConsumerID, installReq.ApplicationID, time.Now().Unix())

		// Create KIND Cluster
		if err := kubernetes.CreateKindCluster(clusterName); err != nil {
			log.Println("❌ Failed to create cluster:", err)
			return err
		}

		// Switch context to the newly created KIND cluster
		if err := switchKubeContext(clusterName); err != nil {
			log.Println("❌ Failed to switch context:", err)
			return err
		}

		// Deploy Helm Chart
		if err := helm.DeployHelmChart(clusterName, installReq.RepoURL, installReq.ChartName, installReq.Application); err != nil {
			log.Println("❌ Failed to deploy Helm chart:", err)
			return err
		}

		// Update Deployment Record
		deployment.ClusterName = clusterName

	case "vm":
		vmName := fmt.Sprintf("vm-%s", installReq.ConsumerID)
		deployment.VMName = vmName
		deployment.VMIP = "10.2.0.1" // Mocking for now

	default:
		log.Println("❌ Invalid deployment type:", installReq.DeployType)
		return errors.New("invalid deployment type")
	}

	// Record Billing Start
	addBillingRecord(installReq)

	return nil
}

func switchKubeContext(clusterName string) error {
	// Run the kubectl command to set the context to the newly created KIND cluster
	cmd := exec.Command("kind", "export", "kubeconfig", "--name", clusterName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch context: %s, output: %s", err, output)
	}
	log.Printf("Switched context to %s", clusterName)
	return nil
}

func addBillingRecord(installReq InstallRequest) {
	applicationID := installReq.ApplicationID
	uintID, err := strconv.ParseUint(applicationID, 10, 32)
	if err != nil {
		log.Println("failed to convert application id to uint:", err)
	}

	var app models.Application
	if err := database.DB.Preload("Publisher").First(&app, uintID).Error; err != nil {
		log.Println("failed to find application:", err)
	}
	billing := models.BillingRecord{
		ID:            fmt.Sprintf("%s-bill", installReq.DeploymentID),
		ConsumerID:    installReq.ConsumerID,
		DeploymentID:  installReq.DeploymentID,
		ApplicationID: app.ID,
		HourlyRate:    app.HourlyRate,
		Amount:        0.0,
		StartTime:     time.Now(),
		CreatedAt:     time.Now(),
	}

	if err := database.DB.Create(&billing).Error; err != nil {
		log.Println("❌ Failed to create deployment record:", err)
	}
}
