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

	filter := bson.D{}
	options := options.Find().SetSkip(int64(offset * count)).SetLimit(int64(count))

	crs, err := coll.Find(ctx, filter, options)
	if err != nil {
		return nil, errors.Wrap(err, "execute find query")
	}
	err = crs.All(ctx, &results)
	if err != nil {
		return nil, errors.Wrap(err, "decode all results")
	}
	return results, nil
}

func (pr *ProviderRepository) bootStrapProviders() error {
	existing, err := pr.GetProviders(context.Background(), 0, 9999)
	if err != nil {
		return errors.Wrap(err, "get providers")
	}
	if len(existing) > 0 {
		log.Print("Skipping bootstrap providers as providers already exist.")
		return nil
	}
	providers := []storage.Provider{
		{
			Type:                 "rss",
			FeedURL:              "http://feeds.bbci.co.uk/news/uk/rss.xml",
			Label:                "BBC News UK",
			PollFrequencySeconds: 0,
		},
		{
			Type:                 "rss",
			FeedURL:              "http://feeds.bbci.co.uk/news/technology/rss.xml",
			Label:                "BBC News Technology",
			PollFrequencySeconds: 1,
		},
		{
			Type:                 "rss",
			FeedURL:              "http://feeds.skynews.com/feeds/rss/uk.xml",
			Label:                "Sky News UK",
			PollFrequencySeconds: 2,
		},
		{
			Type:                 "rss",
			FeedURL:              "http://feeds.skynews.com/feeds/rss/technology.xml",
			Label:                "Sky News Technology",
			PollFrequencySeconds: 3,
		},
	}
	_, err = pr.InsertProviders(context.Background(), providers)
	if err != nil {
		log.Fatal(errors.Wrap(err, "bootstrap providers"))
	}
	log.Print("Successfully bootstrapped provider list")
	return nil
}

// GetProvider returns the provider related to the specified `providerID`. Use `IsNotFound()` to determine if item was not found.
func (pr *ProviderRepository) GetProvider(ctx context.Context, providerID string) (storage.Provider, error) {
	coll := pr.database.Collection(collectionProviders)
	if coll == nil {
		return storage.Provider{}, fmt.Errorf("unable to get collection handler for %s", collectionProviders)
	}

	objID, err := primitive.ObjectIDFromHex(providerID)
	if err != nil {
		return storage.Provider{}, errors.Wrap(err, "ObjectIDFromHex()")
	}

	filter := bson.M{
		"_id": objID,
	}
	options := options.FindOne()
	mgResult := coll.FindOne(ctx, filter, options)
	err = mgResult.Err()
	if err != nil {
		return storage.Provider{}, errors.Wrap(err, "execute find query")
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return storage.Provider{}, storage.ErrNotFound{
				Message: fmt.Sprintf("No provider found for ID `%s`", providerID),
			}
		}
		return storage.Provider{}, errors.Wrap(err, "execute query")
	}

	var result storage.Provider
	err = mgResult.Decode(&result)
	if err != nil {
		return storage.Provider{}, errors.Wrap(err, "decode all results")
	}

	return result, nil
}
