package models

type Merchant struct {
	ID           string  `json:"id" bson:"_id"`
	Name         string  `json:"name" bson:"name"`
	Rating       float64 `json:"rating" bson:"rating"`
	ReviewsCount int64   `json:"reviewsCount" bson:"reviewsCount"`
	PhoneNumber  string  `json:"phoneNumber" bson:"phoneNumber"`
	CreatedAt    string  `json:"createdAt" bson:"createdAt"`
}
