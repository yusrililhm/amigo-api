FROM golang:1.22.4-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o main .

FROM alpine:3.20.0

WORKDIR /app

COPY --from=builder . .

CMD [ "/app/main" ]
