#!/bin/bash
# Rabbit Kubernetes 部署脚本
# 使用方法: ./deploy.sh [环境]
# 环境: dev, test, prod (默认: dev)

set -e

ENVIRONMENT=${1:-dev}
NAMESPACE="rabbit"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================="
echo "Rabbit Kubernetes 部署脚本"
echo "环境: ${ENVIRONMENT}"
echo "命名空间: ${NAMESPACE}"
echo "=========================================="

# 检查 kubectl 是否可用
if ! command -v kubectl &> /dev/null; then
    echo "错误: kubectl 未安装或不在 PATH 中"
    exit 1
fi

# 检查是否已连接到 Kubernetes 集群
if ! kubectl cluster-info &> /dev/null; then
    echo "错误: 无法连接到 Kubernetes 集群"
    exit 1
fi

echo ""
echo "步骤 1/7: 创建 Namespace..."
kubectl apply -f "${SCRIPT_DIR}/namespace.yaml"

echo ""
echo "步骤 2/7: 检查并创建 Secret..."
if ! kubectl get secret rabbit-secrets -n "${NAMESPACE}" &>/dev/null; then
    echo "Secret 不存在，正在创建..."
    echo "提示: 生产环境请使用强密码和密钥！"
    read -sp "请输入 JWT Secret (默认: change-me-in-production): " JWT_SECRET
    JWT_SECRET=${JWT_SECRET:-change-me-in-production}
    
    read -sp "请输入 MySQL 用户名 (默认: root): " MYSQL_USER
    MYSQL_USER=${MYSQL_USER:-root}
    
    read -sp "请输入 MySQL 密码: " MYSQL_PASSWORD
    MYSQL_PASSWORD=${MYSQL_PASSWORD:-your-mysql-password}
    
    kubectl create secret generic rabbit-secrets \
        --from-literal=jwt-secret="${JWT_SECRET}" \
        --from-literal=mysql-username="${MYSQL_USER}" \
        --from-literal=mysql-password="${MYSQL_PASSWORD}" \
        --namespace="${NAMESPACE}"
    echo ""
    echo "Secret 创建成功"
else
    echo "Secret 已存在，跳过创建"
fi

echo ""
echo "步骤 3/7: 创建 ConfigMap..."
kubectl apply -f "${SCRIPT_DIR}/configmap.yaml"

echo ""
echo "步骤 4/7: 创建 ServiceAccount..."
kubectl apply -f "${SCRIPT_DIR}/serviceaccount.yaml"

echo ""
echo "步骤 5/7: 创建 RBAC 配置..."
kubectl apply -f "${SCRIPT_DIR}/rbac.yaml"

echo ""
echo "步骤 6/7: 创建 Deployment..."
kubectl apply -f "${SCRIPT_DIR}/deployment.yaml"

echo ""
echo "步骤 7/7: 创建 Service..."
kubectl apply -f "${SCRIPT_DIR}/service.yaml"

echo ""
read -p "是否创建 Ingress? (y/N): " CREATE_INGRESS
if [[ "${CREATE_INGRESS}" =~ ^[Yy]$ ]]; then
    echo "创建 Ingress..."
    kubectl apply -f "${SCRIPT_DIR}/ingress.yaml"
fi

echo ""
echo "=========================================="
echo "部署完成！"
echo "=========================================="
echo ""
echo "查看部署状态:"
echo "  kubectl get pods -n ${NAMESPACE}"
echo "  kubectl get svc -n ${NAMESPACE}"
echo "  kubectl get deployment -n ${NAMESPACE}"
echo ""
echo "查看日志:"
echo "  kubectl logs -f deployment/rabbit -n ${NAMESPACE}"
echo ""
echo "端口转发（本地测试）:"
echo "  kubectl port-forward svc/rabbit-service 8080:8080 -n ${NAMESPACE}"
echo ""

