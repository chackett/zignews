package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chackett/zignews/pkg/storage"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
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

	result := GenericRepository{
		client:   _client,
		database: _database,
	}

	err = result.createIndexes()
	if err != nil {
		return GenericRepository{}, errors.Wrap(err, "createIndexes()")
	}

	return result, nil
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

// InsertDocuments inserts a collection of documents
func (gr *GenericRepository) InsertDocuments(ctx context.Context, collection string, document []interface{}) ([]string, error) {
	c := gr.database.Collection(collection)
	if c == nil {
		return nil, fmt.Errorf("unable to get collection handler for %s", collection)
	}
	// updateOpts := &options.UpdateOptions{
	// 	Upsert: &[]bool{true}[0],
	// }
	// filter := &bson.D{{}}
	// update := bson.M{
	// 	"$set": document,
	// }
	// _, err := c.UpdateMany(ctx, filter, update, updateOpts)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "insert many documents")
	// }

	updateOpts := &options.UpdateOptions{
		Upsert: &[]bool{true}[0],
	}
	for _, doc := range document {
		d, _ := doc.(storage.Article)
		filter := &bson.M{
			"guid": d.GUID,
		}

		update := bson.M{
			"$set": doc,
		}
		_, err := c.UpdateOne(ctx, filter, update, updateOpts)
		if err != nil {
			return nil, errors.Wrap(err, "insert many documents")
		}
	}

	// var insertedIDs []string

	// Don't really need to check for nil result here. I subscribe to school of thought that because the called function is supposed
	// to return a value, the sooner this breaks the better. As opposed to checking for a nil value and the returned error, potentially being ignored, but the items
	// have been inserted in DB. If the library is broken, it would be better of blowing up the application in test.
	// for _, id := range result. {
	// 	insertedIDs = append(insertedIDs, id.(primitive.ObjectID).String())
	// }

	// Skipping inserted id's as I need to fix how it works with "upsert" operations.

	return nil, nil
}

// UpsertDocument Upserts a document
func (gr *GenericRepository) UpsertDocument(ctx context.Context, collection string, document interface{}, filter interface{}) ([]string, error) {
	c := gr.database.Collection(collection)
	if c == nil {
		return nil, fmt.Errorf("unable to get collection handler for %s", collection)
	}
	updateOpts := &options.UpdateOptions{
		Upsert: &[]bool{true}[0],
	}
	update := bson.M{
		"$set": document,
	}
	_, err := c.UpdateOne(ctx, filter, update, updateOpts)
	if err != nil {
		return nil, errors.Wrap(err, "insert one document")
	}

	// var insertedIDs []string

	// Don't really need to check for nil result here. I subscribe to school of thought that because the called function is supposed
	// to return a value, the sooner this breaks the better. As opposed to checking for a nil value and the returned error, potentially being ignored, but the items
	// have been inserted in DB. If the library is broken, it would be better of blowing up the application in test.
	// for _, id := range result. {
	// 	insertedIDs = append(insertedIDs, id.(primitive.ObjectID).String())
	// }

	// Skipping inserted id's as I need to fix how it works with "upsert" operations.

	return nil, nil
}

// Disconnect closes connection to repository if required
func (gr *GenericRepository) Disconnect(ctx context.Context) error {
	err := gr.client.Disconnect(ctx)
	if err != nil {
		return errors.Wrap(err, "mongo client.Disconnect()")
	}
	return nil
}

func (gr *GenericRepository) createIndexes() error {
	// Bootstrap the Mongo DB repo here. This is very much an afterthought and needs its own home.
	models := mongo.IndexModel{
		Keys:    bson.D{{Key: "guid", Value: 1}},
		Options: options.Index().SetName("guid").SetUnique(true),
	}
	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
	_, err := gr.database.Collection(collectionArticles).Indexes().CreateOne(context.Background(), models, opts)
	if err != nil {
		return errors.Wrap(err, "create article GUID index")
	}

	log.Println("Created MongoDB indexes")
	return nil
}
