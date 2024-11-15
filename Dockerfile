FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN go build -o main .

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/main ./

CMD [ "./main" ]
