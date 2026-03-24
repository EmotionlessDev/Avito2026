FROM golang:1.26.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o booking-api ./cmd/api

FROM alpine:3.19

WORKDIR /root/
COPY --from=builder /app/booking-api .

EXPOSE 8080
CMD ["./booking-api"]
