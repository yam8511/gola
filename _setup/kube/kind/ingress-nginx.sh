kubectl label nodes gola-control-plane ingress-ready=true --overwrite
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
