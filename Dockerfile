FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1

RUN go build -o server main.go

FROM alpine:latest

RUN apk add --no-cache libsqlite3

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 3000

CMD ["./server"]
