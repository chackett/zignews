package mongodb

import (
	"context"

	"github.com/chackett/zignews/pkg/storage"
	"github.com/pkg/errors"
)

const collectionArticles = "articles"

// ArticleRepository is an implementation to create/retrieve articles from MongoBD store
type ArticleRepository struct {
	generic GenericRepository
}

// NewArticleRepository ..
func NewArticleRepository(connection, user, password, dbName string) (*ArticleRepository, error) {
	// A big no no, usually. But in this case, it is the choice of "article repository" to use generic repository.
	// Usually would pass in an implementation instead of creating inside a constructor

	gr, err := NewGenericRepository(connection, user, password, dbName)
	if err != nil {
		return nil, errors.Wrap(err, "NewGenericRepository()")
	}

	return &ArticleRepository{
		generic: gr,
	}, nil
}

// InsertArticles inserts a collection of articles into the article collection
func (pr *ArticleRepository) InsertArticles(ctx context.Context, articles []storage.Article) ([]string, error) {
	// Need to convert []storage.Article to []interface{} manually.
	var artIface []interface{}
	for _, article := range articles {
		artIface = append(artIface, article)
	}

	articleIDs, err := pr.generic.InsertDocuments(ctx, collectionArticles, artIface)
	if err != nil {
		return nil, errors.Wrap(err, "generic insert document")
	}
	return articleIDs, nil
}

// GetArticles returns a collection of article
func (pr *ArticleRepository) GetArticles(ctx context.Context) ([]storage.Article, error) {
	var results []storage.Article
	err := pr.generic.GetCollection(ctx, collectionArticles, &results)
	if err != nil {
		return nil, errors.Wrap(err, "generic get collection")
	}
	return nil, nil
}
