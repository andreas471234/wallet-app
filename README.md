# wallet-app

Golang wallet demo app

## Setting up the app

- Install golang 1.16 or later
  - <https://go.dev/doc/install>
- Install golang-migrate for database migration
  - `brew install golang-migrate`
- Install go dependencies
  - `go get .`

## DB Setup

- Set the DB Config as per local configuration

  ``` bash
    export DBUSER=""
    export DBPASS=""
    export DBHOST=""
    export DBPORT=""
  ```

- Run migration
  - check the db setting in Makefile
  - check db connection
    - `make local-db`
  - run migration
    - `make migrate-up`
  - check for successfull migrations
    - `make run`

## Running the app

- Run the app
  - `go run .`
  - it will run on localhost:8080

## API Documentation

- <https://documenter.getpostman.com/view/32451372/2sA2r8149m>
