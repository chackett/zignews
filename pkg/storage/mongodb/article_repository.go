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

const collectionArticles = "articles"

// ArticleRepository is an implementation to create/retrieve articles from MongoBD store
type ArticleRepository struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewArticleRepository ..
func NewArticleRepository(connection, user, password, dbName string) (*ArticleRepository, error) {
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

	result := ArticleRepository{
		client:   _client,
		database: _database,
	}

	err = result.createIndexes()
	if err != nil {
		return nil, errors.Wrap(err, "createIndexes()")
	}

	return &result, nil
}

// InsertArticles inserts a collection of articles into the article collection
func (pr *ArticleRepository) InsertArticles(ctx context.Context, articles []storage.Article) ([]string, error) {
	c := pr.database.Collection(collectionArticles)
	if c == nil {
		return nil, fmt.Errorf("unable to get collection handler for %s", collectionArticles)
	}

	for _, art := range articles {
		filter := &bson.M{
			"guid": art.GUID,
		}
		updateOpts := &options.UpdateOptions{
			Upsert: &[]bool{true}[0],
		}
		update := bson.M{
			"$set": art,
		}
		_, err := c.UpdateOne(ctx, filter, update, updateOpts)
		if err != nil {
			return nil, errors.Wrap(err, "insert one document")
		}
	}

	return nil, nil
}

// GetArticles returns a collection of article
func (pr *ArticleRepository) GetArticles(ctx context.Context, offset, count int, categories, providers []string) ([]storage.Article, error) {
	var results []storage.Article

	coll := pr.database.Collection(collectionArticles)
	if coll == nil {
		return nil, fmt.Errorf("unable to get collection handler for %s", collectionArticles)
	}
	filter := bson.D{
		// {
		// 	"provider", bson.D{
		// 		{
		// 			"$in", bson.A{providers},
		// 		},
		// 	},
		// },
		// {
		// 	"category", bson.D{
		// 		{
		// 			"$in", bson.A{categories},
		// 		},
		// 	},
		// },
	}
	options := options.Find().SetSkip(int64(offset * count)).SetLimit(int64(count))

	crs, err := coll.Find(ctx, filter, options)
	if err != nil {
		return nil, errors.Wrap(err, "execute find query")
	}
	err = crs.All(ctx, results)
	if err == nil {
		return nil, errors.Wrap(err, "decode all results")
	}

	return results, nil
}
func (pr *ArticleRepository) createIndexes() error {
	// Bootstrap the Mongo DB repo here. This is very much an afterthought and needs its own home and improving.
	models := mongo.IndexModel{
		Keys:    bson.D{{Key: "guid", Value: 1}},
		Options: options.Index().SetName("guid").SetUnique(true),
	}
	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
	_, err := pr.database.Collection(collectionArticles).Indexes().CreateOne(context.Background(), models, opts)
	if err != nil {
		return errors.Wrap(err, "create article GUID index")
	}

	log.Println("Created MongoDB indexes")
	return nil
}
