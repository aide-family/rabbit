# Rabbit Kubernetes 部署配置

本目录包含 Rabbit 消息服务的 Kubernetes 部署配置文件，按照资源类型拆分为多个文件，便于管理和维护。

## 文件结构

```
k8s/
├── namespace.yaml          # Namespace 定义
├── configmap.yaml          # ConfigMap 配置
├── secret.yaml.example     # Secret 配置示例（需要手动创建）
├── serviceaccount.yaml     # ServiceAccount 定义
├── rbac.yaml               # RBAC 权限配置（Role 和 RoleBinding）
├── deployment.yaml         # Deployment 部署配置
├── service.yaml            # Service 服务配置
├── ingress.yaml            # Ingress 入口配置（可选）
├── kustomization.yaml      # Kustomize 配置文件
└── README.md              # 本文档
```

## 快速部署

### 方式一：使用 kubectl 逐个部署

```bash
# 1. 创建 Namespace
kubectl apply -f namespace.yaml

# 2. 创建 Secret（必须先创建，否则 Deployment 会失败）
kubectl create secret generic rabbit-secrets \
  --from-literal=jwt-secret='your-jwt-secret-key' \
  --from-literal=mysql-username='root' \
  --from-literal=mysql-password='your-mysql-password' \
  --namespace=rabbit

# 3. 创建 ConfigMap
kubectl apply -f configmap.yaml

# 4. 创建 ServiceAccount 和 RBAC
kubectl apply -f serviceaccount.yaml
kubectl apply -f rbac.yaml

# 5. 创建 Deployment
kubectl apply -f deployment.yaml

# 6. 创建 Service
kubectl apply -f service.yaml

# 7. 创建 Ingress（可选）
kubectl apply -f ingress.yaml
```

### 方式二：使用 Kustomize 统一部署（推荐）

```bash
# 1. 先创建 Secret
kubectl create secret generic rabbit-secrets \
  --from-literal=jwt-secret='your-jwt-secret-key' \
  --from-literal=mysql-username='root' \
  --from-literal=mysql-password='your-mysql-password' \
  --namespace=rabbit

# 2. 使用 Kustomize 部署所有资源
kubectl apply -k deploy/server/k8s/
```

### 方式三：使用脚本一键部署

```bash
# 创建部署脚本
cat > deploy.sh << 'EOF'
#!/bin/bash
set -e

# 创建 Namespace
kubectl apply -f namespace.yaml

# 创建 Secret（如果不存在）
if ! kubectl get secret rabbit-secrets -n rabbit &>/dev/null; then
  echo "Creating secret..."
  kubectl create secret generic rabbit-secrets \
    --from-literal=jwt-secret='change-me-in-production' \
    --from-literal=mysql-username='root' \
    --from-literal=mysql-password='your-mysql-password' \
    --namespace=rabbit
fi

# 部署其他资源
kubectl apply -f configmap.yaml
kubectl apply -f serviceaccount.yaml
kubectl apply -f rbac.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml

echo "Deployment completed!"
EOF

chmod +x deploy.sh
./deploy.sh
```

## 配置说明

### 1. Namespace

定义了 `rabbit` 命名空间，用于隔离 Rabbit 服务的所有资源。

### 2. ConfigMap

包含 Rabbit 服务的主配置文件 `server.yaml`，支持环境变量替换。

**重要配置项：**
- `useDatabase`: 是否使用数据库模式（false 使用文件配置，true 使用数据库）
- `main.host`: MySQL 服务地址（如果使用数据库模式）
- `configPaths`: 文件配置模式的配置文件路径

### 3. Secret

包含敏感信息，如 JWT 密钥、MySQL 密码等。

**创建方式：**

```bash
# 方式一：使用 kubectl 命令
kubectl create secret generic rabbit-secrets \
  --from-literal=jwt-secret='your-jwt-secret-key' \
  --from-literal=mysql-username='root' \
  --from-literal=mysql-password='your-mysql-password' \
  --namespace=rabbit

# 方式二：使用 YAML 文件（需要 base64 编码）
# 参考 secret.yaml.example 文件
```

**生产环境建议：**
- 使用 Sealed Secrets 或 External Secrets Operator
- 使用密钥管理服务（如 AWS Secrets Manager、HashiCorp Vault）

### 4. ServiceAccount 和 RBAC

为 Rabbit 服务创建了专用的 ServiceAccount 和 RBAC 权限，用于 Kubernetes 服务发现和注册。

**权限说明：**
- 可以获取、列出、监听 endpoints 和 services 资源
- 用于 Kubernetes 注册中心的服务发现

### 5. Deployment

定义了 Rabbit 服务的部署配置。

**关键配置：**
- `replicas`: 副本数（默认 2）
- `image`: 容器镜像（默认 `ghcr.io/aide-family/rabbit:v0.0.1`）
- `resources`: 资源限制和请求
- `livenessProbe` 和 `readinessProbe`: 健康检查

**自定义启动命令：**

如果需要自定义启动命令，可以修改 Deployment 的 `command` 字段：

```yaml
containers:
- name: rabbit
  command: ["/usr/local/bin/rabbit"]
  args: ["run", "--use-database", "true", "--main-host", "mysql-service"]
```

### 6. Service

定义了 ClusterIP 类型的 Service，用于集群内部访问。

**端口：**
- `8080`: HTTP API
- `9090`: gRPC API

### 7. Ingress

定义了外部访问入口（可选）。

**配置说明：**
- 默认使用 Nginx Ingress Controller
- 需要根据实际环境修改 `host` 和 `ingress.class`
- 支持 SSL/TLS 配置（需要取消注释相关配置）

## 环境配置

### 文件配置模式（默认）

不需要数据库，使用配置文件：

```yaml
# 在 deployment.yaml 中设置
env:
- name: MOON_RABBIT_USE_DATABASE
  value: "false"
- name: MOON_RABBIT_CONFIG_PATHS
  value: "./datasource"
```

### 数据库模式

需要 MySQL 数据库：

```yaml
# 在 deployment.yaml 中设置
env:
- name: MOON_RABBIT_USE_DATABASE
  value: "true"
- name: MOON_RABBIT_MAIN_HOST
  value: "mysql-service"  # MySQL Service 名称
```

## 验证部署

```bash
# 查看 Pod 状态
kubectl get pods -n rabbit

# 查看 Service
kubectl get svc -n rabbit

# 查看 Deployment
kubectl get deployment -n rabbit

# 查看日志
kubectl logs -f deployment/rabbit -n rabbit

# 测试服务
kubectl port-forward svc/rabbit-service 8080:8080 -n rabbit
curl http://localhost:8080/health
```

## 更新部署

```bash
# 更新镜像版本
kubectl set image deployment/rabbit rabbit=ghcr.io/aide-family/rabbit:v0.0.2 -n rabbit

# 更新配置
kubectl apply -f configmap.yaml
kubectl rollout restart deployment/rabbit -n rabbit

# 查看更新状态
kubectl rollout status deployment/rabbit -n rabbit
```

## 扩缩容

```bash
# 扩容到 3 个副本
kubectl scale deployment/rabbit --replicas=3 -n rabbit

# 或者修改 deployment.yaml 中的 replicas 字段
```

## 卸载

```bash
# 删除所有资源
kubectl delete -k deploy/server/k8s/

# 或者逐个删除
kubectl delete -f ingress.yaml
kubectl delete -f service.yaml
kubectl delete -f deployment.yaml
kubectl delete -f rbac.yaml
kubectl delete -f serviceaccount.yaml
kubectl delete -f configmap.yaml
kubectl delete secret rabbit-secrets -n rabbit
kubectl delete -f namespace.yaml
```

## 故障排查

### 查看 Pod 状态

```bash
kubectl describe pod <pod-name> -n rabbit
```

### 查看事件

```bash
kubectl get events -n rabbit --sort-by='.lastTimestamp'
```

### 查看日志

```bash
# 查看所有 Pod 日志
kubectl logs -l app=rabbit -n rabbit

# 查看特定 Pod 日志
kubectl logs <pod-name> -n rabbit

# 实时查看日志
kubectl logs -f deployment/rabbit -n rabbit
```

### 常见问题

1. **Pod 无法启动**
   - 检查 Secret 是否已创建
   - 检查 ConfigMap 是否正确
   - 查看 Pod 日志和事件

2. **健康检查失败**
   - 确认 `/health` 端点可用
   - 检查端口配置是否正确

3. **无法连接数据库**
   - 确认 MySQL Service 存在且可访问
   - 检查 Secret 中的数据库密码是否正确

## 生产环境建议

1. **安全性**
   - 使用强密码和密钥
   - 启用 TLS/SSL
   - 使用 NetworkPolicy 限制网络访问
   - 定期轮换 Secret

2. **高可用**
   - 设置多个副本（至少 2 个）
   - 使用 PodDisruptionBudget
   - 配置多可用区部署

3. **监控**
   - 配置 Prometheus 监控
   - 设置告警规则
   - 启用日志聚合

4. **资源管理**
   - 根据实际负载调整资源限制
   - 使用 HPA（Horizontal Pod Autoscaler）自动扩缩容
   - 配置资源配额

