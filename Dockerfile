# --- STEP 1: Build Stage ---
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# --- STEP 2: Production Target ---
FROM scratch AS release
COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]

# --- STEP 3: Debug/Test Target ---
FROM alpine:latest AS debug
RUN apk add --no-cache curl iputils jq
COPY --from=builder /main /main
ENTRYPOINT ["/main"]



# Use the official Go image as a base
FROM golang:1.25.4-alpine3.22

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o testpod .

# Command to run the application
CMD ["./testpod"]
