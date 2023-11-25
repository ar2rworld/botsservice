# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app

# Build server
RUN CGO_ENABLED=0 GOOS=linux go build -o ./server ./app/server
# Build client
RUN CGO_ENABLED=0 GOOS=linux go build -o ./client ./app/client

CMD ["./server"]
