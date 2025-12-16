package services

import (
	"context"
	"errors"
	"log"
	"server/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewUrlService(db *mongo.Database) *UrlService {
	return &UrlService{
		urlCol:   db.Collection("urls"),
		statsCol: db.Collection("timestamps"),
	}
}

type UrlService struct {
	urlCol   *mongo.Collection
	statsCol *mongo.Collection
}

func (service *UrlService) FindRedirection(shortedParam string, sourceIp string) (string, error) {
	var record models.ShortedURL

	err := service.urlCol.FindOne(context.Background(), bson.M{"shortedID": shortedParam}).Decode(&record)

	if err != nil {
		return "", err
	}
	mark := models.VisitedTimestamp{
		ShortedId: shortedParam,
		CreatedAt: time.Now().UTC(),
		SourceIp:  sourceIp,
	}

	_, insertErr := service.statsCol.InsertOne(context.Background(), mark)
	_, updateErr := service.urlCol.UpdateOne(context.Background(), bson.M{"shortedID": shortedParam}, bson.M{"$inc": bson.M{"clicks": 1}})
	if insertErr != nil {
		return "", insertErr
	}
	if updateErr != nil {
		return "", updateErr
	}

	return record.Url, nil
}

func (service *UrlService) InsertShort(urlToInsert models.URL) (*models.ShortedURL, error) {

	if urlToInsert.ShortedId == nil {
		id := uuid.New().String()[:8]
		urlToInsert.ShortedId = &id
	} else if len(*urlToInsert.ShortedId) > 8 {
		log.Printf("Custom ID is too long")
		return nil, errors.New("custom short link id is too long")
	} else if len(*urlToInsert.ShortedId) < 8 {
		log.Printf("Custom ID is too short")
		return nil, errors.New("custom short link id is too short")
	}

	shortedUrl := models.ShortedURL{
		Url:       *urlToInsert.Url,
		ShortedId: *urlToInsert.ShortedId,
		Clicks:    0,
	}
	_, insertErr := service.urlCol.InsertOne(context.TODO(), shortedUrl)

	if insertErr != nil {
		log.Fatalf("Save did not work out: %s", insertErr)
		return nil, insertErr
	}
	return &shortedUrl, nil
}

func (service *UrlService) CheckUrl(url string) (*models.ShortedURL, error) {
	var data *models.ShortedURL

	err := service.urlCol.FindOne(context.Background(), bson.M{"url": url}).Decode(&data)
	if err == nil {
		log.Printf("Found url in database. Returning %s", data.ShortedId)
		return data, nil
	}
	return nil, err
}

func (service *UrlService) RemoveUrl(shortId string) error {

	_, err := service.urlCol.DeleteOne(context.Background(), bson.M{"shortedID": shortId})
	return err

}

func (service *UrlService) GetAll() ([]models.ShortedURL, error) {
	cursor, err := service.urlCol.Find(context.Background(), bson.M{})
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
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
		log.Printf("Cursor did not close correctly: %v", err)
		return nil, err
	}
	return results, nil
}
