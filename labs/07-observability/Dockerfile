FROM golang:1.25-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/demo ./labs/07-observability/cmd/demo

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/demo .

EXPOSE 8087

CMD ["./demo"]