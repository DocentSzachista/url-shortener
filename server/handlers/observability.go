package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Observability struct {
	mongoClient *mongo.Client
}

func NewObservabilityHandler(db *mongo.Client) *Observability {
	return &Observability{
		mongoClient: db,
	}
}

func (client *Observability) IsAvailable(c *gin.Context) {
	c.JSON(
		http.StatusAccepted, gin.H{"status": "ALIVE"},
	)
}
func (client *Observability) IsReady(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := client.mongoClient.Ping(ctx, nil); err != nil {
		c.JSON(503, gin.H{
			"status": "NOT READY",
			"mongo":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"status": "READY"})
}
