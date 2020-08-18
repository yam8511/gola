# 建置 k8s cluster 以及 image registry
sh ./_setup/kind-with-registry.sh

# 建置 ingress-nginx
# 可以用 kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/usage.yaml 驗證
# 使用 curl localhost/foo 檢測，應該會顯示 foo
# 可以用 kubectl delete -f https://kind.sigs.k8s.io/examples/ingress/usage.yaml 刪除
sh ./_setup/ingress-nginx.sh
