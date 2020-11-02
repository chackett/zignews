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
	for _, j := range a.jobs {
		if delay := a.delay(); delay > 0 {
			log.Printf("Starting job `%s` in %s", j.Label, delay)
			time.Sleep(delay)
		}

		err := j.Start()
		if err != nil {
			log.Printf("ERROR: Starting job `%s`. Error=%s", j.Label, err.Error())
		} else {
			log.Printf("Started job: `%s`", j.Label)
		}
	}
}

// Stop stops the aggregation jobs
func (a *Aggregator) Stop() {
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
