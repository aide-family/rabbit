# Docker 部署指南

本文档介绍如何使用 Docker 部署 Rabbit 消息服务。

## 前置要求

- Docker 20.10+
- 已构建的 Rabbit 镜像（参考主 README 中的镜像构建部分）

## 快速开始

### 1. 构建镜像

```bash
docker build -t rabbit-local:latest .
```

### 2. 准备配置文件

在运行容器之前，需要准备配置文件。你可以：

**方式一：使用默认配置**
- 容器会使用内置的默认配置

**方式二：挂载配置文件目录**
- 创建本地配置文件目录：`mkdir -p ./config ./datasource`
- 生成配置文件：`rabbit config -p ./config -N server.yaml`
- 编辑配置文件以满足你的需求

### 3. 运行容器

#### 基础运行

```bash
docker run -d \
  --name rabbit \
  -p 8080:8080 \
  -p 9090:9090 \
  --restart=always \
  rabbit-local:latest run all
```

#### 使用配置文件

```bash
docker run -d \
  --name rabbit \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/config:/moon/config \
  -v $(pwd)/datasource:/moon/datasource \
  --restart=always \
  rabbit-local:latest run all -c /moon/config/server.yaml
```

#### 使用环境变量

```bash
docker run -d \
  --name rabbit \
  -p 8080:8080 \
  -p 9090:9090 \
  -e MOON_RABBIT_ENVIRONMENT=PROD \
  -e MOON_RABBIT_HTTP_ADDRESS=0.0.0.0:8080 \
  -e MOON_RABBIT_GRPC_ADDRESS=0.0.0.0:9090 \
  -e MOON_RABBIT_USE_DATABASE=true \
  -e MOON_RABBIT_MAIN_HOST=mysql \
  -e MOON_RABBIT_MAIN_DATABASE=rabbit \
  -e MOON_RABBIT_MAIN_USERNAME=root \
  -e MOON_RABBIT_MAIN_PASSWORD=your_password \
  --restart=always \
  rabbit-local:latest run all
```

## 配置说明

### 端口映射

- `8080:8080` - HTTP API 端口
- `9090:9090` - gRPC 端口

### 卷挂载

- `./config:/moon/config` - 配置文件目录（可选）
- `./datasource:/moon/datasource` - 数据源文件目录（文件存储模式时使用）

### 环境变量

所有配置都可以通过环境变量设置，环境变量遵循 `MOON_RABBIT_*` 模式。常用环境变量：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `MOON_RABBIT_ENVIRONMENT` | 环境：DEV, TEST, PREVIEW, PROD | `PROD` |
| `MOON_RABBIT_HTTP_ADDRESS` | HTTP 服务器地址 | `0.0.0.0:8080` |
| `MOON_RABBIT_GRPC_ADDRESS` | gRPC 服务器地址 | `0.0.0.0:9090` |
| `MOON_RABBIT_USE_DATABASE` | 启用数据库存储模式 | `false` |
| `MOON_RABBIT_MAIN_HOST` | MySQL 主机地址 | `localhost` |
| `MOON_RABBIT_MAIN_DATABASE` | 数据库名称 | `rabbit` |
| `MOON_RABBIT_MAIN_USERNAME` | MySQL 用户名 | `root` |
| `MOON_RABBIT_MAIN_PASSWORD` | MySQL 密码 | `123456` |

更多环境变量请参考主 README 文档。

## 运行模式

### 文件存储模式

```bash
docker run -d \
  --name rabbit \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/datasource:/moon/datasource \
  --restart=always \
  rabbit-local:latest run all \
    --datasource-paths /moon/datasource
```

### 数据库存储模式

```bash
docker run -d \
  --name rabbit \
  -p 8080:8080 \
  -p 9090:9090 \
  --restart=always \
  rabbit-local:latest run all \
    --use-database true \
    --main-host mysql \
    --main-database rabbit \
    --main-username root \
    --main-password your_password
```

## 常用命令

### 查看日志

```bash
docker logs rabbit
docker logs -f rabbit  # 实时查看日志
```

### 停止容器

```bash
docker stop rabbit
```

### 启动容器

```bash
docker start rabbit
```

### 重启容器

```bash
docker restart rabbit
```

### 删除容器

```bash
docker stop rabbit
docker rm rabbit
```

### 进入容器

```bash
docker exec -it rabbit sh
```

### 查看容器状态

```bash
docker ps | grep rabbit
docker inspect rabbit
```

## 健康检查

容器内置了健康检查，可以通过以下命令查看：

```bash
docker inspect --format='{{.State.Health.Status}}' rabbit
```

## 故障排查

### 查看容器日志

```bash
docker logs rabbit
```

### 检查端口占用

```bash
netstat -tuln | grep -E '8080|9090'
# 或
lsof -i :8080
lsof -i :9090
```

### 检查配置文件

```bash
docker exec rabbit cat /moon/config/server.yaml
```

### 测试服务

```bash
# 测试 HTTP 服务
curl http://localhost:8080/health

# 测试版本信息
docker exec rabbit /usr/local/bin/rabbit version
```

## 生产环境建议

1. **使用数据库存储模式**：生产环境建议使用数据库存储模式而非文件存储模式
2. **配置持久化**：确保配置文件和数据库连接信息正确配置
3. **资源限制**：为容器设置适当的资源限制
4. **日志管理**：配置日志收集和监控
5. **安全配置**：修改默认的 JWT 密钥和认证密码
6. **网络隔离**：使用 Docker 网络进行服务隔离

### 资源限制示例

```bash
docker run -d \
  --name rabbit \
  --memory="512m" \
  --cpus="1.0" \
  -p 8080:8080 \
  -p 9090:9090 \
  --restart=always \
  rabbit-local:latest run all
```

## 与数据库服务连接

如果 Rabbit 需要连接到独立的 MySQL 容器：

```bash
# 创建网络
docker network create rabbit_network

# 运行 MySQL
docker run -d \
  --name mysql \
  --network rabbit_network \
  -e MYSQL_ROOT_PASSWORD=your_password \
  -e MYSQL_DATABASE=rabbit \
  mysql:8.0

# 运行 Rabbit
docker run -d \
  --name rabbit \
  --network rabbit_network \
  -p 8080:8080 \
  -p 9090:9090 \
  -e MOON_RABBIT_USE_DATABASE=true \
  -e MOON_RABBIT_MAIN_HOST=mysql \
  -e MOON_RABBIT_MAIN_DATABASE=rabbit \
  -e MOON_RABBIT_MAIN_USERNAME=root \
  -e MOON_RABBIT_MAIN_PASSWORD=your_password \
  --restart=always \
  rabbit-local:latest run all
```

## 更多信息

- 完整的环境变量列表请参考主 README 文档
- 配置文件格式请参考 `config/server.yaml`
- 命令行参数请使用 `docker exec rabbit /usr/local/bin/rabbit --help` 查看
