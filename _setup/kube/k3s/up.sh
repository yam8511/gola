CRE='\e[0m'
CTEAL='\e[0;36m'
printf "${CTEAL}🛳  K3S UP${CRE}\n"

echo "⚙️  輸入 registry port (映像檔倉庫的port) [預設: 5000]"
printf "> "
read reg_port
reg_port=${reg_port:-5000}
reg_name='registry'

echo "⚙️  請輸入 kubernetes cluster 名稱 [預設: default]"
printf "> "
read cluster
cluster=${cluster:-default}

echo "⚙️  輸入 kubernetes cluster port [預設: 16443]"
printf "> "
read cluster_port
cluster_port=${cluster_port:-16443}

echo "⚙️  輸入 Node IP [預設: 127.0.0.1]"
printf "> "
read node_ip
node_ip=${node_ip:-127.0.0.1}

K3S_FD=$HOME/.gola/k3s
REGISTRY_FILE=$K3S_FD/registry.yaml

mkdir -p $K3S_FD

echo "mirrors:
    \"registry:${reg_port}\":
    endpoint:
        - http://${node_ip}:${reg_port}" > $REGISTRY_FILE

INSTALL_K3S_SKIP_DOWNLOAD=true INSTALL_K3S_NAME=${cluster} INSTALL_K3S_EXEC="--docker --https-listen-port ${cluster_port} --node-ip ${node_ip} --write-kubeconfig-mode 644 --write-kubeconfig $HOME/.kube/config --private-registry $REGISTRY_FILE" sh ./_setup/kube/k3s/k3s.sh

# create registry container unless it already exists
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
    docker run \
    -d --restart=always \
    -p "${reg_port}:5000" \
    --name "${reg_name}" \
    registry:2
fi

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
