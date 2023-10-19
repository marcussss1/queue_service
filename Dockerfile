FROM golang:1.21.3-alpine

WORKDIR /build

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o /api cmd/api/main.go

EXPOSE 8080

ENTRYPOINT ["/api"]
