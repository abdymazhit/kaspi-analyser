package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func getShopsCollection(dhm DBHandlerMongo) *mongo.Collection {
	return dhm.Database(MAIN_DATABASE).Collection(SHOPS_COLLECTION)
}

// -----------------------------------------------------------------------------
// Main methods
// -----------------------------------------------------------------------------

func (dhm DBHandlerMongo) SaveShops(ctx context.Context, shops []interface{}) error {
	if _, err := getShopsCollection(dhm).InsertMany(ctx, shops); err != nil {
		return err
	}
	return nil
}
