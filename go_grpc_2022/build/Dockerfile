FROM golang:1.19.2-alpine3.16 AS builder

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -v -o app ./cmd/app

FROM golang:1.19.2-alpine3.16

RUN addgroup -S user_grp && adduser -S -h /home/user -G user_grp user

RUN wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.13/grpc_health_probe-linux-amd64 && \
  chmod +x /bin/grpc_health_probe

USER user

WORKDIR /home/user/app

ENV ENV="production"

COPY --from=builder /build/app /bin

EXPOSE 3000

ENTRYPOINT ["/bin/app"]