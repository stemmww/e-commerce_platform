FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o inventory ./cmd

EXPOSE 8081
EXPOSE 50051

CMD ["./inventory"]
