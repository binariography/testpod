# BUILD STAGE
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY . .
# Explicitly name the output path
RUN go build -o /bin/testpod .

# FINAL STAGE
FROM alpine:latest
WORKDIR /root/
# Copy from the explicit path in the builder stage
COPY --from=builder /bin/testpod /usr/local/bin/testpod

# Set the entrypoint to the absolute path
ENTRYPOINT ["/usr/local/bin/testpod"]
