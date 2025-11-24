# Rabbit Docker Compose 使用指南

## 快速开始

### 1. 使用文件配置模式（默认，无需数据库）

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f rabbit

# 停止服务
docker-compose down
```

### 2. 使用数据库模式

```bash
# 1. 创建 .env 文件并设置
echo "MOON_RABBIT_USE_DATABASE=true" >> .env

# 2. 启动服务（包括 MySQL）
docker-compose up -d

# 3. 等待 MySQL 就绪后，Rabbit 会自动连接
docker-compose logs -f rabbit
```

## 自定义启动命令

### 方式一：使用 docker-compose.override.yml（推荐）

```bash
# 1. 复制示例文件
cp docker-compose.override.yml.example docker-compose.override.yml

# 2. 编辑 docker-compose.override.yml，修改 command 字段
# 例如：使用数据库模式启动
command: ["run", "--use-database", "true", "--main-host", "mysql"]

# 3. 启动服务
docker-compose up -d
```

### 方式二：使用命令行参数

```bash
# 使用数据库模式启动
docker-compose run --rm rabbit run --use-database true --main-host mysql

# 查看版本信息
docker-compose run --rm rabbit version

# 生成配置文件
docker-compose run --rm rabbit config --path /moon/config

# 发送测试邮件
docker-compose run --rm rabbit send email --uid 1 --to test@example.com --subject "Test" --body "Test email"
```

### 方式三：修改 docker-compose.yml

直接编辑 `docker-compose.yml` 文件中的 `command` 字段：

```yaml
services:
  rabbit:
    command: ["run", "--use-database", "true", "--main-host", "mysql"]
```

## 常用命令示例

### 启动服务
```bash
# 默认启动（文件配置模式）
docker-compose up -d

# 使用数据库模式启动
docker-compose up -d --env-file .env
```

### 查看服务状态
```bash
# 查看所有服务状态
docker-compose ps

# 查看 Rabbit 日志
docker-compose logs -f rabbit
```

### 执行 Rabbit 命令
```bash
# 查看版本
docker-compose exec rabbit rabbit version

# 生成配置文件
docker-compose exec rabbit rabbit config --path /moon/config

# 发送邮件
docker-compose exec rabbit rabbit send email --uid 1 --to test@example.com --subject "Test" --body "Test"
```

### 停止和清理
```bash
# 停止服务
docker-compose stop

# 停止并删除容器
docker-compose down

# 停止并删除容器和卷（会删除数据）
docker-compose down -v
```

## 配置说明

### 环境变量配置

所有配置项都可以通过环境变量设置，详见 `.env.example` 文件。

### 配置文件挂载

- `./datasource` → `/moon/datasource`：文件配置模式的配置文件目录
- `./config` → `/moon/config`：自定义配置文件目录（可选）

### 数据持久化

- `rabbit_data`：Rabbit 数据卷

## 访问服务

启动成功后，可以通过以下地址访问：

- HTTP API: http://localhost:8080
- gRPC API: localhost:9090
- Swagger UI: http://localhost:8080/swagger-ui/
- Metrics: http://localhost:8080/metrics

## 故障排查

### 查看日志
```bash
# Rabbit 日志
docker-compose logs rabbit
```

### 检查服务健康状态
```bash
# 检查 Rabbit 健康状态
docker-compose exec rabbit rabbit version
```

### 重启服务
```bash
# 重启 Rabbit
docker-compose restart rabbit
```

