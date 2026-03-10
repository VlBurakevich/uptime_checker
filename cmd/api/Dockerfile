FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app_bin ./cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app_bin /api

EXPOSE 8080

CMD ["/api"]