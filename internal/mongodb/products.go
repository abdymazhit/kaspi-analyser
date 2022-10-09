package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getProductsCollection(dhm DBHandlerMongo) *mongo.Collection {
	return dhm.Database(MAIN_DATABASE).Collection(PRODUCTS_COLLECTION)
}

// -----------------------------------------------------------------------------
// Main methods
// -----------------------------------------------------------------------------

func (dhm DBHandlerMongo) SaveProduct(ctx context.Context, id string, product interface{}) error {
	res := getProductsCollection(dhm).FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": product}, options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
