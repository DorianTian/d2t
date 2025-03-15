# 1. 基于官方Go镜像
FROM golang:1.20

# 2. 设置工作目录
WORKDIR /app

# 3. 复制 go.mod 和 go.sum 并先下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 4. 复制剩余源文件
COPY . .

# 5. 编译
RUN go build -o main .

# 6. 暴露端口（假设Go在8080端口监听）
EXPOSE 8080

# 7. 配置环境变量 (可选：也可在docker-compose.yml里指定)
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=Dorian
ENV DB_PASS=tqy4468
ENV DB_NAME=d2t_db

# 8. 启动可执行文件
CMD ["./main"]