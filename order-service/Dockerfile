# Use Go 1.23 (or higher)
FROM golang:1.23-alpine

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go app
RUN go build -o main ./cmd/main.go

# Run the app
CMD ["./main"]
