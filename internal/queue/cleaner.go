package queue

import (
	"fmt"
	"os/exec"
	"time"
)

func deleteResource(req DeleteRequest) {
	if req.DeploymentType == "k8s" {
		deleteKINDCluster(req.ClusterName)
	} else if req.DeploymentType == "vm" {
		deleteVM(req.VMName)
	}
}

// Delete KIND Cluster
func deleteKINDCluster(clusterName string) {
	if clusterName == "" {
		fmt.Println("⚠️ No cluster name provided, skipping deletion")
		return
	}
	fmt.Printf("🛑 Deleting KIND cluster: %s\n", clusterName)

	cmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ Failed to delete KIND cluster %s: %v\n", clusterName, err)
	} else {
		fmt.Printf("✅ KIND cluster %s deleted successfully\n", clusterName)
	}
}

// Delete VM (Placeholder: Replace with actual API call)
func deleteVM(vmName string) {
	if vmName == "" {
		fmt.Println("⚠️ No VM name provided, skipping deletion")
		return
	}
	fmt.Printf("🛑 Deleting VM instance: %s\n", vmName)

	time.Sleep(2 * time.Second)
	fmt.Printf("✅ VM instance %s deleted successfully\n", vmName)
}
