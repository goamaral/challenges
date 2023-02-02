# Challenge
## Description
In this challenge I was asked to build a simple users CRUD api using GRPC and Golang.
I should also notify user events (create, update, delete).

## Improvements
- Pagination: We could have added a TotalCount in the response
- Validations: Only basic validations were implemented but could be improved for fields like password, email, country.
- Soft Deletion: This can be easily achieved with a deleted_at column
- Error tracking: This could be achieved using an interceptor in the grpc server. Or adding a hook to logrus.
- CI: If it existed it should, at least, include a test, build and deploy steps.
- Indexing: If a combination of search filters is used often, we should create a index that combination
- Docker: Docker build caching can be improved with a .dockerignore
- Implement GetUser

## How to run
- Run `docker compose -f ./deployment/docker-compose.yml up --build`

## How to run project showcase
- Follow "How to run" instructions
- `go mod download`
- `go run example.go`

## How to run tests
- Follow "How to run" instructions - Repository tests need a database connection
- `go test ./...`

## Features
- [X] Health check
- [X] Create user
- [X] Update user
- [X] Delete user
- [X] List users
- [X] Publish user events (create, update, delete)

## Layers
- Server - Responsible for handling protocol related logic (in this case grpc)
- Service - Responsible for general business logic
- Repository - Responsible for commuticating with a data store
- Provider - Responsible for bridging external services

## Decisions
- HealthCheck: https://github.com/grpc/grpc/blob/master/doc/health-checking.md
- ULID: https://github.com/oklog/ulid. ULIDs are lexicographically sortable. This is helpful for performance and consistency reasons (the database can index id more easily).
- ORM: https://github.com/go-gorm/gorm
- Validations: https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/validator
- Logging: https://github.com/rs/zerolog
- Notification: User changes are being sent to a rabbitmq topic exchange. Other services that are interested in this data can bind a queue to this exchange.
- Dependency injection: The project is simple, so dependencies can be passed in by argument. But in a bigger project we would want to use a dependency injection solution.
- Safety: Only internal services should communicate with this service. Otherwise, the connections should use TLS or mTLS. Wouldn't hurt to also have this, even if communication is only internal.
- Migrations: They were not used since there is only one table, but in a real project they should be. I would use something like https://github.com/pressly/goose to manage migrations.
- Deployment: No deployment solution was implemented. The Dockerfile could be used. This would depend on the chosen infrastructure solution but I would give Kubernetes yamls as an example.
