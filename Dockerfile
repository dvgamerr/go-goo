FROM golang:1.15 as builder

WORKDIR /project
COPY . .

# Production-ready build, without debug information specifically for linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o=goog .

FROM alpine:3.12

# Add CA certificates required for SSL connections
RUN apk add --update --no-cache ca-certificates

COPY --from=builder /project/goog /usr/local/bin/goog

RUN mkdir /app
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/goog"]