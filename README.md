# ESL challenge
TODO

## How to run
- Install docker
- Run `docker compose -f ./deployment/docker-compose.yml up`

## Features
- [X] Health check
- [X] Create user (TODO: return correct grpc status on duplicate keys)
- [X] Update user
- [ ] Delete user
- [ ] List users
- [ ] Publish user events (create, update, delete)

## Decisions
- HealthCheck: Used https://github.com/grpc/grpc/blob/master/doc/health-checking.md
- ULID: Used https://github.com/oklog/ulid. ULIds are lexicographically sortable. This is helpful to avoid using OFFSET+LIMIT in pagination, for performance and consistency reasons.
- ORM: Used https://github.com/go-gorm/gorm
- Validations: Used https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/validator

# Improvements
- Validations: Only basic validations were implemented but could be improved for fields like password, email, country.
- Dependency injection: The project is simple, so dependencies can be passed in by argument. But in a bigger project we would want to use a dependency injection solution.
- Safety: If only internal services communicate with this service there is no major problem. Otherwise, the connections should use TLS or mTLS. Wouldn't hurt also have this even communication is only internal.
- Migrations: They were not used since there is only one table, but in a real project they should be. I would use something like https://github.com/pressly/goose to manage migrations.
- Error tracking: This could be achieved using an interceptor in the grpc server. Or adding a hook to logrus.
- CI: If it existed it should, at least, include a test, build and deploy steps.
- Deployment: The Dockerfile could be used. This would depend on the chosen infrastructure solution but I would give Kubernetes yamls as an example.