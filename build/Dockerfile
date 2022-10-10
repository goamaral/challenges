FROM golang:1.19.2-alpine3.16 AS builder

WORKDIR /build

# RUN             apk add --no-cache git gcc musl-dev make

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -v -o app ./cmd/app

FROM golang:1.19.2-alpine3.16

RUN addgroup -S user_grp && adduser -S -h /home/user -G user_grp user

USER user

WORKDIR /home/user/app

COPY --from=builder /build/app /bin

EXPOSE 3000

ENTRYPOINT ["/bin/app"]