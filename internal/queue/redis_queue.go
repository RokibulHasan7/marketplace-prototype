package queue

import (
	"context"
	"fmt"
	"github.com/RokibulHasan7/marketplace-prototype/internal/helm"
	"github.com/RokibulHasan7/marketplace-prototype/internal/kubernetes"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/database"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/models"
	redis "github.com/redis/go-redis/v9"
	"log"
	"time"
)

// Redis client
var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6370",
})

// Queue Name
const QueueName = "install_queue"

// InstallRequest represents a message in the queue
type InstallRequest struct {
	DeploymentID string
	ConsumerID   string
	Application  string
	DeployType   string
	RepoURL      string
	ChartName    string
}

// PushToQueue pushes a message to Redis
func PushToQueue(req InstallRequest) error {
	ctx := context.Background()

	_, err := redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: QueueName,
		Values: map[string]interface{}{
			"deployment_id": req.DeploymentID,
			"consumer_id":   req.ConsumerID,
			"application":   req.Application,
			"deploy_type":   req.DeployType,
			"repo_url":      req.RepoURL,
			"chart_name":    req.ChartName,
		},
	}).Result()

	if err != nil {
		log.Println("❌ Failed to push to queue:", err)
	}
	return err
}

// StartConsumer processes deployment messages
func StartConsumer() {
	ctx := context.Background()
	log.Println("🚀 Redis Queue Consumer Started...")

	for {
		// Read new messages from the queue
		messages, err := redisClient.XRead(ctx, &redis.XReadArgs{
			Streams: []string{QueueName, "0"},
			Count:   1,
			Block:   0, // Blocks indefinitely
		}).Result()

		if err != nil {
			log.Println("❌ Error reading from queue:", err)
			time.Sleep(2 * time.Second) // Retry after delay
			continue
		}

		for _, stream := range messages {
			for _, message := range stream.Messages {
				deploymentID := message.Values["deployment_id"].(string)
				consumerID := message.Values["consumer_id"].(string)
				application := message.Values["application"].(string)
				deployType := message.Values["deploy_type"].(string)
				repoURL := message.Values["repo_url"].(string)
				chartName := message.Values["chart_name"].(string)

				fmt.Printf("📦 Processing Deployment %s for User %s: %s\n", deploymentID, consumerID, application)

				// Fetch Deployment Record
				var deployment models.Deployment
				if err := database.DB.First(&deployment, deploymentID).Error; err != nil {
					log.Println("❌ Deployment not found:", err)
					continue
				}

				// Process Based on Deployment Type
				switch deployType {
				case "k8s":
					clusterName := fmt.Sprintf("kind-cluster-%s", consumerID)

					// Create KIND Cluster
					if err := kubernetes.CreateKindCluster(clusterName); err != nil {
						log.Println("❌ Failed to create cluster:", err)
						continue
					}

					// Deploy Helm Chart
					if err := helm.DeployHelmChart(clusterName, repoURL, chartName); err != nil {
						log.Println("❌ Failed to deploy Helm chart:", err)
						continue
					}

					// Update Deployment Record
					deployment.ClusterName = clusterName

				case "vm":
					vmName := fmt.Sprintf("vm-%s", consumerID)
					deployment.VMName = vmName
					deployment.VMIP = "10.2.0.1" // Mocking for now

				default:
					log.Println("❌ Invalid deployment type:", deployType)
					continue
				}

				// Mark Deployment as Completed
				if err := database.DB.Save(&deployment).Error; err != nil {
					log.Println("❌ Failed to update deployment record:", err)
				}

				fmt.Printf("✅ Deployment %s completed for %s\n", deploymentID, application)

				// Remove Processed Message from Queue
				_, err := redisClient.XDel(ctx, QueueName, message.ID).Result()
				if err != nil {
					log.Println("❌ Failed to acknowledge message:", err)
				}
			}
		}
	}
}
