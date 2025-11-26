# 架构改造说明 - 服务职责分离

## 📋 概述

本次改造将 Rabbit 服务的启动方式从单一 `run` 命令拆分为 `server` 和 `job` 两个独立命令，实现了服务职责的清晰分离，同时保持向后兼容。

## 🎯 改造原因

### 原有架构的问题

在改造之前，Rabbit 服务使用单一的 `run` 命令同时启动：
- HTTP/gRPC 服务器（对外接口服务）
- EventBus（消息队列处理、后台任务）

这种架构存在以下问题：

1. **资源浪费**：对外接口服务和后台任务处理混合在一起，无法根据实际需求独立扩缩容
2. **职责不清**：所有功能耦合在一个进程中，难以区分服务边界
3. **运维困难**：无法针对不同类型的服务进行独立的监控、日志收集和故障处理
4. **扩展性差**：当需要增加更多后台任务类型时，会进一步增加服务复杂度

### 业界最佳实践

现代微服务架构通常采用职责分离的设计：
- **API Server**：处理外部请求，需要高可用、低延迟
- **Worker/Job**：处理异步任务，需要高吞吐、可扩展

## ✨ 改造后的好处

### 1. **独立扩缩容**

```yaml
# 生产环境部署示例
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rabbit-server
spec:
  replicas: 3  # 根据 API 请求量调整
  template:
    spec:
      containers:
      - name: rabbit
        command: ["/app/rabbit", "server"]
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rabbit-job
spec:
  replicas: 5  # 根据消息队列积压情况调整
  template:
    spec:
      containers:
      - name: rabbit
        command: ["/app/rabbit", "job"]
```

- **Server 服务**：根据 API 请求量独立扩缩容（通常 2-3 个实例）
- **Job 服务**：根据消息队列积压情况独立扩缩容（可以扩展到 10+ 个实例）

### 2. **资源隔离**

- **Server 服务**：主要消耗 CPU 和网络带宽，用于处理 HTTP/gRPC 请求
- **Job 服务**：主要消耗 CPU 和内存，用于处理消息队列任务

可以根据不同服务的资源需求，配置不同的资源限制：

```yaml
# Server 服务资源限制
resources:
  requests:
    cpu: "500m"
    memory: "512Mi"
  limits:
    cpu: "2000m"
    memory: "2Gi"

# Job 服务资源限制
resources:
  requests:
    cpu: "1000m"
    memory: "1Gi"
  limits:
    cpu: "4000m"
    memory: "4Gi"
```

### 3. **故障隔离**

- Server 服务故障不会影响后台任务处理
- Job 服务故障不会影响 API 接口服务
- 可以针对不同服务设置不同的健康检查和重启策略

### 4. **监控和日志分离**

- 可以针对 Server 和 Job 设置不同的监控指标
- 日志可以按服务类型分类收集和分析
- 更容易定位问题所在的服务

### 5. **开发调试便利**

- 开发时可以只启动 Server 服务，快速测试 API 接口
- 可以单独启动 Job 服务，测试消息处理逻辑
- 也可以通过 `http` / `grpc` / `server --job` 等命令按需启动不同组件

## 🚀 启动方式

### 1. Server 模式（组合启动 HTTP/gRPC/Job）

启动可以同时包含对外接口服务和后台任务处理的实例，通过配置/命令行组合开启 HTTP、gRPC、Job：

```bash
# 仅 HTTP + gRPC（默认）
./rabbit server

# 仅 job（单进程后台 Worker）
./rabbit server --job

# HTTP + job
./rabbit server --http --job

# gRPC + job
./rabbit server --grpc --job

# HTTP + gRPC + job（单进程 All-in-One）
./rabbit server --http --grpc --job
```

也可以通过配置文件或环境变量控制，示例：

```bash
MOON_RABBIT_SERVER_ENABLE_HTTP=true \
MOON_RABBIT_SERVER_ENABLE_GRPC=false \
MOON_RABBIT_SERVER_ENABLE_JOB=true \
./rabbit server
```

```bash
# 使用默认配置
./rabbit server

# 指定配置文件路径
./rabbit server --config ./config

# 指定 HTTP 和 gRPC 地址
./rabbit server \
  --http-address 0.0.0.0:8080 \
  --grpc-address 0.0.0.0:9090

# 指定环境
./rabbit server --environment PROD

# 查看帮助
./rabbit server --help
```

**适用场景**：
- 生产环境部署 API 服务
- 需要高可用的对外接口服务
- 根据 API 请求量进行水平扩展

### 2. Job 模式（仅启动后台任务处理）

启动只包含消息队列处理和后台任务的实例，不包含对外接口。

```bash
# 使用默认配置
./rabbit job

# 指定配置文件路径
./rabbit job --config ./config

# 指定环境
./rabbit job --environment PROD

# 查看帮助
./rabbit job --help
```

**适用场景**：
- 生产环境部署后台任务处理服务
- 需要高吞吐量的消息处理
- 根据消息队列积压情况进行水平扩展

### 3. HTTP / gRPC 独立模式

支持只启动单一协议的服务实例，方便按需拆分或本地调试：

```bash
# 仅 HTTP
./rabbit http --config ./config

# 仅 gRPC
./rabbit grpc --config ./config
```

## 📊 架构对比

### 改造前

```
┌─────────────────────────────────┐
│      rabbit run                  │
│  ┌──────────┐  ┌─────────────┐ │
│  │  HTTP    │  │  EventBus   │ │
│  │  gRPC    │  │  Workers    │ │
│  └──────────┘  └─────────────┘ │
└─────────────────────────────────┘
         │
   所有功能耦合
   无法独立扩展
```

### 改造后

```
┌─────────────────────┐      ┌─────────────────┐
│  rabbit server      │      │   rabbit job    │
│  ┌──────────┐       │      │  ┌─────────────┐│
│  │  HTTP    │ ◀──┐  │      │  │  EventBus   ││
│  │  gRPC    │ ◀┐ │  │      │  │  Workers    ││
│  │  Job     │ ┌┴─┘  │      │  └─────────────┘│
│  └──────────┘ │     │      └─────────────────┘
└─────────────────────┘
       ▲        ▲
       │        │
   rabbit http  rabbit grpc

Server 支持按需组合（HTTP/gRPC/Job），同时保留 `http`、`grpc`、`job` 三个独立模式，方便在不同场景下独立扩展。
```

## 🔧 技术实现

### 代码结构

```
cmd/
├── server/          # Server 命令（组合 HTTP/gRPC/Job）
│   ├── server.go    # 主逻辑
│   ├── flags.go     # 命令行参数（支持 --http/--grpc/--job）
│   ├── wire.go      # 依赖注入定义
│   ├── wire_gen.go  # 生成的依赖注入代码
│   └── client.go    # 客户端配置生成
├── job/             # Job 命令（仅 EventBus/后台任务）
│   ├── job.go       # 主逻辑
│   ├── flags.go     # 命令行参数
│   ├── wire.go      # 依赖注入定义
│   └── wire_gen.go  # 生成的依赖注入代码
├── http/            # HTTP 命令（仅 HTTP 服务）
│   ├── http.go      # 主逻辑
│   ├── flags.go     # 命令行参数
│   ├── wire.go      # 依赖注入定义
│   └── wire_gen.go  # 生成的依赖注入代码
└── grpc/            # gRPC 命令（仅 gRPC 服务）
    ├── grpc.go      # 主逻辑
    ├── flags.go     # 命令行参数
    ├── wire.go      # 依赖注入定义
    └── wire_gen.go  # 生成的依赖注入代码
```

### 核心改动

1. **`internal/server/server.go`**：
   - 新增 `ServerOptions` / `NewServerOptions()`，从配置或命令行解析 `enableHttp/enableGrpc/enableJob`
   - `RegisterService()` 根据 `ServerOptions` 按需注册 HTTP、gRPC、Job（EventBus），`Servers` 支持自动发现 HTTP server
   - `ProviderSetServer` 始终构建 HTTP/gRPC/EventBus，再由 `ServerOptions` 决定是否加入到运行实例

2. **依赖注入分离**：
   - Server 模式：通过开关来启用/禁用 HTTP、gRPC、Job（`enableHttp/enableGrpc/enableJob`）
   - Job 模式：使用 `ProviderSetJob`，只注入 EventBus 相关依赖
   - HTTP/GRPC 独立模式：同样复用 `ProviderSetServer`，但在命令内部强制只开启对应组件

## 📝 迁移指南

### 从旧版本迁移

如果你之前使用 `rabbit run` 命令，现在可以迁移到新的命令体系：

**选项 1：使用 server 命令（推荐，用于开发/测试/小规模部署）**
```bash
# 等价于原来的 run：同时启动 HTTP/gRPC/Job
./rabbit server --http --grpc --job
```

**选项 2：分离部署（推荐用于生产环境）**
```bash
# 部署 Server 服务（HTTP/gRPC）
./rabbit server --http --grpc

# 部署 Job 服务（可以多个实例）
./rabbit job
```

### 配置文件

所有命令使用相同的配置文件格式，无需修改配置文件。只需要确保配置文件路径正确：

```bash
# 所有命令都支持 --config 参数
./rabbit server --config /path/to/config
./rabbit job --config /path/to/config
./rabbit http --config /path/to/config
./rabbit grpc --config /path/to/config
```

## 🎓 最佳实践

### 生产环境部署

1. **Server 服务**：
   - 部署 2-3 个实例（根据 API 请求量）
   - 配置负载均衡器（如 Nginx、Traefik）
   - 设置健康检查端点
   - 监控 HTTP/gRPC 请求延迟和错误率

2. **Job 服务**：
   - 部署 5-10+ 个实例（根据消息队列积压情况）
   - 配置消息队列监控
   - 监控消息处理速度和失败率
   - 设置合理的 worker 数量（通过配置文件 `eventBus.workerCount`）

### 开发环境

使用 `run` 命令，方便本地开发和调试：

```bash
# 本地开发，同时启动所有服务
./rabbit run --config ./config
```

### 监控指标

**Server 服务监控**：
- HTTP 请求 QPS
- gRPC 请求 QPS
- 请求延迟（P50, P95, P99）
- 错误率

**Job 服务监控**：
- 消息处理速度
- 消息队列积压数量
- Worker 处理时间
- 消息处理失败率

## 🔍 常见问题

### Q: 是否必须同时部署 Server 和 Job？

A: 不是。根据你的需求：
- 如果只需要 API 服务，只部署 Server
- 如果只需要处理消息，只部署 Job
- 如果需要完整功能，同时部署两者

### Q: Run 命令会被废弃吗？

A: 不会。Run 命令会一直保留，用于开发、测试和小规模部署场景。

### Q: Server 和 Job 可以共享同一个数据库吗？

A: 可以。它们使用相同的配置文件，可以连接到同一个数据库。

### Q: 如何实现 Server 和 Job 之间的通信？

A: 目前通过数据库进行通信（消息日志表）。Server 接收请求后写入消息日志，Job 从消息日志中读取并处理。

