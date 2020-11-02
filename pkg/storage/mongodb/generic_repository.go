package mongodb

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GenericRepository is a generic implementation to retrieve collections from MongoDB
type GenericRepository struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewGenericRepository returns a repository for generic retrieval
func NewGenericRepository(connection, user, password, dbName string) (GenericRepository, error) {
	URI := fmt.Sprintf("mongodb://%s:%s@%s", user, password, connection)
	_client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		return GenericRepository{}, errors.Wrap(err, "mongo.NewClient()")
	}
	err = _client.Connect(context.Background())
	if err != nil {
		return GenericRepository{}, errors.Wrap(err, "mongo client.Connect()")
	}

	_database := _client.Database(dbName)

	return GenericRepository{
		client:   _client,
		database: _database,
	}, nil
}

// GetCollection returns a named collection
func (gr *GenericRepository) GetCollection(ctx context.Context, collection string, results interface{}) error {
	coll := gr.database.Collection(collection)
	if coll == nil {
		return fmt.Errorf("unable to get collection handler for %s", collection)
	}

	filter := bson.D{{}}
	options := options.Find()
	options.SetLimit(100) // #Hack - Magic number, hardcoded for now.

	crs, err := coll.Find(ctx, filter, options)
	if err != nil {
		return errors.Wrap(err, "execute find query")
	}
	err = crs.All(ctx, results)
	if err == nil {
		return errors.Wrap(err, "decode all results")
	}

	return nil
}

// InsertDocuments inserts a document
func (gr *GenericRepository) InsertDocuments(ctx context.Context, collection string, document []interface{}) ([]string, error) {
	c := gr.database.Collection(collection)
	if c == nil {
		return nil, fmt.Errorf("unable to get collection handler for %s", collection)
	}

	result, err := c.InsertMany(ctx, document)
	if err != nil {
		return nil, errors.Wrap(err, "insert single document")
	}

	var insertedIDs []string

	// Don't really need to check for nil result here. I subscribe to school of thought that because the called function is supposed
	// to return a value, the sooner this breaks the better. As opposed to checking for a nil value and the returned error, potentially being ignored, but the items
	// have been inserted in DB. If the library is broken, it would be better of blowing up the application in test.
	for _, id := range result.InsertedIDs {
		insertedIDs = append(insertedIDs, id.(primitive.ObjectID).String())
	}

	return insertedIDs, nil
}

// Disconnect closes connection to repository if required
func (gr *GenericRepository) Disconnect(ctx context.Context) error {
	err := gr.client.Disconnect(ctx)
	if err != nil {
		return errors.Wrap(err, "mongo client.Disconnect()")
	}
	return nil
}
