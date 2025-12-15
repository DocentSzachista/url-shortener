package handlers

import (
	"log"
	"net/http"
	"os"
	"server/models"
	"server/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UrlShortedHandler struct {
	hostname     string
	mongoService *services.UrlService
}

func NewShortedURLHandler(db *mongo.Database) *UrlShortedHandler {
	hostname := os.Getenv("HOSTNAME")
	return &UrlShortedHandler{
		hostname:     hostname,
		mongoService: services.NewUrlService(db),
	}
}

func (dbHandler *UrlShortedHandler) ResolveShort(c *gin.Context) {

	shortedParam := c.Param("id")
	sourceIp := c.ClientIP()
	link, err := dbHandler.mongoService.FindRedirection(shortedParam, sourceIp)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, link)
}

func (dbHandler *UrlShortedHandler) AddShortURL(c *gin.Context) {
	var urlToShort models.URL

	if err := c.ShouldBindJSON(&urlToShort); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := dbHandler.mongoService.CheckUrl(*urlToShort.Url)

	if err == nil {
		c.JSON(
			http.StatusOK, data,
		)
		return
	}

	data, err = dbHandler.mongoService.InsertShort(
		urlToShort,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK, data,
	)

}

func (dbHandler *UrlShortedHandler) RemoveUrl(c *gin.Context) {
	shortedParam := c.Param("id")

	err := dbHandler.mongoService.RemoveUrl(shortedParam)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)

	log.Print("Removed URL from database")
}

func (dbHandler *UrlShortedHandler) GetAllUrls(c *gin.Context) {

	results, err := dbHandler.mongoService.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Print("Returned list of urls")
	c.JSON(http.StatusOK, results)

}
