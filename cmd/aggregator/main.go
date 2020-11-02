package main

import (
	"log"

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

	agg.Start()
}
