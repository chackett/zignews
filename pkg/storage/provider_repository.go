package storage

import (
	"context"
)

// Provider represents a provider
type Provider struct {
	Type                 string `json:"type,omitempty"`
	Label                string `json:"label,omitempty"`
	FeedURL              string `json:"feedURL,omitempty"`
	PollFrequencySeconds int    `json:"pollFrequencySeconds,omitempty"`
}

// ProviderRepository defines functionality to CRUD providers in underlying store
type ProviderRepository interface {
	InsertProviders(ctx context.Context, p []Provider) ([]string, error)
	GetProviders(ctx context.Context, offset, count int) ([]Provider, error)
}
