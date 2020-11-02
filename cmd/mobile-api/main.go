package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	mobileapi "github.com/chackett/zignews/pkg/mobile-api"
	"github.com/chackett/zignews/pkg/storage/mongodb"
	"github.com/pkg/errors"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

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

	svc, err := mobileapi.NewService(artRepo, provRepo)
	if err != nil {
		log.Fatal(errors.Wrap(err, "create new mobileapi service"))
	}

	httpAPI, err := mobileapi.NewHandler(svc, config.APIAddress)
	if err != nil {
		log.Fatal(errors.Wrap(err, "create new mobileapi http handler"))
	}

	fmt.Println("mobile api running..")

	log.Fatal(httpAPI.Start())

}
