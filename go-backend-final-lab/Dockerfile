
# ===== 阶段 1：编译 Go 程序 =====
FROM golang:1.22-alpine AS builder

WORKDIR /src

# 先复制依赖清单
COPY go.mod go.sum ./

# 下载 Go 模块
RUN go mod download

# 再复制项目源码
COPY . .

# 编译为 Linux 可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api


# ===== 阶段 2：运行程序 =====
FROM alpine:3.20

WORKDIR /app

# 只复制编译产物和网页资源
COPY --from=builder /out/api ./api
COPY --from=builder /src/web ./web

# 声明容器内服务端口
EXPOSE 8080

# 容器启动时执行 Go 程序
CMD ["./api"]
