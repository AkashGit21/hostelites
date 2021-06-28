# Build Stage
FROM golang:1.16.5-alpine3.13 AS builder

WORKDIR /app

COPY . $GOPATH/src/github.com/AkashGit21/hostelites

RUN go get -u github.com/AkashGit21/hostelites && mv $GOPATH/bin/hostelites .

# Run stage
FROM alpine:3.13

COPY --from=builder /app /app

ENTRYPOINT ["/app/hostelites"]

EXPOSE 8080
