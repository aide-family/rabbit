# 多阶段构建 - 构建阶段
FROM golang:1.25.3-alpine AS builder

# 安装必要的构建工具
RUN apk add --no-cache \
    protobuf-dev \
    protoc \
    git \
    make

WORKDIR /moon

# 复制构建文件
COPY Makefile Makefile

# 初始化环境
RUN make init

# 复制源代码
COPY . .

# 构建应用
RUN make build

# 最终运行阶段 - 使用 Alpine
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

# 创建非 root 用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /moon

# 复制二进制文件和配置
COPY --from=builder /moon/bin/ /usr/sbin/
COPY --from=builder /moon/config/server.yaml /moon/config/server.yaml

# 设置权限
RUN chown -R appuser:appgroup /moon

# 切换到非 root 用户
USER appuser

# 设置卷
VOLUME /moon/config
VOLUME /moon/.rabbit

# 暴露端口
EXPOSE 8080
EXPOSE 9090

# 运行应用
CMD ["rabbit", "run", "-c", "/moon/config/", "--rabbit-config", "/moon/.rabbit/"]
