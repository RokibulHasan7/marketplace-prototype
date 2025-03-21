package provisioner

import (
	"fmt"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/models"
	"log"
)

type VMProvisioner struct {
	InstallReq InstallRequest
}

func (vp *VMProvisioner) Provision() error {
	vmName := fmt.Sprintf("vm-%s", vp.InstallReq.ConsumerID)
	log.Printf("🚀 Provisioning VM: %s", vmName)

	// Update deployment record
	return updateDeploymentVM(vp.InstallReq.DeploymentID, vmName, "10.2.0.1")
}

func updateDeploymentVM(deploymentID, vmName, vmIP string) error {
	return database.DB.Model(&models.Deployment{}).
		Where("id = ?", deploymentID).
		Updates(map[string]interface{}{
			"vm_name": vmName,
			"vm_ip":   vmIP,
		}).Error
}
