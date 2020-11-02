package mobileapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/chackett/zignews/pkg/storage"
	"github.com/pkg/errors"
)

const (
	// maxPageSize defines the maximum number of results that can be returned in a response
	maxPageSize = 50
	// defaultPageSize defines the page size if unspecified or exceeds the maximum value
	defaultPageSize = 20
	// defaultOffset defines the default offset value if unspecified or unsafe
	defaultOffset = 0
	// minPollFrequencyMinutes defines the minimum allowable polling frequency
	minPollFrequencySeconds = 10 // Useful to prevent possible rate limiting / abuse
	// maxPollFrequencyMinutes defines the maximum allowable polling frequency
	maxPollFrequencySeconds = 60 * 60 * 24 // Not really sure we need an upper limit but nice to have configurability
)

// SupportedProviders contains a list of "type" of providers that are supported by the aggregator.
// This map could later on use a factory function to create an instance of associated provider, instead of the placeholder `interface{}`
var SupportedProviders map[string]interface{} = map[string]interface{}{"rss": nil}

// ServiceImpl implements the domain logic for the mobile api
type ServiceImpl struct {
	articles  storage.ArticleRepository
	providers storage.ProviderRepository
}

// NewService returns a new instance of the service mobile api service domain logic implementation
func NewService(articleRepo storage.ArticleRepository, providerRepo storage.ProviderRepository) (*ServiceImpl, error) {
	if articleRepo == nil {
		return nil, errors.New("article repository is nil")
	}
	if providerRepo == nil {
		return nil, errors.New("provider repository is nil")
	}
	return &ServiceImpl{
		articles:  articleRepo,
		providers: providerRepo,
	}, nil
}

// GetArticles retrieves a list of articles from underlying storage. Filter parameters can be used to reduce the results
func (s *ServiceImpl) GetArticles(ctx context.Context, offset, count int, category, provider []string) ([]storage.Article, error) {
	if count < 0 || count > maxPageSize {
		count = defaultPageSize
	}
	if offset < 0 {
		offset = defaultOffset
	}

	articles, err := s.articles.GetArticles(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get articles from repository")
	}

	return articles, nil
}

// SaveProvider saves a provider to underlying storage
func (s *ServiceImpl) SaveProvider(ctx context.Context, provider storage.Provider) (string, error) {
	if provider.Label == "" {
		return "", errors.New("provider label is required")
	}
	if _, ok := SupportedProviders[provider.Type]; !ok {
		return "", fmt.Errorf("`%s` is an unsupported provider type", provider.Type)
	}
	_, err := url.ParseRequestURI(provider.FeedURL)
	if err != nil {
		return "", errors.New("provider feedURL must be a valid URL")
	}

	// If the poll frequency is invalid we don't want to use default and result in unexpected behaviour by the operator.
	// Better to let them know.
	if provider.PollFrequencySeconds < minPollFrequencySeconds || provider.PollFrequencySeconds > maxPollFrequencySeconds {
		return "", fmt.Errorf("Invalid polling frequency, %f seconds. Must be inside range %d-%d", provider.PollFrequencySeconds, minPollFrequencySeconds, maxPollFrequencySeconds)
	}

	providerID, err := s.providers.InsertProviders(ctx, []storage.Provider{provider})
	if err != nil {
		return "", errors.Wrap(err, "save provider to repository")
	}

	return providerID[0], nil
}
