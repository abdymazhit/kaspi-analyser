package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getShopReviewsCollection(dhm DBHandlerMongo) *mongo.Collection {
	return dhm.Database(MAIN_DATABASE).Collection(SHOP_REVIEWS_COLLECTION)
}

// -----------------------------------------------------------------------------
// Main methods
// -----------------------------------------------------------------------------

func (dhm DBHandlerMongo) SaveShopReview(ctx context.Context, id string, shopReview interface{}) error {
	res := getShopReviewsCollection(dhm).FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": shopReview}, options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
