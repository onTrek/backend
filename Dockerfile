FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
ENTRYPOINT ["./entrypoint.sh"]
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
COPY ./db/migrations ./migrations

COPY ./entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

RUN mkdir -p maps gpxs avatars /root/db

EXPOSE 3000

ENTRYPOINT ["./entrypoint.sh"]

CMD ["./server"]
