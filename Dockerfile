# Build
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

# Run
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY test_in.txt .
COPY test1_in.txt .

ENTRYPOINT [ "./main" ]