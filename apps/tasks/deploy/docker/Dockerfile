FROM golang:1.24-alpine AS build-env

# git อาจจำเป็นถ้า Go modules ของคุณมีการดึง dependency จาก private repo หรือ git
RUN apk update && apk add --no-cache git

# ตั้งค่า Working Directory หลักสำหรับ Go Workspace ใน build stage
WORKDIR /app

# Copy go.work และ go.mod/go.sum ของแต่ละโมดูลก่อน
# สิ่งนี้ช่วยให้ Docker cache layer ถูกนำกลับมาใช้ใหม่ได้ หากไม่มีการเปลี่ยนแปลงในไฟล์เหล่านี้
# การแยก COPY นี้จะช่วยให้ build เร็วขึ้นในภายหลัง
COPY go.work .
COPY apps/tasks/go.mod apps/tasks/go.mod
COPY apps/tasks/go.sum apps/tasks/go.sum
COPY apps/users/go.mod apps/users/go.mod
COPY apps/users/go.sum apps/users/go.sum

# ถ้ามี libs/modules อื่นๆ ใน go.work ที่ project นี้ใช้ ให้ copy go.mod/go.sum ของ libs นั้นๆ ด้วย
# ตัวอย่าง:
# COPY libs/utils/go.mod libs/utils/go.mod
# COPY libs/utils/go.sum libs/utils/go.sum
# COPY libs/core/go.mod libs/core/go.mod
# COPY libs/core/go.sum libs/core/go.sum


# ดาวน์โหลด dependencies โดยใช้ go work sync (สำหรับ Go 1.18+)
# การใช้ go work sync จะดีกว่าสำหรับ Go Workspace
RUN go work sync


# Copy source code ของแต่ละโมดูลที่เกี่ยวข้อง
# Copy เฉพาะสิ่งที่จำเป็นสำหรับการ build เพื่อลดขนาด context ที่ส่งไป Docker daemon
COPY apps/tasks apps/tasks
COPY apps/users apps/users
# ถ้ามี libs/modules อื่นๆ ใน go.work ที่ project นี้ใช้ ให้ copy มาด้วย
# COPY libs/utils libs/utils
# COPY libs/core libs/core


# Build the Go application
# CGO_ENABLED=0 เพื่อสร้าง Static Binary ที่ไม่ต้องพึ่งพา C libraries ใน Runtime (ทำให้ image เล็กลงและ portable)
# -ldflags="-s -w" เพื่อลดขนาด binary โดยการลบ debug info และ symbol table
# -o /usr/local/bin/tasks-service คือ output binary path ใน Container
# ./apps/tasks/cmd/api/main.go คือ entry point ของ main app
# WORKDIR ยังอยู่ที่ /app ทำให้ path "./apps/tasks/cmd/api/main.go" ทำงานได้
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /usr/local/bin/tasks-service ./apps/tasks/cmd/api/main.go

# Stage 2: production-ready (Create a minimal and secure image)
FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=build-env /usr/local/bin/tasks-service /usr/local/bin/tasks-service

RUN mkdir -p /app/config
# viper.AddConfigPath("/app/config") in apps/tasks/internal/config/config.go
COPY apps/tasks/internal/config/config.yml /app/config/config.yml

USER appuser

EXPOSE 50051

CMD ["/usr/local/bin/tasks-service"]