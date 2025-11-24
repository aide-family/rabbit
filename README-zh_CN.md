# Rabbit (玉兔) 🐰

[![Go 版本](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![许可证](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Kratos](https://img.shields.io/badge/Kratos-v2-00ADD8?style=flat&logo=go)](https://github.com/go-kratos/kratos)

> 基于 Kratos 框架构建的分布式消息服务平台，提供统一的消息发送和管理能力。

## 📖 项目介绍

Rabbit (玉兔) 是一个高性能、高可用、高扩展的分布式消息服务平台。它支持多种消息通道（邮件、Webhook、短信、飞书等）的统一管理和发送，通过命名空间实现多租户隔离，支持配置文件和数据库两种存储模式，满足不同场景的部署需求。

## ✨ 核心特性

- **多通道消息发送**：统一管理邮件、Webhook、短信、飞书等多种消息通道
- **模板化发送**：支持消息模板配置，实现消息内容的动态渲染和复用
- **异步消息处理**：基于消息队列实现异步发送，提升系统吞吐量和可靠性
- **配置管理**：支持邮件服务器、Webhook 端点等通道配置的集中管理
- **多租户隔离**：通过命名空间实现不同业务或租户的配置和数据隔离
- **灵活存储**：支持配置文件和数据库两种存储模式
- **丰富的 CLI 工具**：提供完整的命令行接口，支持服务管理、消息发送、配置生成等
- **热加载**：支持配置文件热加载，无需重启服务

## 🚀 快速开始

### 前置要求

- Go 1.25+ (从源码构建)
- Docker & Docker Compose (容器化部署)
- MySQL 5.7+ (可选，用于数据库存储模式)
- etcd (可选，用于服务注册)

### 安装

#### 从源码安装

```bash
# 克隆仓库
git clone https://github.com/aide-family/rabbit.git
cd rabbit

# 初始化环境
make init

# 构建二进制文件
make build

# 运行服务
./bin/rabbit run
```

#### 使用 Docker

```bash
# 构建 Docker 镜像
docker build -t rabbit:latest .

# 运行容器
docker run -d \
  --name rabbit \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/config:/moon/config \
  -v $(pwd)/datasource:/moon/datasource \
  rabbit:latest
```

#### 使用 Docker Compose

```bash
cd deploy/server/docker
docker-compose up -d
```

### 生成配置文件

```bash
# 生成默认配置文件
rabbit config --path ./config --name server.yaml

# 或使用自定义路径
rabbit config -p ./config -n server.yaml
```

## 📦 部署

### Docker 部署

详细说明请参考 [Docker Compose 文档](deploy/server/docker/README-docker-compose.md)。

```bash
cd deploy/server/docker
docker-compose up -d
```

### Kubernetes 部署

详细说明请参考 [Kubernetes 部署指南](deploy/server/k8s/README.md)。

#### 快速部署

```bash
cd deploy/server/k8s
./deploy.sh
```

#### 使用 Kustomize

```bash
kubectl apply -k deploy/server/k8s/
```

### 手动部署

1. **构建二进制文件**：
   ```bash
   make build
   ```

2. **生成配置文件**：
   ```bash
   rabbit config -p ./config
   ```

3. **编辑配置**：
   根据环境编辑 `config/server.yaml`。

4. **运行服务**：
   ```bash
   ./bin/rabbit run -c ./config/server.yaml
   ```

## ⚙️ 配置说明

### 配置文件

默认配置文件为 `config/server.yaml`。可以使用 `--config` 或 `-c` 参数指定自定义路径。

### 环境变量

Rabbit 支持通过环境变量进行配置。所有环境变量遵循 `MOON_RABBIT_*` 模式。

#### 服务器配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `MOON_RABBIT_ENVIRONMENT` | `PROD` | 环境：DEV, TEST, PREVIEW, PROD |
| `MOON_RABBIT_NAME` | `moon.rabbit` | 服务名称 |
| `MOON_RABBIT_HTTP_ADDRESS` | `0.0.0.0:8080` | HTTP 服务器地址 |
| `MOON_RABBIT_GRPC_ADDRESS` | `0.0.0.0:9090` | gRPC 服务器地址 |
| `MOON_RABBIT_HTTP_TIMEOUT` | `10s` | HTTP 请求超时时间 |
| `MOON_RABBIT_GRPC_TIMEOUT` | `10s` | gRPC 请求超时时间 |

#### 数据库配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `MOON_RABBIT_USE_DATABASE` | `false` | 启用数据库存储模式 |
| `MOON_RABBIT_MAIN_HOST` | `localhost` | MySQL 主机地址 |
| `MOON_RABBIT_MAIN_PORT` | `3306` | MySQL 端口 |
| `MOON_RABBIT_MAIN_DATABASE` | `rabbit` | 数据库名称 |
| `MOON_RABBIT_MAIN_USERNAME` | `root` | MySQL 用户名 |
| `MOON_RABBIT_MAIN_PASSWORD` | `123456` | MySQL 密码 |
| `MOON_RABBIT_MAIN_DEBUG` | `false` | 启用数据库调试模式 |

#### JWT 配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `MOON_RABBIT_JWT_SECRET` | `xxx` | JWT 密钥 |
| `MOON_RABBIT_JWT_EXPIRE` | `600s` | JWT 过期时间 |
| `MOON_RABBIT_JWT_ISSUER` | `rabbit` | JWT 签发者 |

#### 注册中心配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `MOON_RABBIT_REGISTRY_TYPE` | `` | 注册中心类型：etcd, kubernetes |
| `MOON_RABBIT_ETCD_ENDPOINTS` | `127.0.0.1:2379` | etcd 端点（逗号分隔） |
| `MOON_RABBIT_ETCD_USERNAME` | `` | etcd 用户名 |
| `MOON_RABBIT_ETCD_PASSWORD` | `` | etcd 密码 |
| `MOON_RABBIT_KUBERNETES_NAMESPACE` | `moon` | Kubernetes 命名空间 |
| `MOON_RABBIT_KUBERNETES_KUBECONFIG` | `~/.kube/config` | Kubernetes kubeconfig 路径 |

#### 功能开关

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `MOON_RABBIT_ENABLE_CLIENT_CONFIG` | `true` | 启用客户端配置 |
| `MOON_RABBIT_ENABLE_SWAGGER` | `true` | 启用 Swagger UI |
| `MOON_RABBIT_ENABLE_METRICS` | `true` | 启用指标端点 |
| `MOON_RABBIT_CONFIG_PATHS` | `./datasource` | 配置文件路径（逗号分隔） |

### 命令行参数

#### 全局参数

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--namespace` | `-n` | `` | 服务命名空间 |
| `--rabbit-config` | | `./.rabbit/` | Rabbit 配置文件目录 |

#### Run 命令参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--config`, `-c` | `` | 配置文件路径 |
| `--environment` | `PROD` | 环境：DEV, TEST, PREVIEW, PROD |
| `--http-address` | `0.0.0.0:8080` | HTTP 服务器地址 |
| `--grpc-address` | `0.0.0.0:9090` | gRPC 服务器地址 |
| `--use-database` | `false` | 启用数据库存储模式 |
| `--config-paths` | `./datasource` | 配置文件路径 |

更多参数请使用 `rabbit run --help` 查看。

### 使用示例

```bash
# 使用自定义配置文件运行
rabbit run -c ./config/server.yaml

# 使用环境变量运行
MOON_RABBIT_HTTP_ADDRESS=0.0.0.0:8080 \
MOON_RABBIT_USE_DATABASE=true \
rabbit run

# 使用命令行参数运行
rabbit run \
  --http-address 0.0.0.0:8080 \
  --grpc-address 0.0.0.0:9090 \
  --use-database true \
  --config-paths ./datasource,./config
```

## 📚 命令说明

### 基础命令

- `rabbit config` - 生成默认配置文件
- `rabbit version` - 显示版本信息

### 消息命令

- `rabbit send email` - 发送邮件消息
- `rabbit send sms` - 发送短信消息
- `rabbit send feishu` - 发送飞书消息
- `rabbit apply` - 提交消息到队列
- `rabbit get` - 获取消息信息
- `rabbit delete` - 删除消息

### 服务命令

- `rabbit run` - 启动 Rabbit 服务

### 数据库命令

- `rabbit gorm migrate` - 迁移数据库表结构
- `rabbit gorm gen` - 生成 GORM 查询代码

详细命令说明请使用 `rabbit --help` 查看。

## 🔧 开发指南

### 前置要求

- Go 1.25+
- Make
- Protocol Buffers 编译器 (protoc)
- MySQL 8.0+ (数据库模式需要)

### 设置开发环境

```bash
# 克隆仓库
git clone https://github.com/aide-family/rabbit.git
cd rabbit

# 初始化环境
make init

# 生成所有代码
make all

# 运行测试
make test

# 开发模式运行
make dev
```

### 项目结构

```
rabbit/
├── cmd/              # 命令行接口
├── internal/         # 内部包
│   ├── biz/         # 业务逻辑
│   ├── data/        # 数据层
│   ├── server/      # 服务器实现
│   └── conf/        # 配置
├── pkg/             # 公共包
├── proto/           # Protocol Buffer 定义
├── config/          # 配置文件
├── deploy/          # 部署配置
└── Makefile         # 构建脚本
```

## 🤝 贡献指南

我们欢迎贡献！提交 PR 前请先阅读贡献指南。

### Pull Request 流程

1. **Fork 仓库**并从 `main` 分支创建你的分支
2. **创建 Issue** 讨论你的更改（如果是重大更改）
3. **进行更改**，遵循我们的代码风格指南
4. **添加测试**（新功能或 bug 修复）
5. **更新文档**（如需要）
6. **确保所有测试通过** (`make test`)
7. **提交 Pull Request**，附上清晰的描述

#### PR 标题格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型：**
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更改
- `style`: 代码风格更改（格式化等）
- `refactor`: 代码重构
- `test`: 测试添加或更改
- `chore`: 构建过程或辅助工具更改

**示例：**
```
feat(message): 添加邮件模板支持

添加对邮件模板的支持，支持动态变量替换。
模板可以在配置文件中定义，发送邮件时通过名称引用。

Closes #123
```

#### PR 检查清单

- [ ] 代码遵循项目的风格指南
- [ ] 已完成自我审查
- [ ] 为复杂代码添加了注释
- [ ] 已更新文档
- [ ] 已添加/更新测试
- [ ] 所有测试通过
- [ ] 未引入新的警告
- [ ] 更改向后兼容（或提供了迁移指南）

### Issue 报告

报告问题时，请包含：

1. **问题类型**：Bug、功能请求、问题等
2. **描述**：问题的清晰描述
3. **复现步骤**：对于 bug，提供复现步骤
4. **预期行为**：你期望发生什么
5. **实际行为**：实际发生了什么
6. **环境**：操作系统、Go 版本、Rabbit 版本
7. **配置**：相关配置（已脱敏）
8. **日志**：相关日志输出
9. **截图**：如适用

#### Issue 模板

```markdown
**问题类型**: [Bug/功能请求/问题]

**描述**:
<!-- 问题的清晰描述 -->

**复现步骤** (针对 bug):
1. 
2. 
3. 

**预期行为**:
<!-- 你期望发生什么 -->

**实际行为**:
<!-- 实际发生了什么 -->

**环境**:
- 操作系统: 
- Go 版本: 
- Rabbit 版本: 

**配置**:
```yaml
<!-- 相关配置（已脱敏） -->
```

**日志**:
```
<!-- 相关日志输出 -->
```

**其他信息**:
<!-- 任何其他相关信息 -->
```

## 📄 许可证

本项目采用 Apache License 2.0 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [Kratos](https://github.com/go-kratos/kratos) - 微服务框架
- [Cobra](https://github.com/spf13/cobra) - Go 命令行框架

## 📞 联系方式

- **仓库**: https://github.com/aide-family/rabbit
- **Issues**: https://github.com/aide-family/rabbit/issues
- **邮箱**: aidecloud@163.com

---

由 [Aide Family](https://github.com/aide-family) 用 ❤️ 制作
