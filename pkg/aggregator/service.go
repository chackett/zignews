package aggregator

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/chackett/zignews/pkg/rssprovider"
	"github.com/chackett/zignews/pkg/storage"

	"github.com/chackett/zignews/pkg/events"
	"github.com/pkg/errors"

	"github.com/nats-io/nats.go"
)

// Aggregator orchestrates the creation and running of aggregation jobs
type Aggregator struct {
	jobs           []Job
	delayStarts    bool
	msgBus         *nats.Conn
	subNewProvider *nats.Subscription
	providers      storage.ProviderRepository
	articles       storage.ArticleRepository
}

// NewAggregator returns a new instance of Aggregator, which is used for aggregating news providers that implement `aggregator.NewProvider`
func NewAggregator(jobs []Job, delayStarts bool, msgBus *nats.Conn, providerRepo storage.ProviderRepository, articles storage.ArticleRepository) (Aggregator, error) {
	return Aggregator{
		jobs:        jobs,
		delayStarts: delayStarts,
		msgBus:      msgBus,
		providers:   providerRepo,
		articles:    articles,
	}, nil
}

// Start starts the aggregation jobs
func (a *Aggregator) Start() error {
	log.Print("Starting Aggregator")

	// Subscribe to event bus for new providers
	sub, err := a.msgBus.Subscribe(events.NewProvider, a.handleNewProvider)
	if err != nil {
		return errors.Wrap(err, "msgBus subscribe")
	}
	a.subNewProvider = sub

	// Start the jobs
	for i := range a.jobs {
		job := a.jobs[i]
		if delay := a.delay(); delay > 0 {
			log.Printf("Starting job `%s` in %s", job.Label, delay)
			time.Sleep(delay)
		}

		// Lazy way of starting a goroutine
		go func() {
			err := job.Start()
			if err != nil {
				log.Printf("ERROR: Starting job `%s`. Error=%s", job.Label, err.Error())
			}
		}()
	}
	return nil
}

// Stop stops the aggregation jobs
func (a *Aggregator) Stop() {
	log.Print("Aggregator stopping")

	// Unsubscribe from queue
	err := a.subNewProvider.Unsubscribe()
	if err != nil {
		log.Printf("ERROR: Unable to subscribe from queue `%s` - %s", events.NewProvider, err.Error())
		// don't return, continue stopping the jobs.
	}

	// Stop jobs
	for _, j := range a.jobs {
		j.Stop()
	}
}

// Delay calculates a delay if needed
func (a *Aggregator) delay() time.Duration {
	if a.delayStarts {
		return time.Duration(time.Second * time.Duration(rand.Intn(60)))
	}
	return 0
}

func (a *Aggregator) handleNewProvider(m *nats.Msg) {
	providerID := string(m.Data)
	log.Printf("New provider discovered with id `%s`", providerID)

	provider, err := a.providers.GetProvider(context.Background(), providerID)
	if err != nil {
		switch err := errors.Cause(err).(type) {
		case storage.ErrNotFound:
			fmt.Printf("ERROR: Unable to find the specified provider - %s", err.Error())
		default:
			fmt.Printf("ERROR: Get provider error - %s", err.Error())
		}
		return
	}

	// Build a new instance of provider based on newly received configuration
	pollingFrequency := time.Duration(time.Second * time.Duration(provider.PollFrequencySeconds))
	rssp, err := rssprovider.NewRSSProvider(provider.Label, provider.FeedURL, pollingFrequency)
	if err != nil {
		fmt.Printf("ERROR: New RSS Provider - %s", err.Error())
		return
	}
	job, err := NewJob(provider.Label, rssp, a.articles)
	if err != nil {
		fmt.Printf("ERROR: Create new job - %s", err.Error())
		return
	}

	// Lazy way to start a goroutine
	go func() {
		err := job.Start()
		if err != nil {
			log.Printf("ERROR: Starting job `%s`. Error=%s", job.Label, err.Error())
		}
	}()
}
