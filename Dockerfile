# building app
FROM golang:1.27.2 AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/thanos_queue

# building image app
FROM alpine:3.24.3
RUN adduser -S appuser
WORKDIR /home/appuser
COPY --from=builder /app/app .
USER appuser
CMD ["./app"]
