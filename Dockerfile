FROM golang:1.23.4-alpine

# Install git and air for hot reloading
RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod tidy 

# Copy the rest of the application
COPY . .

# Expose port
EXPOSE 8080

# Use air for development with hot reloading
CMD ["air", "-c", ".air.toml"]