# ZigNews

[![Go Report Card](https://goreportcard.com/badge/github.com/chackett/zignews)](https://goreportcard.com/report/github.com/chackett/zignews)

A news aggregator service. This tool will aggregate news articles from various sources and provide them to clients via a single API.

## Running the system

To bring up the entire system, use `$ docker-compose up`. The result will be a running aggregator with bootstrapped RSS feeds, which should be available via the API.

View the API documentation as hosted by Postman [here](https://documenter.getpostman.com/view/820576/TVYNYvQ1).

## High level function

* Aggregator polls list of pre-defined news sites (providers) and saves the article meta data to database.
* Mobile API reads the article data from database and provides them to clients in a single streamlined API.
* MongoDB database stores both the articles and provider configuration such as feed URLs, poll frequency, ttl for cache etc.
* Redis will cache API calls from clients so new queries are not executed for each client API call.
* The intention is to create an event based system. For example, the aggregator could send `new_article` for a specific provider, which might invalidate cached items. If a new provider is provided via API, it could send an event for the aggregator to pick it up and start a new job for that provider.

## Known issues / Design notes

* Currently the aggregator and the mobile api services connect to the database, breaking the 1:1 db:service rule. If I did this correctly/with more time, I would create a "news service" which would have an interface for configuration/saving and retrieving articles.
* As a result of the previous comment, there is a bit of a bodge, in that the "mobile" api is also where the api to save providers is hosted.
* The aggregator and apis are loosly coupled, meaning if the aggregator stops then news is still available with the caveat that the articles become "stale".
* Caching - At scale, the system would be read heavy i.e. More API clients reading articles than articles being submitted to the database. Unless, it ends up reading many many news sites for a small number of clients, but still.

## Project Structure

Went for a mono repo here that hosts a full "system"

* _root_ - hold some meta/build files etc.
* `cmd` - entrypoints to launch the various services.
* `Docker` - Store multiple docker files.
* `pkg` - main implementation files. Typically, re-usable packages and since this is a monorepo, there is a package that matches each executable (cmd) holding main implementation for that "service".
  * `aggregator` - Implementation of the news aggregator.
  * `mobile-api` - Implementation of mobile api service.
  * `storage` - Persistence implementation. Think cache, db, memory etc.
  * `rssprovider` - Implementation of `NewsProvider` that consumes RSS feeds.
  * `news` - A lightweight implementation for CRUDing news related information. Would ordinarily be it's own service and be single point of access to the database.
  * `cache` - Implement caching.
  * `events` - Implement event messaging and signalling between components.

## Tech stack

* [Go](https://golang.org/) v1.15.2 - application code
* [MongoDB](http://mongodb.com/) v3.6 - Persistence
* [NATS](https://nats.io/) v2.1.8 - Event messaging
* [Redis](https://redis.io) v6.0.9 - Caching
* [Docker](https://docker.com) v19.03.13- Containerisation

## Libraries used

* [Gofeed](github.com/mmcdole/gofeed) - RSS Parser
* [Gorilla Mux](https://github.com/gorilla/mux) - Power HTTP router
* [Mongo-Driver](https://github.com/mongodb/mongo-go-driver) - Official MongoDB Go driver
* [Env](https://github.com/caarlos0/env) - Envar parser
* [Go-Redis](https://github.com/go-redis/redis) - Official Redis Go client
* [Nats.go](https://github.com/nats-io/nats.go) - Official NATS Go client
* [Errors](https://github.com/pkg/errors) - Great package for exposing errors
