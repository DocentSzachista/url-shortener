package models

import (
	"time"
)

type URL struct {
	Url       *string `json:"url" binding:"required"`
	ShortedId *string `json:"shortedId"`
}

type ShortedURL struct {
	Url       string `json:"url" bson:"url"`
	ShortedId string `json:"shortedId" bson:"shortedID"`
}

type VisitedTimestamp struct {
	ShortedId string    `bson:"shortedId"`
	CreatedAt time.Time `bson:"created_at"`
	SourceIp  string    `bson:"source_ip"`
}
