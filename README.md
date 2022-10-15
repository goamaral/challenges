# ESL challenge

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
- [X] Create user (TODO: return correct grpc status on duplicate keys)
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
- HealthCheck: Used https://github.com/grpc/grpc/blob/master/doc/health-checking.md
- ULID: Used https://github.com/oklog/ulid. ULIds are lexicographically sortable. This is helpful for performance and consistency reasons (the database can index id more easily).
- ORM: Used https://github.com/go-gorm/gorm
- Validations: Used https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/validator
- Notification: User changes are being sent to a rabbitmq topic exchange. Other services that are interested in this data can bind a queue to this exchange.

# Improvements
- Pagination: We could have added a TotalCount in the response
- Soft Deletion: This can be easily achieved with a deleted_at column
- Validations: Only basic validations were implemented but could be improved for fields like password, email, country.
- Dependency injection: The project is simple, so dependencies can be passed in by argument. But in a bigger project we would want to use a dependency injection solution.
- Safety: If only internal services communicate with this service there is no major problem. Otherwise, the connections should use TLS or mTLS. Wouldn't hurt also have this even communication is only internal.
- Migrations: They were not used since there is only one table, but in a real project they should be. I would use something like https://github.com/pressly/goose to manage migrations.
- Error tracking: This could be achieved using an interceptor in the grpc server. Or adding a hook to logrus.
- CI: If it existed it should, at least, include a test, build and deploy steps.
- Deployment: The Dockerfile could be used. This would depend on the chosen infrastructure solution but I would give Kubernetes yamls as an example.
- Indexing: If a combination of search filters is used often, we should create a index that combination
- Docker: Docker build caching can be improved with a .dockerignore