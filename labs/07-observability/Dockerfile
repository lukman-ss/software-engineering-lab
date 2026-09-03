FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/demo ./labs/07-observability/cmd/demo

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/demo .

EXPOSE 8087

CMD ["./demo"]