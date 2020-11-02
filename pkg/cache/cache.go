package cache

import "context"

// Cache defines functionality to cache and retrieve items from an underlying store
type Cache interface {
	// Get tries to retrieve an item from cache using the specified key. If the item does not exist, nil is returned.
	// Value is returned by means of a reference argument `typ`
	Get(ctx context.Context, key string, typ interface{}) error
	// Set will add an item to the cache using the key. If an item exists, it is overwritten.
	Set(ctx context.Context, key string, value interface{}) error
}
