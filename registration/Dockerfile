FROM golang:1.23-alpine

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Expose the port your app runs on
EXPOSE 8080

# Use Air to run the application
CMD ["air", "-c", "./registration/.air.toml"]