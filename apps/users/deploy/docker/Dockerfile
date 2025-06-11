# Stage 1: build-env (Build the Go application)
FROM golang:1.24-alpine AS build-env

# ติดตั้งแพ็คเกจที่จำเป็นสำหรับการคอมไพล์ (ถ้ามี เช่น git, CGO_ENABLED=0 ไม่ต้องใช้ gcc)
# ถ้าคุณใช้ CGO และต้องการ build ที่ซับซ้อน อาจจะต้องเพิ่ม build-base หรือ gcc
RUN apk update && apk add --no-cache git

# ตั้งค่า Working Directory สำหรับ Go Modules และ Go Workspace
WORKDIR /app

# Copy go.work และ go.mod/go.sum ของแต่ละโมดูล
# เพื่อให้ Go Modules สามารถ resolve dependencies ได้อย่างถูกต้อง
COPY go.work .
# Copy go.mod และ go.sum ของ main app
COPY apps/users/go.mod apps/users/go.mod
# Copy go.mod และ go.sum ของ libs (ถ้ามี)
# หาก libs/utils และ libs/core มี go.mod/go.sum ของตัวเอง ให้ copy มาด้วย
# COPY libs/utils/go.mod libs/utils/go.mod
# COPY libs/core/go.mod libs/core/go.mod

# ดาวน์โหลด dependencies โดยใช้ go work sync (สำหรับ Go 1.18+)
# หรือ go mod download (ถ้าไม่ใช้ go work sync)
# นี่จะดึง dependencies มาไว้ใน cache ก่อนที่จะ copy source code ทั้งหมด
RUN go work sync

# Copy source code ของแต่ละโมดูลที่เกี่ยวข้อง
COPY apps/users apps/users
# COPY libs/utils libs/utils
# COPY libs/core libs/core

# Build the Go application
# CGO_ENABLED=0 เพื่อสร้าง Static Binary ที่ไม่ต้องพึ่งพา C libraries ใน Runtime (ทำให้ image เล็กลง)
# -o /usr/local/bin/users-service คือ output binary path
# ./apps/users/cmd/main.go คือ entry point ของ main app
WORKDIR /app/apps/users/cmd
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /usr/local/bin/users-service ./main.go

# Stage 2: production-ready (Create a minimal image)
FROM alpine:latest

# กำหนด User ID/Group ID เพื่อความปลอดภัย (non-root user)
# สร้าง user และ group ที่จะใช้รันแอป
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy the compiled binary from the build-env stage
COPY --from=build-env /usr/local/bin/users-service /usr/local/bin/users-service

# Expose the gRPC port (matching your const port = ":50051")
EXPOSE 50051

# Define the command to run your application
CMD ["/usr/local/bin/users-service"]