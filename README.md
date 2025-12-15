# Rabbit (Jade Rabbit) 🐰

<div align="right">

[English](README.md) | [中文](README-zh_CN.md)

</div>

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Kratos](https://img.shields.io/badge/Kratos-v2-00ADD8?style=flat&logo=go)](https://github.com/go-kratos/kratos)

> A distributed messaging platform built on the Kratos framework, providing unified message delivery and management capabilities.

## 📖 Introduction

Rabbit (Jade Rabbit) is a high-performance, highly available, and highly scalable distributed messaging service platform. It supports unified management and delivery of multiple message channels (email, Webhook, SMS, Feishu, etc.), implements multi-tenant isolation through namespaces, and supports both file-based and database storage modes to meet different deployment requirements.

## ✨ Features

- **Multi-channel Messaging**: Unified management of email, Webhook, SMS, Feishu, and other message channels
- **Template-based Delivery**: Support for message template configuration with dynamic content rendering and reuse
- **Asynchronous Processing**: Queue-based asynchronous message delivery for improved throughput and reliability
- **Configuration Management**: Centralized management of channel configurations (email servers, Webhook endpoints, etc.)
- **Multi-tenant Isolation**: Namespace-based isolation of configurations and data for different businesses or tenants
- **Flexible Storage**: Support for both file-based and database storage modes
- **Rich CLI Tools**: Comprehensive command-line interface for service management, message sending, and configuration generation
- **Hot Reload**: Support for hot reloading of configurations without service restart

## 🚀 Quick Start

### Prerequisites

- Go 1.25+ (for building from source)
- Docker & Docker Compose (for containerized deployment)
- MySQL 5.7+ (optional, for database storage mode)
- etcd (optional, for service registry)

### Installation

#### From Source

```bash
# Clone the repository
git clone https://github.com/aide-family/rabbit.git
cd rabbit

# Initialize the environment
make init

# Build the binary
make build

# Run the service
./bin/rabbit run all
```

#### Using Docker

```bash
# Build the Docker image
docker build -t rabbit:latest .

# Run the container
docker run -d \
  --name rabbit \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/config:/moon/config \
  -v $(pwd)/datasource:/moon/datasource \
  rabbit:latest
```

#### Using Docker Compose

```bash
cd deploy/server/docker
docker-compose up -d
```

### Generate Configuration

```bash
# Generate default configuration file
rabbit config --path ./config --name server.yaml

# Or with custom path
rabbit config -p ./config -n server.yaml

# Force overwrite existing file
rabbit config -p ./config -n server.yaml --force

# Generate client configuration file
rabbit config -p ./config -n client.yaml --client
```

## 📦 Deployment

### Docker Deployment

See [Docker Compose Documentation](deploy/server/docker/README-docker-compose.md) for detailed instructions.

```bash
cd deploy/server/docker
docker-compose up -d
```

### Kubernetes Deployment

See [Kubernetes Deployment Guide](deploy/server/k8s/README.md) for detailed instructions.

#### Quick Deploy

```bash
cd deploy/server/k8s
./deploy.sh
```

#### Using Kustomize

```bash
kubectl apply -k deploy/server/k8s/
```

### Manual Deployment

1. **Build the binary**:
   ```bash
   make build
   ```

2. **Generate configuration**:
   ```bash
   rabbit config -p ./config
   ```

3. **Edit configuration**:
   Edit `config/server.yaml` according to your environment.

4. **Run the service**:
   ```bash
   ./bin/rabbit run all -c ./config/server.yaml
   ```

## ⚙️ Configuration

### Configuration File

The default configuration file is `config/server.yaml`. You can specify custom paths using the `--config` or `-c` flag (can be used multiple times).

**Note**: The `--use-database` and `--datasource-paths` options are mutually exclusive:
- Use `--use-database true` for database storage mode (recommended for production)
- Use `--datasource-paths` for file-based storage mode (useful for development and testing)

### Environment Variables

Rabbit supports configuration through environment variables. All environment variables follow the pattern `MOON_RABBIT_*`.

#### Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_ENVIRONMENT` | `PROD` | Environment: DEV, TEST, PREVIEW, PROD |
| `MOON_RABBIT_NAME` | `moon.rabbit` | Service name |
| `MOON_RABBIT_USE_RANDOM_ID` | `false` | Use random service ID |
| `MOON_RABBIT_METADATA_TAG` | `rabbit` | Service metadata tag |
| `MOON_RABBIT_METADATA_REPOSITORY` | `https://github.com/aide-family/rabbit` | Service metadata repository |
| `MOON_RABBIT_METADATA_AUTHOR` | `Aide Family` | Service metadata author |
| `MOON_RABBIT_METADATA_EMAIL` | `aidecloud@163.com` | Service metadata email |
| `MOON_RABBIT_HTTP_ADDRESS` | `0.0.0.0:8080` | HTTP server address |
| `MOON_RABBIT_HTTP_NETWORK` | `tcp` | HTTP server network |
| `MOON_RABBIT_HTTP_TIMEOUT` | `10s` | HTTP request timeout |
| `MOON_RABBIT_GRPC_ADDRESS` | `0.0.0.0:9090` | gRPC server address |
| `MOON_RABBIT_GRPC_NETWORK` | `tcp` | gRPC server network |
| `MOON_RABBIT_GRPC_TIMEOUT` | `10s` | gRPC request timeout |
| `MOON_RABBIT_JOB_ADDRESS` | `0.0.0.0:9091` | Job server address |
| `MOON_RABBIT_JOB_NETWORK` | `grpc` | Job server network |
| `MOON_RABBIT_JOB_TIMEOUT` | `10s` | Job request timeout |

#### Database Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_USE_DATABASE` | `false` | Enable database storage mode (mutually exclusive with MOON_RABBIT_DATASOURCE_PATHS) |
| `MOON_RABBIT_MAIN_HOST` | `localhost` | MySQL host |
| `MOON_RABBIT_MAIN_PORT` | `3306` | MySQL port |
| `MOON_RABBIT_MAIN_DATABASE` | `rabbit` | Database name |
| `MOON_RABBIT_MAIN_USERNAME` | `root` | MySQL username |
| `MOON_RABBIT_MAIN_PASSWORD` | `123456` | MySQL password |
| `MOON_RABBIT_MAIN_DEBUG` | `false` | Enable database debug mode |
| `MOON_RABBIT_MAIN_USE_SYSTEM_LOGGER` | `true` | Use system logger for database |

#### JWT Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_JWT_SECRET` | `xxx` | JWT secret key |
| `MOON_RABBIT_JWT_EXPIRE` | `600s` | JWT expiration time |
| `MOON_RABBIT_JWT_ISSUER` | `rabbit` | JWT issuer |

#### Registry Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_REGISTRY_TYPE` | `` | Registry type: etcd, kubernetes |
| `MOON_RABBIT_ETCD_ENDPOINTS` | `127.0.0.1:2379` | etcd endpoints (comma-separated) |
| `MOON_RABBIT_ETCD_USERNAME` | `` | etcd username |
| `MOON_RABBIT_ETCD_PASSWORD` | `` | etcd password |
| `MOON_RABBIT_KUBERNETES_NAMESPACE` | `moon` | Kubernetes namespace |
| `MOON_RABBIT_KUBERNETES_KUBECONFIG` | `~/.kube/config` | Kubernetes kubeconfig path |

#### Cluster Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_CLUSTER_NAME` | `moon.rabbit` | Cluster name |
| `MOON_RABBIT_CLUSTER_ENDPOINTS` | `` | Cluster endpoints |
| `MOON_RABBIT_CLUSTER_PROTOCOL` | `GRPC` | Cluster protocol: GRPC, HTTP |
| `MOON_RABBIT_CLUSTER_TIMEOUT` | `10s` | Cluster request timeout |

#### Job Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_JOB_CORE_WORKER_TOTAL` | `10` | Total number of job workers |
| `MOON_RABBIT_JOB_CORE_TIMEOUT` | `10s` | Job core timeout |
| `MOON_RABBIT_JOB_CORE_BUFFER_SIZE` | `1000` | Job core buffer size |

#### Feature Flags

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_ENABLE_CLIENT_CONFIG` | `false` | Enable client configuration |
| `MOON_RABBIT_ENABLE_SWAGGER` | `false` | Enable Swagger UI |
| `MOON_RABBIT_ENABLE_METRICS` | `false` | Enable metrics endpoint |
| `MOON_RABBIT_DATASOURCE_PATHS` | `` | Data source file paths (comma-separated, mutually exclusive with MOON_RABBIT_USE_DATABASE) |
| `MOON_RABBIT_MESSAGE_LOG_PATH` | `` | Message log file path |

#### Swagger Basic Auth

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_SWAGGER_BASIC_AUTH_ENABLED` | `true` | Enable Swagger basic authentication |
| `MOON_RABBIT_SWAGGER_BASIC_AUTH_USERNAME` | `moon.rabbit` | Swagger basic auth username |
| `MOON_RABBIT_SWAGGER_BASIC_AUTH_PASSWORD` | `rabbit.swagger` | Swagger basic auth password |

#### Metrics Basic Auth

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_METRICS_BASIC_AUTH_ENABLED` | `true` | Enable metrics basic authentication |
| `MOON_RABBIT_METRICS_BASIC_AUTH_USERNAME` | `moon.rabbit` | Metrics basic auth username |
| `MOON_RABBIT_METRICS_BASIC_AUTH_PASSWORD` | `rabbit.metrics` | Metrics basic auth password |

### Command Line Arguments

#### Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--namespace` | `-n` | `` | The namespace of the service |
| `--rabbit-config` | | `./.rabbit/` | The config file directory of the rabbit |
| `--log-format` | | `TEXT` | Log format: TEXT, JSON |
| `--log-level` | | `DEBUG` | Log level: DEBUG, INFO, WARN, ERROR |

#### Config Command Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--path`, `-p` | | `.` | Output path for the config file |
| `--name` | | `config.yaml` | Output file name |
| `--force`, `-f` | | `false` | Overwrite existing file if it exists |
| `--client` | | `false` | Generate client config file instead of server config |

#### Run Command Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config`, `-c` | `` | Configuration file path (can be used multiple times) |
| `--enable-client-config` | `false` | Enable client configuration |
| `--server-name` | `rabbit` | Server name |
| `--use-random-node-id` | `false` | Use random node ID |
| `--server-metadata` | `` | Server metadata (format: key=value, can be used multiple times) |
| `--environment` | `PROD` | Environment: DEV, TEST, PREVIEW, PROD |
| `--jwt-secret` | `xxx` | JWT secret key |
| `--jwt-expire` | `600s` | JWT expiration time |
| `--jwt-issuer` | `rabbit` | JWT issuer |
| `--main-username` | `root` | MySQL username |
| `--main-password` | `123456` | MySQL password |
| `--main-host` | `localhost` | MySQL host |
| `--main-port` | `3306` | MySQL port |
| `--main-database` | `rabbit` | Database name |
| `--main-debug` | `false` | Enable database debug mode |
| `--main-use-system-logger` | `true` | Use system logger for database |
| `--registry-type` | `` | Registry type: ETCD, KUBERNETES |
| `--etcd-endpoints` | `127.0.0.1:2379` | etcd endpoints |
| `--etcd-username` | `` | etcd username |
| `--etcd-password` | `` | etcd password |
| `--kubernetes-kubeconfig` | `~/.kube/config` | Kubernetes kubeconfig path |
| `--use-database` | `false` | Enable database storage mode (mutually exclusive with --datasource-paths) |
| `--datasource-paths` | `` | Data source file paths (comma-separated, mutually exclusive with --use-database) |
| `--message-log-path` | `` | Message log file path |

#### Run All Command Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--http-address` | `0.0.0.0:8080` | HTTP server address |
| `--http-network` | `tcp` | HTTP server network |
| `--http-timeout` | `10s` | HTTP request timeout |
| `--grpc-address` | `0.0.0.0:9090` | gRPC server address |
| `--grpc-network` | `tcp` | gRPC server network |
| `--grpc-timeout` | `10s` | gRPC request timeout |
| `--job-address` | `0.0.0.0:9091` | Job server address |
| `--job-network` | `grpc` | Job server network |
| `--job-timeout` | `10s` | Job request timeout |
| `--job-core-worker-total` | `10` | Total number of job workers |
| `--job-core-timeout` | `10s` | Job core timeout |
| `--job-core-buffer-size` | `1000` | Job core buffer size |
| `--enable-swagger` | `false` | Enable Swagger UI |
| `--enable-swagger-basic-auth` | `true` | Enable Swagger basic authentication |
| `--swagger-basic-auth-username` | `moon.rabbit` | Swagger basic auth username |
| `--swagger-basic-auth-password` | `rabbit.swagger` | Swagger basic auth password |
| `--enable-metrics` | `false` | Enable metrics endpoint |
| `--enable-metrics-basic-auth` | `true` | Enable metrics basic authentication |
| `--metrics-basic-auth-username` | `moon.rabbit` | Metrics basic auth username |
| `--metrics-basic-auth-password` | `rabbit.metrics` | Metrics basic auth password |
| `--enable-client-config` | `false` | Enable client configuration |

#### GORM Command Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--config`, `-c` | | `./config` | Config file path |
| `--force-gen`, `-f` | | `false` | Force generate code, overwrite existing |
| `--username` | | `root` | MySQL username |
| `--password` | | `123456` | MySQL password |
| `--host` | | `localhost` | MySQL host |
| `--port` | | `3306` | MySQL port |
| `--database` | | `rabbit` | MySQL database |
| `--params` | | `charset=utf8mb4,parseTime=true,loc=Asia/Shanghai` | MySQL connection parameters |
| `--biz`, `-b` | | `false` | Use biz namespace configuration |

See `rabbit run --help` and `rabbit run all --help` for all available flags.

### Example Usage

```bash
# Run all services (HTTP, gRPC, Job) with custom configuration file
rabbit run all -c ./config/server.yaml

# Run only HTTP server
rabbit run http -c ./config/server.yaml

# Run only gRPC server
rabbit run grpc -c ./config/server.yaml

# Run only Job server
rabbit run job -c ./config/server.yaml

# Run with multiple configuration files
rabbit run all -c ./config/server.yaml -c ./config/override.yaml

# Run with environment variables
MOON_RABBIT_HTTP_ADDRESS=0.0.0.0:8080 \
MOON_RABBIT_USE_DATABASE=true \
rabbit run all

# Run with database storage mode
rabbit run all \
  --http-address 0.0.0.0:8080 \
  --grpc-address 0.0.0.0:9090 \
  --job-address 0.0.0.0:9091 \
  --use-database true

# Run with file-based storage mode
rabbit run all \
  --http-address 0.0.0.0:8080 \
  --grpc-address 0.0.0.0:9090 \
  --job-address 0.0.0.0:9091 \
  --datasource-paths ./datasource,./config
```

## 📚 Commands

### Basic Commands

- `rabbit config` - Generate default configuration file
- `rabbit version` - Display version information

### Message Commands

- `rabbit send email` - Send email messages
- `rabbit send sms` - Send SMS messages
- `rabbit send feishu` - Send Feishu messages
- `rabbit apply` - Apply messages to queue
- `rabbit get` - Get message information
- `rabbit delete` - Delete messages

### Service Commands

- `rabbit run` - Start the Rabbit service
  - `rabbit run all` - Start all services (HTTP, gRPC, Job)
  - `rabbit run http` - Start only HTTP server
  - `rabbit run grpc` - Start only gRPC server
  - `rabbit run job` - Start only Job server
- `rabbit gorm` - GORM code generation and database migration tools
  - `rabbit gorm gen` - Generate GORM query code
  - `rabbit gorm migrate` - Migrate database schema

See `rabbit --help` for detailed command information.

## 🔧 Development

### Prerequisites

- Go 1.25+
- Make
- Protocol Buffers compiler (protoc)
- MySQL 8.0+ (for database mode)

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/aide-family/rabbit.git
cd rabbit

# Initialize the environment
make init

# Generate all code
make all

# Run tests
make test

# Run in development mode
make dev
```

### Project Structure

```
rabbit/
├── cmd/              # Command-line interface
├── internal/         # Internal packages
│   ├── biz/         # Business logic
│   ├── data/        # Data layer
│   ├── server/      # Server implementation
│   └── conf/        # Configuration
├── pkg/             # Public packages
├── proto/           # Protocol buffer definitions
├── config/          # Configuration files
├── deploy/          # Deployment configurations
└── Makefile         # Build scripts
```

## 🤝 Contributing

We welcome contributions! Please read our contributing guidelines before submitting PRs.

### Pull Request Process

1. **Fork the repository** and create your branch from `main`
2. **Create an issue** to discuss your changes (if it's a significant change)
3. **Make your changes** following our code style guidelines
4. **Add tests** for new features or bug fixes
5. **Update documentation** as needed
6. **Ensure all tests pass** (`make test`)
7. **Submit a Pull Request** with a clear description

#### PR Title Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Build process or auxiliary tool changes

**Example:**
```
feat(message): add email template support

Add support for email templates with dynamic variable substitution.
Templates can be defined in the configuration file and referenced
by name when sending emails.

Closes #123
```

#### PR Checklist

- [ ] Code follows the project's style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests pass
- [ ] No new warnings introduced
- [ ] Changes are backward compatible (or migration guide provided)

### Issue Reporting

When reporting issues, please include:

1. **Issue Type**: Bug, Feature Request, Question, etc.
2. **Description**: Clear description of the issue
3. **Steps to Reproduce**: For bugs, provide steps to reproduce
4. **Expected Behavior**: What you expected to happen
5. **Actual Behavior**: What actually happened
6. **Environment**: OS, Go version, Rabbit version
7. **Configuration**: Relevant configuration (sanitized)
8. **Logs**: Relevant log output
9. **Screenshots**: If applicable

#### Issue Template

```markdown
**Issue Type**: [Bug/Feature Request/Question]

**Description**:
<!-- Clear description of the issue -->

**Steps to Reproduce** (for bugs):
1. 
2. 
3. 

**Expected Behavior**:
<!-- What you expected to happen -->

**Actual Behavior**:
<!-- What actually happened -->

**Environment**:
- OS: 
- Go Version: 
- Rabbit Version: 

**Configuration**:
```yaml
<!-- Relevant configuration (sanitized) -->
```

**Logs**:
```
<!-- Relevant log output -->
```

**Additional Context**:
<!-- Any other relevant information -->
```

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Kratos](https://github.com/go-kratos/kratos) - A microservice-oriented framework
- [Cobra](https://github.com/spf13/cobra) - A CLI framework for Go

## 📞 Contact

- **Repository**: https://github.com/aide-family/rabbit
- **Issues**: https://github.com/aide-family/rabbit/issues
- **Email**: aidecloud@163.com

---

Made with ❤️ by [Aide Family](https://github.com/aide-family)
