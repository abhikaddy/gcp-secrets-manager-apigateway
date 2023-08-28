FROM golang:1.21.0-alpine

# Set the working directory
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the source code
COPY main.go ./

# Build the Go binary
RUN go build -o gcpsm-apigateway .

# Start the application
CMD [ "./gcpsm-apigateway" ]
