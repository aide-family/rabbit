# Kubernetes 部署指南

本文档介绍如何在 Kubernetes 集群中部署 Rabbit 消息服务。

## 前置要求

- Kubernetes 1.20+
- kubectl 已配置并可以访问集群
- 已构建的 Rabbit Docker 镜像并推送到镜像仓库

## 快速部署

### 1. 准备镜像

首先需要将 Rabbit 镜像推送到可访问的镜像仓库：

```bash
# 构建镜像
docker build -t rabbit-local:latest .

# 标记镜像（替换为你的镜像仓库地址）
docker tag rabbit-local:latest your-registry/rabbit:latest

# 推送镜像
docker push your-registry/rabbit:latest
```

### 2. 创建配置文件

创建 Kubernetes 部署配置文件 `rabbit.yaml`：

```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: rabbit-config
  namespace: moon
data:
  server.yaml: |
    environment: PROD
    enableClientConfig: false
    enableSwagger: true
    enableMetrics: true
    useDatabase: false
    
    server:
      name: rabbit
      namespace: moon
      http:
        address: "0.0.0.0:8080"
      grpc:
        address: "0.0.0.0:9090"
      job:
        address: "0.0.0.1:9091"
    
    jwt:
      secret: "your-jwt-secret"
      expire: "600s"
      issuer: "rabbit"
---
apiVersion: v1
kind: Secret
metadata:
  name: rabbit-secret
  namespace: moon
type: Opaque
stringData:
  jwt-secret: "your-jwt-secret-key"
  mysql-password: "your-mysql-password"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rabbit
  namespace: moon
  labels:
    app: rabbit
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rabbit
  template:
    metadata:
      labels:
        app: rabbit
    spec:
      containers:
      - name: rabbit
        image: your-registry/rabbit:latest
        imagePullPolicy: Always
        ports:
        - name: http
          containerPort: 8080
          protocol: TCP
        - name: grpc
          containerPort: 9090
          protocol: TCP
        env:
        - name: MOON_RABBIT_ENVIRONMENT
          value: "PROD"
        - name: MOON_RABBIT_HTTP_ADDRESS
          value: "0.0.0.0:8080"
        - name: MOON_RABBIT_GRPC_ADDRESS
          value: "0.0.0.0:9090"
        - name: MOON_RABBIT_USE_DATABASE
          value: "false"
        - name: MOON_RABBIT_JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: rabbit-secret
              key: jwt-secret
        volumeMounts:
        - name: config
          mountPath: /moon/config
          readOnly: true
        - name: datasource
          mountPath: /moon/datasource
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
      volumes:
      - name: config
        configMap:
          name: rabbit-config
      - name: datasource
        emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: rabbit
  namespace: moon
  labels:
    app: rabbit
spec:
  type: ClusterIP
  ports:
  - name: http
    port: 8080
    targetPort: http
    protocol: TCP
  - name: grpc
    port: 9090
    targetPort: grpc
    protocol: TCP
  selector:
    app: rabbit
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: rabbit-ingress
  namespace: moon
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: rabbit.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: rabbit
            port:
              number: 8080
```

### 3. 创建命名空间

```bash
# 创建命名空间（如果不存在）
kubectl create namespace moon --dry-run=client -o yaml | kubectl apply -f -

# 或者直接创建
kubectl create namespace moon
```

### 4. 部署服务

```bash
kubectl apply -f deploy/server/k8s/rabbit.yaml
```

### 5. 检查部署状态

```bash
# 查看命名空间
kubectl get namespace moon

# 查看 Pod 状态
kubectl get pods -n moon

# 查看服务状态
kubectl get svc -n moon

# 查看部署状态
kubectl get deployment -n moon

# 查看详细日志
kubectl logs -f deployment/rabbit -n moon
```

## 配置说明

### 命名空间

所有资源都部署在 `moon` 命名空间中，便于管理和隔离。

### ConfigMap

用于存储配置文件，可以通过修改 ConfigMap 来更新配置：

```bash
# 编辑 ConfigMap
kubectl edit configmap rabbit-config -n moon

# 查看 ConfigMap
kubectl get configmap rabbit-config -n moon -o yaml
```

### Secret

用于存储敏感信息（如 JWT 密钥、数据库密码）：

```bash
# 创建 Secret
kubectl create secret generic rabbit-secret \
  --from-literal=jwt-secret=your-secret \
  --from-literal=mysql-password=your-password \
  -n moon

# 查看 Secret
kubectl get secret rabbit-secret -n moon
```

### Deployment

- **副本数**：默认 1 个副本，可以根据需要调整
- **资源限制**：建议设置适当的资源请求和限制
- **健康检查**：配置了存活探针和就绪探针

### Service

- **类型**：ClusterIP（集群内部访问）
- **端口**：HTTP 8080，gRPC 9090

### Ingress

如果需要外部访问，可以配置 Ingress。示例中使用 Nginx Ingress Controller。

## 数据库存储模式

如果需要使用数据库存储模式，需要：

1. **部署 MySQL**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
  namespace: moon
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:8.0
        env:
        - name: MYSQL_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: rabbit-secret
              key: mysql-password
        - name: MYSQL_DATABASE
          value: "rabbit"
        ports:
        - containerPort: 3306
---
apiVersion: v1
kind: Service
metadata:
  name: mysql
  namespace: moon
spec:
  selector:
    app: mysql
  ports:
  - port: 3306
    targetPort: 3306
```

2. **更新 Rabbit 配置**

在 Deployment 的环境变量中添加：

```yaml
env:
- name: MOON_RABBIT_USE_DATABASE
  value: "true"
- name: MOON_RABBIT_MAIN_HOST
  value: "mysql"
- name: MOON_RABBIT_MAIN_DATABASE
  value: "rabbit"
- name: MOON_RABBIT_MAIN_USERNAME
  value: "root"
- name: MOON_RABBIT_MAIN_PASSWORD
  valueFrom:
    secretKeyRef:
      name: rabbit-secret
      key: mysql-password
```

## 扩缩容

### 手动扩缩容

```bash
# 扩展到 3 个副本
kubectl scale deployment rabbit --replicas=3 -n moon

# 查看副本状态
kubectl get deployment rabbit -n moon
```

### 自动扩缩容（HPA）

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: rabbit-hpa
  namespace: moon
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: rabbit
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

应用 HPA：

```bash
kubectl apply -f rabbit-hpa.yaml
kubectl get hpa -n moon
```

## 服务管理

### 更新部署

```bash
# 更新镜像
kubectl set image deployment/rabbit rabbit=your-registry/rabbit:v1.1.0 -n moon

# 查看滚动更新状态
kubectl rollout status deployment/rabbit -n moon

# 回滚到上一个版本
kubectl rollout undo deployment/rabbit -n moon

# 查看更新历史
kubectl rollout history deployment/rabbit -n moon
```

### 查看日志

```bash
# 查看所有 Pod 日志
kubectl logs -f -l app=rabbit -n moon

# 查看特定 Pod 日志
kubectl logs -f rabbit-xxxxx -n moon

# 查看前 100 行日志
kubectl logs --tail=100 rabbit-xxxxx -n moon
```

### 进入容器

```bash
kubectl exec -it deployment/rabbit -n moon -- sh
```

### 删除部署

```bash
# 删除所有资源
kubectl delete -f deploy/server/k8s/rabbit.yaml

# 或删除命名空间（会删除命名空间下的所有资源）
kubectl delete namespace moon
```

## 监控和健康检查

### 健康检查端点

- **存活探针**：`/health`
- **就绪探针**：`/health`

### 查看 Pod 状态

```bash
# 查看 Pod 详细信息
kubectl describe pod rabbit-xxxxx -n moon

# 查看事件
kubectl get events -n moon --sort-by='.lastTimestamp'
```

## 生产环境建议

1. **使用 StatefulSet**：如果需要持久化存储，考虑使用 StatefulSet
2. **配置资源限制**：为 Pod 设置适当的资源请求和限制
3. **使用持久化存储**：数据库存储模式建议使用 PVC
4. **配置网络策略**：限制 Pod 之间的网络访问
5. **启用 TLS**：在 Ingress 中配置 TLS 证书
6. **配置监控**：集成 Prometheus 和 Grafana
7. **日志收集**：配置日志收集系统（如 ELK、Loki）

### 持久化存储示例

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: rabbit-data
  namespace: moon
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: standard
```

在 Deployment 中使用：

```yaml
volumeMounts:
- name: datasource
  mountPath: /moon/datasource
volumes:
- name: datasource
  persistentVolumeClaim:
    claimName: rabbit-data
```

## 故障排查

### Pod 无法启动

```bash
# 查看 Pod 状态
kubectl get pods -n moon

# 查看 Pod 详细信息
kubectl describe pod rabbit-xxxxx -n moon

# 查看日志
kubectl logs rabbit-xxxxx -n moon
```

### 服务无法访问

```bash
# 检查 Service
kubectl get svc rabbit -n moon

# 检查 Endpoints
kubectl get endpoints rabbit -n moon

# 测试服务连接
kubectl run -it --rm debug --image=busybox --restart=Never -- sh
# 在容器内执行
wget -qO- http://rabbit:8080/health
```

### 配置问题

```bash
# 查看 ConfigMap
kubectl get configmap rabbit-config -n moon -o yaml

# 查看环境变量
kubectl exec rabbit-xxxxx -n moon -- env | grep MOON_RABBIT
```

## 更多信息

- Kubernetes 官方文档：https://kubernetes.io/docs/
- 完整的配置选项请参考主 README 文档
- 环境变量列表请参考主 README 中的环境变量部分
