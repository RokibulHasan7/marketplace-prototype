package helm

import (
	"fmt"
	"os/exec"
)

// DeployHelmChart deploys a Helm chart onto a Kind cluster
func DeployHelmChart(clusterName, repoURL, chartName string) error {
	// Add the repo first
	addRepoCmd := exec.Command("helm", "repo", "add", "myrepo", repoURL)
	if output, err := addRepoCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add repo: %v\n%s", err, string(output))
	}

	// Update the Helm repo
	updateRepoCmd := exec.Command("helm", "repo", "update")
	if output, err := updateRepoCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to update repo: %v\n%s", err, string(output))
	}

	// Install the Helm chart
	installCmd := exec.Command("helm", "install", chartName, "myrepo/"+chartName, "--kube-context", "kind-"+clusterName)
	output, err := installCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to deploy Helm chart: %v\n%s", err, string(output))
	}

	fmt.Printf("✅ Helm chart %s deployed successfully on cluster %s\n", chartName, clusterName)
	return nil
}
