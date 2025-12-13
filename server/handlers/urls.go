package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"server/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UrlShortedHandler struct {
	Col      *mongo.Collection
	hostname string
}

func NewShortedURLHandler(db *mongo.Database) *UrlShortedHandler {
	hostname := os.Getenv("HOSTNAME")
	return &UrlShortedHandler{
		hostname: hostname,
		Col:      db.Collection("urls"),
	}
}

func (dbHandler *UrlShortedHandler) ResolveShort(c *gin.Context) {

	shortedParam := c.Param("id")
	var record models.ShortedURL

	err := dbHandler.Col.FindOne(context.Background(), bson.M{"shortedID": shortedParam}).Decode(&record)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, record.Url)
}

func (dbHandler *UrlShortedHandler) AddShortURL(c *gin.Context) {
	var urlToShort models.URL

	if err := c.ShouldBindJSON(&urlToShort); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.ShortedURL
	err := dbHandler.Col.FindOne(context.Background(), bson.M{"url": urlToShort.Url}).Decode(&existing)
	if err == nil {
		log.Printf("Found url in database. Returning %s", existing.ShortedId)
		c.JSON(
			http.StatusOK, existing,
		)
		return
	}

	existing = models.ShortedURL{
		Url:       urlToShort.Url,
		ShortedId: uuid.New().String()[:8],
	}

	_, insertErr := dbHandler.Col.InsertOne(context.TODO(), existing)

	if insertErr != nil {
		log.Fatalf("Save did not work out: %s", err)
		c.JSON(
			http.StatusInternalServerError, gin.H{"error": insertErr.Error()},
		)
		return
	}
	c.JSON(
		http.StatusOK, existing,
	)

}

func (dbHandler *UrlShortedHandler) RemoveUrl(c *gin.Context) {
	shortedParam := c.Param("id")

	deleted, err := dbHandler.Col.DeleteOne(context.TODO(), bson.M{"shortedID": shortedParam})
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}
	c.JSON(http.StatusOK, deleted)

	log.Print("Removed URL from database")
}

func (dbHandler *UrlShortedHandler) GetAllUrls(c *gin.Context) {
	cursor, err := dbHandler.Col.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var results []models.ShortedURL
	for cursor.Next(context.Background()) {
		var record models.ShortedURL
		if err := cursor.Decode(&record); err != nil {
			log.Printf("Failed to decode record: %v", err)
			continue
		}
		results = append(results, record)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
