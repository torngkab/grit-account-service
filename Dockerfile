FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o accounts-service main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/accounts-service .

CMD ["./accounts-service"]