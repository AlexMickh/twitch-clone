FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /env-setter ./cmd/env-setter/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /twitch-clone ./cmd/app/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /twitch-clone .
COPY --from=builder /env-setter .
COPY ./config/prod.yml .

RUN ./env-setter --config=./prod.yml

EXPOSE 8080

CMD [ "./twitch-clone" ]
