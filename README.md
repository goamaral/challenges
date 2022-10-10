# ESL challenge
TODO

## How to run
- Install docker
- Run `docker compose -f ./deployment/docker-compose.yml up`

## Features
- [ ] Create user
- [ ] Update user
- [ ] Delete user
- [ ] List users
- [ ] Publish user events (create, update, delete)

## Decisions and improvements
- No migrations where used, Since there is only one table. But in a real project they should be used.
I would use something like https://github.com/pressly/goose to manage migrations.
- No error tracking solution (besides logs) was implemented. This could be achieved using an interceptor in the grpc server.
- No type of CI pipeline file was created. If it existed it should at least include a test, build and deploy steps.
- No deployment files were created but the Dockerfile could be used. This depends on the chosen infrastructure solution but I would give Kubernetes yamls as an example.