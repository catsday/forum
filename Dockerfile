FROM golang:1.23

RUN apt-get update && apt-get install -y gcc

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o forum ./cmd/main.go

EXPOSE 8080

CMD ["./forum"]
