FROM golang:1.25.6 AS builder
WORKDIR /app

# Copy dependency files and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .
# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

EXPOSE 8080

# Command to run the application
CMD ["./main"]