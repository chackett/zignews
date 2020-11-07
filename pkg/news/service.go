package service

import (
	"context"

	"github.com/chackett/zignews/pkg/cache"
	"github.com/chackett/zignews/pkg/storage"
	"github.com/pkg/errors"
)

// Service defines functionality for saving, retrieving and processing news articles and news providers
type Service struct {
	articles  storage.ArticleRepository
	providers storage.ProviderRepository
	cache     cache.Cache
}

// NewService returns a new instance of service
func NewService(articleRepo storage.ArticleRepository, providerRepo storage.ProviderRepository, cache cache.Cache) (Service, error) {
	return Service{}, nil
}

// GetArticles retrieves a collection of articles from underlying store
func (s *Service) GetArticles(ctx context.Context) ([]storage.Article, error) {
	articles, err := s.articles.GetArticles(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "articleRepo.GetArticles")
	}
	return articles, nil
}

// InsertArticles saves slice of articles to underlying store
func (s *Service) InsertArticles(ctx context.Context, articles []storage.Article) ([]string, error) {
	insertedIDs, err := s.articles.InsertArticles(ctx, articles)
	if err != nil {
		return nil, errors.Wrap(err, "articleRepo.InsertArticles")
	}
	return insertedIDs, nil
}

// GetProviders returns a list of providers
func (s *Service) GetProviders(ctx context.Context) ([]storage.Provider, error) {
	providers, err := s.providers.GetProviders(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "articleRepo.GetProviders")
	}
	return providers, nil
}
