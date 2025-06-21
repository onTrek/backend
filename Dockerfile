FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ./docs ./docs

ENV GIN_MODE=release
ENV CGO_ENABLED=1

RUN go build -o server main.go

FROM alpine:latest

RUN apk add --no-cache sqlite

WORKDIR /root/

COPY --from=builder /app/server .

COPY --from=builder /app/docs ./docs

RUN mkdir -p maps gpxs

EXPOSE 3000

CMD ["./server"]
