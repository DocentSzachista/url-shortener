package main

import (
	"log"
	"os"
	"server/handlers"

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

	r := gin.Default()
	r.GET("/:id", handler.ResolveShort)
	r.POST("/addShort", handler.AddShortURL)
	r.DELETE("/deleteShort/:id", handler.RemoveUrl)
	r.Run("localhost:8080")
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
