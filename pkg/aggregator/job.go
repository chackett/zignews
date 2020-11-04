package aggregator

import (
	"context"
	"log"
	"time"

	mobileapi "github.com/chackett/zignews/pkg/mobile-api"
	"github.com/chackett/zignews/pkg/rssprovider"
	"github.com/chackett/zignews/pkg/storage"
	"github.com/pkg/errors"
)

// Job is a runnable task that will retrieve news articles
type Job struct {
	articleRepo storage.ArticleRepository
	provider    NewsProvider
	Label       string
	chStop      chan struct{}
}

// NewsProvider defines functionality to retrieve news articles
type NewsProvider interface {
	Latest() ([]storage.Article, error)
	PollingFrequency() time.Duration
}

// NewJob returns a job based on the specified NewsProvider
func NewJob(label string, provider NewsProvider, articleRepo storage.ArticleRepository) (Job, error) {
	return Job{
		Label:       label,
		provider:    provider,
		articleRepo: articleRepo,
		chStop:      make(chan struct{}),
	}, nil
}

// Start a job
func (j *Job) Start() error {
	log.Printf("Starting job: %s", j.Label)
	for {
		select {
		case <-j.chStop:
			log.Printf("Stopping job: %s", j.Label)
			return nil
		default:
			latest, err := j.provider.Latest()
			if err != nil {
				log.Print(errors.Wrapf(err, "%s - get latest - %s", j.Label, err.Error()))
			}
			log.Printf("%s - Received %d articles", j.Label, len(latest))
			_, err = j.articleRepo.InsertArticles(context.Background(), latest)
			if err != nil {
				log.Printf("ERROR: Saving articles - Job: %s Error: %s", j.Label, err.Error())
			}
			time.Sleep(j.provider.PollingFrequency())
		}
	}
}

// Stop a job
func (j *Job) Stop() {
	j.chStop <- struct{}{}
}

// BuildJobs take an instance of `storage.ProviderRepository` and will return a slice of `Job`
// based on all valid providers that are available
func BuildJobs(provRepo storage.ProviderRepository, artRepo storage.ArticleRepository) ([]Job, error) {
	offset := 0
	count := 99999
	providers, err := provRepo.GetProviders(context.Background(), offset, count)
	if err != nil {
		return nil, errors.Wrap(err, "get providers")
	}

	var result []Job

	for _, prov := range providers {
		if _, ok := mobileapi.SupportedProviders[prov.Type]; !ok {
			log.Printf("Skipping provider `%s` - type `%s` not supported", prov.Label, prov.Type)
			continue
		}

		// A hack, creating RSS provider as I know only RSS is supported - Need to make this dynamic when others are supported.
		// Will use the supported providers map to expose a factory function.

		pollingFrequency := time.Duration(time.Second * time.Duration(prov.PollFrequencySeconds))
		p, err := rssprovider.NewRSSProvider(prov.Label, prov.FeedURL, pollingFrequency)
		if err != nil {
			log.Printf("ERROR: Unable to create RSSProvider - %s", err.Error())
			continue
		}
		j, err := NewJob(p.Label, p, artRepo)
		if err != nil {
			log.Printf("ERROR: Unable to create new Job - %s", err.Error())
			continue
		}
		result = append(result, j)
	}

	return result, nil
}
