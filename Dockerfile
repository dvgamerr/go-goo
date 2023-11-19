FROM golang:1.15 as builder

WORKDIR /project
COPY . .

# Production-ready build, without debug information specifically for linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o=goo .

FROM alpine:3.12

# Add CA certificates required for SSL connections
RUN apk add --update --no-cache ca-certificates

COPY --from=builder /project/goo /usr/local/bin/ggoo

RUN mkdir /app
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/goo"]