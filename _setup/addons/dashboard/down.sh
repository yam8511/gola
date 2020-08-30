# create k8s dashboard
kubectl delete -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta8/aio/deploy/recommended.yaml
kubectl delete clusterrolebinding default-admin 
kubectl delete serviceaccount default
