FROM golang:latest

WORKDIR /app
COPY .env /app/.env
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
COPY ./cmd/empty/main.go /app/cmd/empty/main.go
RUN go build -o main /app/cmd/empty/main.go
RUN rm /app/main

COPY ./pkg/v2 /app/pkg/v2
RUN go build -o relay /app/pkg/v2/relay/cmd/server/main.go

COPY ./js /app/js
COPY ./views /app/views

CMD ["./relay"]
