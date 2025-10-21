# --- Build stage ---
FROM golang:1.25.3-alpine AS build
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o server ./cmd/server

# --- Runtime stage: pandoc/core (Alpine-based)
FROM pandoc/core:latest
# ติดตั้ง certs/timezone 
RUN apk add --no-cache \
    ca-certificates tzdata 

WORKDIR /app
COPY --from=build /app/server /app/server
RUN mkdir -p /app/tmp && chown -R 65532:65532 /app
USER 65532:65532
ENV PORT=2000
EXPOSE 2000
ENTRYPOINT ["/app/server"]