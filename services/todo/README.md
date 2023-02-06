# TODO Service

The todo service is effectively the middleware service that acts as the glue between a UI and the backing database. It can be accessed via
a RESTful api for your basic CRUD operations.

## Database
After there's a running Mysql instance (via docker compose is an easy solution). Make sure you've connected to it and have created the todo database. See [the create_database sql script](_infra/db/create_database.sql) for a dead simple command. 

Schemas are located in db/migrations using the tool [go-migrate](https://github.com/golang-migrate/migrate). This tool generates the files for migrations,
and you as the developer implement your changes in plain old SQL. An example to create a migration is like this:
` migrate create -dir db/migrations -ext sql create_todo_table  `

To run the docker container migration, run:
`docker build -f Dockerfile.migrate -t todo/migrate:latest .`
`docker run -it --rm  --network todo_default todo/migrate`

from the [db folder](_infra/db)

These migrations can be run in a docker container before the service is launched. I know some people like to integrate their migrations into the
main service on startup. But that's always irked me as it too heavily couples service startup to the database. It's a nice check to be sure,
but I've had enough ops people get annoyed at this as they may want to validate various scenarios before attempting such a thing.
See Dockerfile.migrate [migration dockerfile](_infra/db/Dockerfile.migrate) for how these migrations are run. but it basically boils down to
`migrate -database "mysql://root:password@tcp(127.0.0.1:3306)/todo" -path db/migrations up`

### Seeding data
The todo service includes its own seed command for development.
after building the service, run `./todo seed --connstring="root:password@(localhost:3306)/todo"` to seed the database with some todo items that are in various states.

## Building:
`make build`

## Testing:
`make test`

`make cover` to generate an hteml frenly cover profile. Found in bin/cover.html


## Overview
The service does its best to divide our work by concern. We have a RESTful server that calls upon a "core" todo domain layer. This domain validates the request, and then calls the todo datastore. This data store is responsible for interacting with the database.

![todo service high level overview](../../diagrams/todo-plooto-todo-service.drawio.png)


### A quick note on go and sql
Go only recently included generics into its language spec. While there have been a few attempts at creating a go based ORM, none of them have been particularly great. Some rely _heavily_ on reflection and slow down the program (i.e: GORM) while others have a lot of boilerplate code generation that includes tests that seemingly require a running database (i.e: sqlboiler). 

I opted instead to use the popular `sqlx` library instead. This makes my datastore layer less flexible, but kept the code down. I believe I've avoided any obvious sql injections, but I'm bringing it up here for now.

If given a bit more time, I'd have probably chosen to go with sqlBoiler.