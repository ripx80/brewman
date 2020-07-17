#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk update && \
    apk upgrade && \
    apk add --no-cache git
RUN go get -d -v ./...
RUN CGO_ENABLED=1 GOOS=linux go build -v -ldflags '-s -w -extldflags "-static"' -o /go/bin/app main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT ./app
LABEL Name=brewman Version=1.0