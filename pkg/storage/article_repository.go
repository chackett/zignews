package storage

import (
	"context"
	"time"
)

// Article defines a news article
type Article struct {
	Title       string    `json:"title,omitempty"`
	Link        string    `json:"link,omitempty"`
	Description string    `json:"description,omitempty"`
	Published   time.Time `json:"published,omitempty"`
	GUID        string    `json:"guid,omitempty"`
	Thumbnail   string    `json:"thumbnail,omitempty"`
	Categories  []string  `json:"categories,omitempty"`
	Provider    string    `json:"provider,omitempty"`
}

// ArticleRepository defines functionality to CRUD articles in underlying store
type ArticleRepository interface {
	InsertArticles(ctx context.Context, p []Article) ([]string, error)
	GetArticles(ctx context.Context, offset, count int, category, provider []string) ([]Article, error)
}
