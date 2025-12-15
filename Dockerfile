FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN apk add git && git config --global url."https://torngkab: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@github.com/".insteadOf "https://github.com/"

RUN go env -w GOPRIVATE=github.com/torngkab/* && go mod download && go mod verify

COPY . .

COPY .env.example .env

RUN go build -o accounts-service main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/accounts-service .

CMD ["./accounts-service"]