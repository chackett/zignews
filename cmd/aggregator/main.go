package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/caarlos0/env"
	"github.com/chackett/zignews/pkg/aggregator"
	"github.com/chackett/zignews/pkg/storage/mongodb"
	"github.com/pkg/errors"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Print("Zignews aggregator")

	var config Config
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "parse config"))
	}

	artRepo, err := mongodb.NewArticleRepository(config.MongoAddress, config.MongoUser, config.MongoPass, config.MongoDatabase)
	if err != nil {
		log.Fatal(errors.Wrap(err, "create article repository"))
	}

	provRepo, err := mongodb.NewProviderRepository(config.MongoAddress, config.MongoUser, config.MongoPass, config.MongoDatabase)
	if err != nil {
		log.Fatal(errors.Wrap(err, "create article repository"))
	}

	jobs, err := aggregator.BuildJobs(provRepo, artRepo)
	if err != nil {
		log.Fatal(errors.Wrap(err, "aggregator BuildJobs()"))
	}
	log.Printf("Found %d jobs.", len(jobs))

	agg, err := aggregator.NewAggregator(jobs, config.DelayJobStart)
	if err != nil {
		log.Fatal(errors.Wrap(err, "create aggregator"))
	}

	go agg.Start()

	// Hold application until instructed to quit - it's a long running process
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	<-c
	agg.Stop()
	log.Println("shutting down")
	os.Exit(0)
}
