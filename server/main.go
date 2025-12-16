package main

import (
	"log"
	"os"
	"server/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	client := makeConnectionToDB()
	handler := handlers.NewShortedURLHandler(client.Database("urls"))
	observeHandler := handlers.NewObservabilityHandler(client)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Twój frontend
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/:id", handler.ResolveShort)
	r.GET("/links", handler.GetAllUrls)
	r.POST("/addShort", handler.AddShortURL)
	r.DELETE("/deleteShort/:id", handler.RemoveUrl)
	r.Run("localhost:8080")
	r.GET("/health/liveness", observeHandler.IsAvailable)
	r.GET("/health/readiness", observeHandler.IsReady)
}

func makeConnectionToDB() *mongo.Client {

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your MONGODB_URI variable")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}
