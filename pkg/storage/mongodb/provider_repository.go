package mongodb

import (
	"context"

	"github.com/chackett/zignews/pkg/storage"
	"github.com/pkg/errors"
)

const collectionProviders = "providers"

// ProviderRepository is an implementation to create/retrieve providers from MongoBD store
type ProviderRepository struct {
	generic GenericRepository
}

// NewProviderRepository ..
func NewProviderRepository(connection, user, password, dbName string) (*ProviderRepository, error) {
	// A big no no, usually. But in this case, it is the choice of "provider repository" to use generic repository.
	// Usually would pass in an implementation instead of creating inside a constructor

	gr, err := NewGenericRepository(connection, user, password, dbName)
	if err != nil {
		return nil, errors.Wrap(err, "NewGenericRepository()")
	}

	return &ProviderRepository{
		generic: gr,
	}, nil
}

// InsertProviders inserts a provider into the provider collection
func (pr *ProviderRepository) InsertProviders(ctx context.Context, providers []storage.Provider) ([]string, error) {
	// Need to convert []storage.Provider to []interface{} manually.
	var provIface []interface{}
	for _, provider := range providers {
		provIface = append(provIface, provider)
	}
	insertedIDs, err := pr.generic.InsertDocuments(ctx, collectionProviders, provIface)
	if err != nil {
		return nil, errors.Wrap(err, "generic insert document")
	}
	return insertedIDs, nil
}

// GetProviders returns a collection of providers
func (pr *ProviderRepository) GetProviders(ctx context.Context) ([]storage.Provider, error) {
	var results []storage.Provider
	err := pr.generic.GetCollection(ctx, collectionProviders, &results)
	if err != nil {
		return nil, errors.Wrap(err, "generic get collection")
	}
	return results, nil
}
