package provisioner

import (
	"fmt"
	"github.com/RokibulHasan7/marketplace-prototype/internal/services/deployments/utils/helm"
	"github.com/RokibulHasan7/marketplace-prototype/internal/services/deployments/utils/kubernetes"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/models"
	"log"
	"os/exec"
	"time"
)

type KubernetesProvisioner struct {
	InstallReq InstallRequest
}

func (kp *KubernetesProvisioner) Provision() error {
	clusterName := fmt.Sprintf("kind-cluster-%s-%s-%d", kp.InstallReq.ConsumerID, kp.InstallReq.ApplicationID, time.Now().Unix())
	log.Printf("🚀 Provisioning Kubernetes Cluster: %s", clusterName)

	// Fetch Deployment Record
	var deployment models.Deployment
	if err := database.DB.First(&deployment, kp.InstallReq.DeploymentID).Error; err != nil {
		log.Println("❌ Deployment not found:", err)
		return err
	}

	if err := kubernetes.CreateKindCluster(clusterName); err != nil {
		// Update status to "failed"
		database.DB.Model(&deployment).Update("status", "failed")
		return fmt.Errorf("❌ failed to create KIND cluster: %w", err)
	}

	if err := switchKubeContext(clusterName); err != nil {
		// Update status to "failed"
		database.DB.Model(&deployment).Update("status", "failed")
		return fmt.Errorf("❌ failed to switch context: %w", err)
	}

	if err := helm.DeployHelmChart(clusterName, kp.InstallReq.RepoURL, kp.InstallReq.ChartName, kp.InstallReq.Application); err != nil {
		// Update status to "failed"
		database.DB.Model(&deployment).Update("status", "failed")
		return fmt.Errorf("❌ failed to deploy Helm chart: %w", err)
	}

	// Update deployment record
	if err := updateDeploymentCluster(kp.InstallReq.DeploymentID, clusterName); err != nil {
		log.Println("failed to update deployment cluster")
		return err
	}

	// Mark Deployment as Completed
	if err := database.DB.Save(&deployment).Error; err != nil {
		log.Println("❌ Failed to update deployment record:", err)
	}

	// Update status to "installed"
	database.DB.Model(&deployment).Update("status", "installed")

	return nil
}

func updateDeploymentCluster(deploymentID, clusterName string) error {
	return database.DB.Model(&models.Deployment{}).
		Where("id = ?", deploymentID).
		Update("cluster_name", clusterName).Error
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
