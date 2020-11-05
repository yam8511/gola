#!/bin/sh
CRE='\033[0m'
CTEAL='\033[0;36m'
echo "${CTEAL}🛳  K3D UP${CRE}"

echo "⚙️  請輸入 kubernetes cluster 名稱 [預設: default]"
printf "> "
read cluster
cluster=${cluster:-default}

if k3d cluster ls ${cluster} 1>/dev/null 2>&1; then
    k3d cluster start ${cluster}
    exit 0
fi

echo "⚙️  輸入 kubernetes cluster port [預設: 16443]"
printf "> "
read cluster_port
cluster_port=${cluster_port:-16443}

echo "⚙️  請輸入 registry port (映像檔倉庫的port) [預設: 5000]"
printf "> "
read reg_port
reg_port=${reg_port:-5000}
reg_name='registry'

K3D_FD=$HOME/.gola/k3d
K8S_CLUSTER_FD=$HOME/.gola/k8s/cluster-${cluster}
REGISTRY_FILE=$K3D_FD/registry.yaml

mkdir -p $K3D_FD
mkdir -p $K8S_CLUSTER_FD/data

printf "mirrors:
    registry:${reg_port}:
        endpoint:
            - http://${reg_name}:${reg_port}" > $REGISTRY_FILE

k3d cluster create ${cluster} \
    --servers 1 \
    --agents 3 \
    --api-port 0.0.0.0:${cluster_port} \
    -v $K8S_CLUSTER_FD/data:/data \
    -v $REGISTRY_FILE:/etc/rancher/k3s/registries.yaml \
    -p 80:80@loadbalancer \
    -p 443:443@loadbalancer \
    -p 30306:30306@server[0] \
    -p 30679:30679@server[0] \
    -p 30001-30009:30001-30009@server[0] || exit 1

# create registry container unless it already exists
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
    docker run \
    -d --restart=always \
    -p "${reg_port}:5000" \
    --name "${reg_name}" \
    registry:2
fi

docker network connect "k3d-${cluster}" "${reg_name}"

echo "⚙️  請輸入要使用的 kubernetes namespace [預設: default]"
printf "> "
read namespace
namespace=${namespace:-default}

# 建立新的namespace
kubectl create namespace ${namespace}
kubectl create namespace devops
kubectl create namespace monitoring
kubectl create namespace logging
kubectl config set-context --current --namespace ${namespace}
