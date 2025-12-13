package models

type URL struct {
	Url string
}

type ShortedURL struct {
	Url       string `bson:"url"`
	ShortedId string `bson:"shortedID"`
}
