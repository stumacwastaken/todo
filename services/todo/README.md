# TODO Service

The todo service is effectively the middleware service that acts as the glue between a UI and the backing database. It can be accessed via
a RESTful api for your basic CRUD operations

## Database
schemas are located in db/migrations using the tool go-migrate <<include URL here>>. This tool generates the files for migrations,
and you as the developer implement your changes in plain old SQL. An example to create a migration is like this:
` migrate create -dir db/migrations -ext sql create_todo_table  `

These migrations are run in a docker container before the service is launched. I know some people like to integrate their migrations into the
main service on startup. But that's always irked me as it too heavily couples service startup to the database. It's a nice check to be sure,
but I've had enough ops people get annoyed at this as they may want to validate various scenarios before attempting such a thing.
See Dockerfile.migrate <<insert link to Dockerfile.migrate>> for how these migrations are run. but it basically boils down to
`migrate -database "mysql://root:password@tcp(127.0.0.1:3306)/todo" -path db/migrations up`

## Building:
`make build`

## Testing:
`make test`

## Additional points to come here.