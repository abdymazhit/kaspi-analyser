package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// main database
	MAIN_DATABASE           = "new_main"
	PRODUCTS_COLLECTION     = "products"
	OFFERS_COLLECTION       = "offers"
	SHOPS_COLLECTION        = "shops"
	SHOP_REVIEWS_COLLECTION = "shop_reviews"
)

type DBConfigMongo struct {
	URI string
}

type DBHandlerMongo struct {
	*mongo.Client
}

func NewDBHandlerMongo(ctx context.Context, cfg DBConfigMongo) (*DBHandlerMongo, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return &DBHandlerMongo{client}, nil
}
