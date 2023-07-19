# Base image for building the app
FROM golang:1.20 AS builder
WORKDIR /go/src/app
COPY . .
RUN go build -o admin cmd/admin/main.go && \
    go build -o eventlistener cmd/eventlistener/main.go

# Base image for running the app
FROM golang:1.20
WORKDIR /app
COPY --from=builder /go/src/app/admin ./admin
COPY --from=builder /go/src/app/eventlistener ./eventlistener
