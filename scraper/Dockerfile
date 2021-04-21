FROM golang:1.16

WORKDIR /app

COPY go.mod .
RUN go mod download -x

COPY main.go .
RUN go build main.go
CMD ["./main", "http://server:8000"]
