# Docker Compose 部署指南

本文档介绍如何使用 Docker Compose 部署 Rabbit 消息服务。

## 前置要求

- Docker 20.10+
- Docker Compose 2.0+

## 快速开始

### 1. 构建镜像

```bash
docker build -t rabbit-local:latest .
```

### 2. 准备配置文件

在 `deploy/server/docker` 目录下创建 `datasource` 目录（如果使用文件存储模式）：

```bash
cd deploy/server/docker
mkdir -p datasource
```

### 3. 配置环境变量（可选）

创建 `.env` 文件来自定义配置：

```bash
cat > .env << EOF
# 端口配置
RABBIT_HTTP_PORT=8080
RABBIT_GRPC_PORT=9090

# 环境配置
MOON_RABBIT_ENVIRONMENT=PROD
MOON_RABBIT_ENABLE_SWAGGER=true
MOON_RABBIT_ENABLE_METRICS=true

# JWT 配置
MOON_RABBIT_JWT_SECRET=your-secret-key
MOON_RABBIT_JWT_EXPIRE=600s
MOON_RABBIT_JWT_ISSUER=rabbit

# 数据库配置（如果使用数据库存储模式）
MOON_RABBIT_USE_DATABASE=false
MOON_RABBIT_MAIN_HOST=mysql
MOON_RABBIT_MAIN_DATABASE=rabbit
MOON_RABBIT_MAIN_USERNAME=root
MOON_RABBIT_MAIN_PASSWORD=your_password

# Swagger 认证
MOON_RABBIT_SWAGGER_BASIC_AUTH_ENABLED=true
MOON_RABBIT_SWAGGER_BASIC_AUTH_USERNAME=moon.rabbit
MOON_RABBIT_SWAGGER_BASIC_AUTH_PASSWORD=rabbit.swagger

# Metrics 认证
MOON_RABBIT_METRICS_BASIC_AUTH_ENABLED=true
MOON_RABBIT_METRICS_BASIC_AUTH_USERNAME=moon.rabbit
MOON_RABBIT_METRICS_BASIC_AUTH_PASSWORD=rabbit.metrics
EOF
```

### 4. 启动服务

```bash
docker-compose -f deploy/server/docker/docker-compose.yml up -d
```

## 服务管理

### 启动服务

```bash
docker-compose -f deploy/server/docker/docker-compose.yml up -d
```

### 停止服务

```bash
docker-compose -f deploy/server/docker/docker-compose.yml stop
```

### 重启服务

```bash
docker-compose -f deploy/server/docker/docker-compose.yml restart
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose -f deploy/server/docker/docker-compose.yml logs

# 实时查看日志
docker-compose -f deploy/server/docker/docker-compose.yml logs -f

# 查看特定服务日志
docker-compose -f deploy/server/docker/docker-compose.yml logs rabbit
```

### 查看服务状态

```bash
docker-compose -f deploy/server/docker/docker-compose.yml ps
```

### 停止并删除服务

```bash
docker-compose -f deploy/server/docker/docker-compose.yml down
```

### 停止并删除服务及数据卷

```bash
docker-compose -f deploy/server/docker/docker-compose.yml down -v
```

## 配置说明

### docker-compose.yml 结构

```yaml
services:
  rabbit:
    image: rabbit-local:latest
    container_name: rabbit
    ports:
      - "8080:8080"  # HTTP 端口
      - "9090:9090"  # gRPC 端口
    environment:
      # 环境变量配置
    volumes:
      - ./datasource:/moon/datasource:ro  # 数据源目录（只读）
    networks:
      - rabbit_network
    command: ["run", "all"]  # 启动命令
```

### 端口配置

- **HTTP 端口**：默认 `8080`，可通过 `RABBIT_HTTP_PORT` 环境变量修改
- **gRPC 端口**：默认 `9090`，可通过 `RABBIT_GRPC_PORT` 环境变量修改

### 卷挂载

- `./datasource:/moon/datasource:ro` - 数据源文件目录（文件存储模式时使用，只读挂载）

### 网络配置

服务运行在独立的 Docker 网络 `rabbit_network` 中，便于与其他服务通信。

### 启动命令

可以通过修改 `command` 字段来改变启动行为：

```yaml
# 默认启动所有服务
command: ["run", "all"]

# 仅启动 HTTP 服务
command: ["run", "http"]

# 仅启动 gRPC 服务
command: ["run", "grpc"]

# 仅启动 Job 服务
command: ["run", "job"]

# 使用配置文件
command: ["run", "all", "-c", "/moon/config/server.yaml"]

# 使用数据库存储模式
command: ["run", "all", "--use-database", "true", "--main-host", "mysql"]
```

## 运行模式

### 文件存储模式（默认）

使用默认配置即可，数据会存储在 `./datasource` 目录中。

### 数据库存储模式

需要修改 `docker-compose.yml` 或使用环境变量：

**方式一：修改 docker-compose.yml**

```yaml
services:
  rabbit:
    environment:
      - MOON_RABBIT_USE_DATABASE=true
      - MOON_RABBIT_MAIN_HOST=mysql
      - MOON_RABBIT_MAIN_DATABASE=rabbit
      - MOON_RABBIT_MAIN_USERNAME=root
      - MOON_RABBIT_MAIN_PASSWORD=your_password
    command: ["run", "all", "--use-database", "true"]
```

**方式二：使用 .env 文件**

在 `.env` 文件中设置相应的环境变量。

## 与 MySQL 服务集成

如果需要连接 MySQL 数据库，可以在 `docker-compose.yml` 中添加 MySQL 服务：

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: rabbit_mysql
    environment:
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: rabbit
    volumes:
      - mysql_data:/var/lib/mysql
    networks:
      - rabbit_network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  rabbit:
    image: rabbit-local:latest
    container_name: rabbit
    depends_on:
      mysql:
        condition: service_healthy
    ports:
      - "${RABBIT_HTTP_PORT:-8080}:8080"
      - "${RABBIT_GRPC_PORT:-9090}:9090"
    environment:
      - MOON_RABBIT_USE_DATABASE=true
      - MOON_RABBIT_MAIN_HOST=mysql
      - MOON_RABBIT_MAIN_DATABASE=rabbit
      - MOON_RABBIT_MAIN_USERNAME=root
      - MOON_RABBIT_MAIN_PASSWORD=your_password
    networks:
      - rabbit_network
    command: ["run", "all", "--use-database", "true"]

networks:
  rabbit_network:
    driver: bridge

volumes:
  mysql_data:
```

## 健康检查

容器内置了健康检查功能，Docker Compose 会自动监控服务健康状态：

```bash
# 查看健康状态
docker-compose -f deploy/server/docker/docker-compose.yml ps
```

健康检查配置：
- **检查命令**：`/usr/local/bin/rabbit version`
- **检查间隔**：30 秒
- **超时时间**：10 秒
- **重试次数**：3 次
- **启动等待期**：40 秒

## 故障排查

### 查看服务日志

```bash
docker-compose -f deploy/server/docker/docker-compose.yml logs rabbit
```

### 检查服务状态

```bash
docker-compose -f deploy/server/docker/docker-compose.yml ps
```

### 进入容器调试

```bash
docker-compose -f deploy/server/docker/docker-compose.yml exec rabbit sh
```

### 测试服务

```bash
# 测试 HTTP 健康检查
curl http://localhost:8080/health

# 测试版本信息
docker-compose -f deploy/server/docker/docker-compose.yml exec rabbit /usr/local/bin/rabbit version
```

### 常见问题

1. **端口被占用**
   - 检查端口是否被其他服务占用：`lsof -i :8080`
   - 修改 `.env` 文件中的端口配置

2. **容器无法启动**
   - 查看日志：`docker-compose logs rabbit`
   - 检查配置文件格式是否正确

3. **数据库连接失败**
   - 确认 MySQL 服务已启动
   - 检查数据库连接参数是否正确
   - 确认网络连接正常

## 生产环境建议

1. **使用环境变量文件**：将敏感信息存储在 `.env` 文件中，并确保不被提交到版本控制
2. **配置资源限制**：在 `docker-compose.yml` 中添加资源限制
3. **使用数据库存储**：生产环境建议使用数据库存储模式
4. **配置日志收集**：配置日志驱动和日志轮转
5. **网络安全**：使用内部网络，避免暴露不必要的端口

### 资源限制示例

```yaml
services:
  rabbit:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### 日志配置示例

```yaml
services:
  rabbit:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 更多信息

- 完整的配置选项请参考主 README 文档
- Docker Compose 官方文档：https://docs.docker.com/compose/
- 环境变量列表请参考主 README 中的环境变量部分
