#!/bin/sh
CRE='\033[0m'
CTEAL='\033[0;36m'
echo "${CTEAL}🚢  KinD UP${CRE}"

echo "⚙️  請輸入 kubernetes cluster 名稱 [預設: default]"
printf "> "
read cluster
cluster=${cluster:-default}

echo "⚙️  輸入 kubernetes cluster port [預設: 16443]"
printf "> "
read cluster_port
cluster_port=${cluster_port:-16443}

echo "⚙️  輸入 registry port (映像檔倉庫的port) [預設: 5000]"
printf "> "
read reg_port
reg_port=${reg_port:-5000}
reg_name='registry'

# create registry container unless it already exists
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    -d --restart=always -p "${reg_port}:5000" --name "${reg_name}" \
    registry:2
fi

K8S_CLUSTER_FD=$HOME/.gola/k8s/cluster-${cluster}

mkdir -p $K8S_CLUSTER_FD/data

# create a cluster with the local registry enabled in containerd
cat <<EOF | kind create cluster --name ${cluster} --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerAddress: "127.0.0.1"
  apiServerPort: ${cluster_port}
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraMounts:
  - hostPath: ${K8S_CLUSTER_FD}/data
    containerPath: /data
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
  - containerPort: 30001
    hostPort: 30001
    protocol: TCP
  - containerPort: 30002
    hostPort: 30002
    protocol: TCP
  - containerPort: 30003
    hostPort: 30003
    protocol: TCP
  - containerPort: 30004
    hostPort: 30004
    protocol: TCP
  - containerPort: 30005
    hostPort: 30005
    protocol: TCP
  - containerPort: 30006
    hostPort: 30006
    protocol: TCP
  - containerPort: 30007
    hostPort: 30007
    protocol: TCP
  - containerPort: 30008
    hostPort: 30008
    protocol: TCP
  - containerPort: 30009
    hostPort: 30009
    protocol: TCP
  - containerPort: 30306
    hostPort: 30306
    protocol: TCP
  - containerPort: 30679
    hostPort: 30679
    protocol: TCP
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."registry:${reg_port}"]
    endpoint = ["http://${reg_name}:${reg_port}"]
EOF

# connect the registry to the cluster network
docker network connect "kind" "${reg_name}"

# tell https://tilt.dev to use the registry
# https://docs.tilt.dev/choosing_clusters.html#discovering-the-registry
for node in $(kind get nodes --name ${cluster}); do
  kubectl annotate node "${node}" "kind.x-k8s.io/registry=registry:${reg_port}";
done

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

# 建置 ingress-nginx
# 可以用 kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/usage.yaml 驗證
# 使用 curl localhost/foo 檢測，應該會顯示 foo
# 可以用 kubectl delete -f https://kind.sigs.k8s.io/examples/ingress/usage.yaml 刪除
kubectl label nodes ${cluster}-control-plane ingress-ready=true --overwrite
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
# 等待ingress建立完成
# kubectl wait --namespace ingress-nginx \
#   --for=condition=ready pod \
#   --selector=app.kubernetes.io/component=controller
kubectl delete -A ValidatingWebhookConfiguration ingress-nginx-admission
