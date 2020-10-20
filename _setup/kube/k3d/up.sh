#!/bin/sh
CRE='\033[0m'
CTEAL='\033[0;36m'
echo "${CTEAL}🛳  K3D UP${CRE}"

if k3d cluster ls gola 1>/dev/null 2>&1; then
    k3d cluster start gola
    exit 0
fi

echo "⚙️  輸入 registry port (映像檔倉庫的port) [預設: 5000]"
printf "> "
read reg_port
reg_port=${reg_port:-5000}
reg_name='registry'

echo "mirrors:
    \"localhost:${reg_port}\":
    endpoint:
        - http://${reg_name}:${reg_port}" > $HOME/.cache/gola/kube/k3d/registry.yaml


k3d cluster create gola \
    --servers 1 \
    --agents 3 \
    --api-port 0.0.0.0:16443 \
    -v $HOME/.cache/gola/data:/data \
    -v $HOME/.cache/gola/kube/k3d/registry.yaml:/etc/rancher/k3s/registries.yaml \
    -p 80:80@loadbalancer \
    -p 443:443@loadbalancer \
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

docker network connect "k3d-gola" "${reg_name}"

# 建立新的namespace
kubectl create namespace devops
