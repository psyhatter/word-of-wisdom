FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/client/main.go

CMD ["/app/main", "-address", "server:8080"]
