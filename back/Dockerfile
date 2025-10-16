# --- Build stage ---
FROM golang:1.25.3-alpine AS build
WORKDIR /app

# ติดตั้ง git (บาง package go อาจต้องใช้)
RUN apk add --no-cache git

# Copy mod files และดาวน์โหลด dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy source code ทั้งหมด
COPY . .

# Build binary (แบบ static และเล็ก)
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o server ./cmd/server

# --- Runtime stage ---
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

# ติดตั้ง pandoc + ca-certificates (จำเป็นสำหรับ TLS และการแปลงเอกสาร)
RUN apt-get update && \
    apt-get install -y --no-install-recommends pandoc ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# คัดลอก binary ที่ build เสร็จจาก stage ก่อนหน้า
COPY --from=build /app/server /app/server

# สร้างโฟลเดอร์ tmp สำหรับเก็บไฟล์ .docx ชั่วคราว
RUN mkdir -p /app/tmp && chown -R 65532:65532 /app

# รันด้วย user non-root
USER 65532:65532

# ตั้งค่า environment variable และ port
ENV PORT=2000
EXPOSE 2000

ENTRYPOINT ["/app/server"]