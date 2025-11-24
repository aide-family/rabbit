# Rabbit (Jade Rabbit) 🐰

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
./bin/rabbit run
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
   ./bin/rabbit run -c ./config/server.yaml
   ```

## ⚙️ Configuration

### Configuration File

The default configuration file is `config/server.yaml`. You can specify a custom path using the `--config` or `-c` flag.

### Environment Variables

Rabbit supports configuration through environment variables. All environment variables follow the pattern `MOON_RABBIT_*`.

#### Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_ENVIRONMENT` | `PROD` | Environment: DEV, TEST, PREVIEW, PROD |
| `MOON_RABBIT_NAME` | `moon.rabbit` | Service name |
| `MOON_RABBIT_HTTP_ADDRESS` | `0.0.0.0:8080` | HTTP server address |
| `MOON_RABBIT_GRPC_ADDRESS` | `0.0.0.0:9090` | gRPC server address |
| `MOON_RABBIT_HTTP_TIMEOUT` | `10s` | HTTP request timeout |
| `MOON_RABBIT_GRPC_TIMEOUT` | `10s` | gRPC request timeout |

#### Database Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_USE_DATABASE` | `false` | Enable database storage mode |
| `MOON_RABBIT_MAIN_HOST` | `localhost` | MySQL host |
| `MOON_RABBIT_MAIN_PORT` | `3306` | MySQL port |
| `MOON_RABBIT_MAIN_DATABASE` | `rabbit` | Database name |
| `MOON_RABBIT_MAIN_USERNAME` | `root` | MySQL username |
| `MOON_RABBIT_MAIN_PASSWORD` | `123456` | MySQL password |
| `MOON_RABBIT_MAIN_DEBUG` | `false` | Enable database debug mode |

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

#### Feature Flags

| Variable | Default | Description |
|----------|---------|-------------|
| `MOON_RABBIT_ENABLE_CLIENT_CONFIG` | `true` | Enable client configuration |
| `MOON_RABBIT_ENABLE_SWAGGER` | `true` | Enable Swagger UI |
| `MOON_RABBIT_ENABLE_METRICS` | `true` | Enable metrics endpoint |
| `MOON_RABBIT_CONFIG_PATHS` | `./datasource` | Configuration file paths (comma-separated) |

### Command Line Arguments

#### Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--namespace` | `-n` | `` | The namespace of the service |
| `--rabbit-config` | | `./.rabbit/` | The config file directory of the rabbit |

#### Run Command Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config`, `-c` | `` | Configuration file path |
| `--environment` | `PROD` | Environment: DEV, TEST, PREVIEW, PROD |
| `--http-address` | `0.0.0.0:8080` | HTTP server address |
| `--grpc-address` | `0.0.0.0:9090` | gRPC server address |
| `--use-database` | `false` | Enable database storage mode |
| `--config-paths` | `./datasource` | Configuration file paths |

See `rabbit run --help` for all available flags.

### Example Usage

```bash
# Run with custom configuration file
rabbit run -c ./config/server.yaml

# Run with environment variables
MOON_RABBIT_HTTP_ADDRESS=0.0.0.0:8080 \
MOON_RABBIT_USE_DATABASE=true \
rabbit run

# Run with command line flags
rabbit run \
  --http-address 0.0.0.0:8080 \
  --grpc-address 0.0.0.0:9090 \
  --use-database true \
  --config-paths ./datasource,./config
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

### Database Commands

- `rabbit gorm migrate` - Migrate database schema
- `rabbit gorm gen` - Generate GORM query code

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
