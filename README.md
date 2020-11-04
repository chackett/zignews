# ZigNews

A news aggregator service. This tool will aggregate news articles from various sources and provide them to clients via a single API.

## Running the system

## High level function

* Aggregator polls list of pre-defined news sites and saves the article meta data to database.
* Mobile API reads the article data from database and provides them to clients in a single streamlined API.
* PostgresSQL database stores both the articles and the configuration for application such as news sites, poll frequency, ttl for cache.
* Memcached will cachce API calls from clients so new queries are not executed for each client API call.
* The intention is to create an event based system. For example, the aggregator could send `new_article` for a specific provider, which might invalidate cached items. If a new news provider is provided, it could send an event for the aggregator to pick it up.

## Known issues / Design notes

* Currently the aggregator and the mobile api services connect to the database, breaking the 1:1 db:service rule. If I did this correctly/with more time, I would create a "news service" which would have an interface for configuration/saving and retrieving articles.
* As a result of the previous comment, there is a bit of a bodge, in that the "mobile" api is also where the api to save providers is hosted.
* The aggregator and apis are loosly coupled, meaning if the aggregator stops then news is still available with the caveat that the articles become "stale".
* Caching - At scale, the system would be read heavy i.e. More API clients reading articles than articles being submitted to the database. Unless, it ends up reading many many news sites for a small number of clients, but still.

## Project Structure

Went for a mono repo here that hosts a full "system"

* _root_ - hold some meta/build files etc.
* `cmd` - entrypoints to launch the various services.
* `pkg` - main implementation files. Typically, re-usable packages and since this is a monorepo, there is a package that matches each executable (cmd) holding main implementation for that "service".
  * `aggregator` - Implementation of the news aggregator.
  * `mobile-api` - Implementation of mobile api service.
  * `storage` - Persistence implementation. Think cache, db, memory etc.

## Tech stack

* TO DO

## Library choice

* RSS