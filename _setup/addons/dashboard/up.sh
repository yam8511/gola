# create k8s dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta8/aio/deploy/recommended.yaml
kubectl create clusterrolebinding default-admin --clusterrole cluster-admin --serviceaccount=default:default

echo "=========================================="
echo "== you can type below, then look dashboard"
echo "== kubectl proxy  # port 8001 listening..."
echo "== open http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/ # open browser"
echo "== kubectl get secrets # find 'default-token-xxxxx'"
echo "== kubectl describe secrets default-token-xxxxx # will find 'token' to login k8s dashboard"
echo "== reference: https://istio.io/latest/zh/docs/setup/platform-setup/kind/"
echo "=========================================="
