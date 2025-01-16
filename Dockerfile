FROM golang:1.23.4 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/api ./cmd/api

FROM debian:bookworm-slim

# Set the Current Working Directory inside the container
WORKDIR /app


# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/api .
COPY --from=builder /app/cmd/api/config.yaml .

RUN useradd -u 4201 -ms /bin/bash gorunner

USER gorunner

# Command to run the executable
CMD ["./api"]