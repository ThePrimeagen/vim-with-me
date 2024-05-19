FROM golang:latest

WORKDIR /app
COPY .env /app/.env
COPY ./pkg/v2 /app/pkg/v2
RUN go build /app/pkg/v2/relay/cmd/server/main.go -o relay
CMD ["./relay"]
