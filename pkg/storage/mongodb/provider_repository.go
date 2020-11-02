package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/chackett/zignews/pkg/storage"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionProviders = "providers"

// ProviderRepository is an implementation to create/retrieve providers from MongoBD store
type ProviderRepository struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewProviderRepository ..
func NewProviderRepository(connection, user, password, dbName string) (*ProviderRepository, error) {
	URI := fmt.Sprintf("mongodb://%s:%s@%s", user, password, connection)
	_client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		return nil, errors.Wrap(err, "mongo.NewClient()")
	}
	err = _client.Connect(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "mongo client.Connect()")
	}

	_database := _client.Database(dbName)

	result := &ProviderRepository{
		client:   _client,
		database: _database,
	}

	result.bootStrapProviders()

	return result, nil
}

// InsertProviders inserts a provider into the provider collection
func (pr *ProviderRepository) InsertProviders(ctx context.Context, providers []storage.Provider) ([]string, error) {
	// Need to convert []storage.Provider to []interface{} manually.
	var provIface []interface{}
	for _, provider := range providers {
		provIface = append(provIface, provider)
	}
	c := pr.database.Collection(collectionProviders)
	if c == nil {
		return nil, fmt.Errorf("unable to get collection handler for %s", collectionProviders)
	}

	imr, err := c.InsertMany(ctx, provIface, nil)
	if err != nil {
		return nil, errors.Wrap(err, "insert many")
	}
	var insertedIDs []string
	for _, id := range imr.InsertedIDs {
		insertedIDs = append(insertedIDs, id.(primitive.ObjectID).Hex())
	}

	return insertedIDs, nil
}

// GetProviders returns a collection of providers
func (pr *ProviderRepository) GetProviders(ctx context.Context, offset, count int) ([]storage.Provider, error) {
	var results []storage.Provider
	coll := pr.database.Collection(collectionProviders)
	if coll == nil {
		return nil, fmt.Errorf("unable to get collection handler for %s", collectionProviders)
	}

	filter := bson.D{{}}
	options := options.Find().SetSkip(int64(offset * count)).SetLimit(int64(count))

	crs, err := coll.Find(ctx, filter, options)
	if err != nil {
		return nil, errors.Wrap(err, "execute find query")
	}
	err = crs.All(ctx, results)
	if err == nil {
		return nil, errors.Wrap(err, "decode all results")
	}

	return nil, nil
}

func (pr *ProviderRepository) bootStrapProviders() {
	providers := []storage.Provider{
		{
			Type:                 "rss",
			FeedURL:              "http://feeds.bbci.co.uk/news/uk/rss.xml",
			Label:                "BBC News UK",
			PollFrequencySeconds: 10,
		},
		{
			Type:                 "rss",
			FeedURL:              "http://feeds.bbci.co.uk/news/technology/rss.xml",
			Label:                "BBC News Technology",
			PollFrequencySeconds: 10,
		},
		{
			Type:                 "rss",
			FeedURL:              "http://feeds.skynews.com/feeds/rss/uk.xml",
			Label:                "Sky News UK",
			PollFrequencySeconds: 10,
		},
		{
			Type:                 "rss",
			FeedURL:              "http://feeds.skynews.com/feeds/rss/technology.xml",
			Label:                "Sky News Technology",
			PollFrequencySeconds: 10,
		},
	}
	_, err := pr.InsertProviders(context.Background(), providers)
	if err != nil {
		log.Fatal(errors.Wrap(err, "bootstrap providers"))
	}
}
