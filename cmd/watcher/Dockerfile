FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /watcher_bin ./cmd/watcher/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /watcher_bin /watcher

EXPOSE 8081

CMD ["/watcher"]