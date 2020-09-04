#!/bin/sh
CRE='\033[0m'
CTEAL='\033[0;36m'
echo "${CTEAL}🚢  KinD UP${CRE}"

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

# create a cluster with the local registry enabled in containerd
cat <<EOF | kind create cluster --name gola --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
  - containerPort: 30007
    hostPort: 30007
    protocol: TCP
  - containerPort: 30008
    hostPort: 30008
    protocol: TCP
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
    endpoint = ["http://${reg_name}:${reg_port}"]
EOF

# connect the registry to the cluster network
docker network connect "kind" "${reg_name}"

# tell https://tilt.dev to use the registry
# https://docs.tilt.dev/choosing_clusters.html#discovering-the-registry
for node in $(kind get nodes); do
  kubectl annotate node "${node}" "kind.x-k8s.io/registry=localhost:${reg_port}";
done

# 建置 ingress-nginx
# 可以用 kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/usage.yaml 驗證
# 使用 curl localhost/foo 檢測，應該會顯示 foo
# 可以用 kubectl delete -f https://kind.sigs.k8s.io/examples/ingress/usage.yaml 刪除
kubectl label nodes gola-control-plane ingress-ready=true --overwrite
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
