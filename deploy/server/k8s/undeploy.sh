#!/bin/bash
# Rabbit Kubernetes 卸载脚本
# 使用方法: ./undeploy.sh

set -e

NAMESPACE="rabbit"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================="
echo "Rabbit Kubernetes 卸载脚本"
echo "命名空间: ${NAMESPACE}"
echo "=========================================="

# 确认操作
read -p "确定要删除 Rabbit 服务吗？这将删除所有相关资源 (y/N): " CONFIRM
if [[ ! "${CONFIRM}" =~ ^[Yy]$ ]]; then
    echo "操作已取消"
    exit 0
fi

echo ""
echo "正在删除资源..."

# 删除资源（按依赖顺序）
echo "删除 Ingress..."
kubectl delete -f "${SCRIPT_DIR}/ingress.yaml" --ignore-not-found=true

echo "删除 Service..."
kubectl delete -f "${SCRIPT_DIR}/service.yaml" --ignore-not-found=true

echo "删除 Deployment..."
kubectl delete -f "${SCRIPT_DIR}/deployment.yaml" --ignore-not-found=true

echo "删除 RBAC..."
kubectl delete -f "${SCRIPT_DIR}/rbac.yaml" --ignore-not-found=true

echo "删除 ServiceAccount..."
kubectl delete -f "${SCRIPT_DIR}/serviceaccount.yaml" --ignore-not-found=true

echo "删除 ConfigMap..."
kubectl delete -f "${SCRIPT_DIR}/configmap.yaml" --ignore-not-found=true

echo "删除 Secret..."
kubectl delete secret rabbit-secrets -n "${NAMESPACE}" --ignore-not-found=true

read -p "是否删除 Namespace? (y/N): " DELETE_NS
if [[ "${DELETE_NS}" =~ ^[Yy]$ ]]; then
    echo "删除 Namespace..."
    kubectl delete -f "${SCRIPT_DIR}/namespace.yaml" --ignore-not-found=true
fi

echo ""
echo "=========================================="
echo "卸载完成！"
echo "=========================================="

