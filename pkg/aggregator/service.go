package aggregator

import (
	"log"
	"math/rand"
	"time"
)

// Aggregator orchestrates the creation and running of aggregation jobs
type Aggregator struct {
	jobs        []Job
	delayStarts bool
}

// NewAggregator returns a new instance of Aggregator, which is used for aggregating news providers that implement `aggregator.NewProvider`
func NewAggregator(jobs []Job, delayStarts bool) (Aggregator, error) {
	return Aggregator{
		jobs:        jobs,
		delayStarts: delayStarts,
	}, nil
}

// Start starts the aggregation jobs
func (a *Aggregator) Start() {
	log.Print("Starting Aggregator")
	for i := range a.jobs {
		job := a.jobs[i]
		if delay := a.delay(); delay > 0 {
			log.Printf("Starting job `%s` in %s", job.Label, delay)
			time.Sleep(delay)
		}

		go func() {
			err := job.Start()
			if err != nil {
				log.Printf("ERROR: Starting job `%s`. Error=%s", job.Label, err.Error())
			}
		}()
	}
}

// Stop stops the aggregation jobs
func (a *Aggregator) Stop() {
	log.Print("Aggregator stopping")
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
