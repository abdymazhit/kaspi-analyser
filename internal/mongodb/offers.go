package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kaspi-analyser/internal/models"
)

func getOffersCollection(dhm DBHandlerMongo) *mongo.Collection {
	return dhm.Database(MAIN_DATABASE).Collection(OFFERS_COLLECTION)
}

// -----------------------------------------------------------------------------
// Main methods
// -----------------------------------------------------------------------------

func (dhm DBHandlerMongo) SaveOffer(ctx context.Context, id string, offer interface{}) error {
	res := getOffersCollection(dhm).FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": offer}, options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (dhm DBHandlerMongo) GetMerchantsFromOffers(ctx context.Context) ([]models.Merchant, error) {
	cur, err := getOffersCollection(dhm).Aggregate(ctx, []bson.M{
		{"$group": bson.M{
			"_id":                     "$merchantId",
			"merchantId":              bson.M{"$first": "$merchantId"},
			"merchantName":            bson.M{"$first": "$merchantName"},
			"merchantRating":          bson.M{"$first": "$merchantRating"},
			"merchantReviewsQuantity": bson.M{"$first": "$merchantReviewsQuantity"},
		}},
	})
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cur.All(ctx, &results); err != nil {
		return nil, err
	}

	var merchants []models.Merchant
	for _, result := range results {
		merchant := models.Merchant{
			ID:   result["merchantId"].(string),
			Name: result["merchantName"].(string),
		}
		if result["merchantRating"] != nil {
			merchant.Rating = result["merchantRating"].(float64)
		}
		if result["merchantReviewsQuantity"] != nil {
			merchant.ReviewsCount = int64(result["merchantReviewsQuantity"].(float64))
		}
		merchants = append(merchants, merchant)
	}
	return merchants, nil
}
