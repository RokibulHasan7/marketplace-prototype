package queue

import (
	"context"
	"fmt"
	redis "github.com/redis/go-redis/v9"
	"log"
	"time"
)

// Stream Name
const UninstallQueue = "delete_queue"

// DeleteRequest represents a deletion message in the queue
type DeleteRequest struct {
	DeploymentID   string
	DeploymentType string
	ClusterName    string
	VMName         string
}

// PushToDeleteQueue pushes a delete message to Redis
func PushToDeleteQueue(req DeleteRequest) error {
	ctx := context.Background()

	_, err := redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: UninstallQueue,
		Values: map[string]interface{}{
			"deployment_id":   req.DeploymentID,
			"deployment_type": req.DeploymentType,
			"cluster_name":    req.ClusterName,
			"vm_name":         req.VMName,
		},
	}).Result()

	if err != nil {
		log.Println("❌ Failed to push delete request to queue:", err)
	}
	return err
}

// StartDeleteConsumer processes delete messages
func StartDeleteConsumer() {
	ctx := context.Background()
	log.Println("🚀 Redis Delete Queue Consumer Started...")

	for {
		// Read messages from the delete queue
		messages, err := redisClient.XRead(ctx, &redis.XReadArgs{
			Streams: []string{UninstallQueue, "0"},
			Count:   1,
			Block:   0, // Blocks indefinitely
		}).Result()

		if err != nil {
			log.Println("❌ Error reading from delete queue:", err)
			time.Sleep(2 * time.Second) // Retry delay
			continue
		}

		for _, stream := range messages {
			for _, message := range stream.Messages {
				deleteReq := DeleteRequest{
					DeploymentID:   message.Values["deployment_id"].(string),
					DeploymentType: message.Values["deployment_type"].(string),
					ClusterName:    message.Values["cluster_name"].(string),
					VMName:         message.Values["vm_name"].(string),
				}

				fmt.Printf("🗑️  Processing Delete Request for Deployment %s\n", deleteReq.DeploymentID)

				// Perform deletion
				deleteResource(deleteReq)

				// Acknowledge message deletion
				_, err := redisClient.XDel(ctx, UninstallQueue, message.ID).Result()
				if err != nil {
					log.Println("❌ Failed to acknowledge delete message:", err)
				}
			}
		}
	}
}
