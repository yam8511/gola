#!/bin/sh

# 把 myhostname的值 加入到 /etc/host
# e.g.
# 192.168.124.73 rancher.me
myhostname='rancher.me' # 注意不要使用 .local 結尾, ex. rancher.local
version='stable'

# 安裝rancher憑證
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.15.0/cert-manager.crds.yaml
kubectl create namespace cert-manager
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v0.15.0

# 等待 rancher憑證 部署完成
kubectl -n cert-manager rollout status deploy/cert-manager
kubectl -n cert-manager rollout status deploy/cert-manager-cainjector
kubectl -n cert-manager rollout status deploy/cert-manager-webhook

# 查看 rancher憑證 部署
kubectl -n cert-manager get deploy

# 安裝 rancher
helm repo add rancher-${version} https://releases.rancher.com/server-charts/${version}
kubectl create namespace cattle-system
helm install rancher rancher-${version}/rancher \
  --namespace cattle-system \
  --set hostname=${myhostname}

# 等待 rancher 部署完成
kubectl -n cattle-system rollout status deploy/rancher

# 查看 rancher 部署
kubectl -n cattle-system get deploy rancher
