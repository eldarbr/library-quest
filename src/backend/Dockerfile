# ---- Stage 1: Build ----
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o quest-api-server ./cmd/quest-api-server

# ---- Stage 2: Final Image ----
FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/quest-api-server /app/quest-api-server

USER appuser
WORKDIR /app

EXPOSE 8080

ENTRYPOINT [ "/app/quest-api-server", "-f", "words.txt", "--port", "8080", "--host", "0.0.0.0" ]
