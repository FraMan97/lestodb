FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o lestodb ./internal/server-tcp/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/lestodb .

EXPOSE 2001

ENV PORT=2001
ENV SHARDING_COUNT=36

CMD ["./lestodb"]