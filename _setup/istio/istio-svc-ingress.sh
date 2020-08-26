# 如果連得到容器IP，直接使用以下指令即可
# export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
# export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
# export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}')
# export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT

# 1. 找出 Istio 的Ingress Port
INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')

# 2. 把 istio 的 ingress service 匯出來
kubectl get svc istio-ingressgateway -n istio-system -o yaml > gateway.yaml

# 3. 調整Yaml
#     - http2.nodePort: `3xxxx`
#     - https.nodePort: `3xxxx`
#     - type: `NodePort`
cat gateway.yaml | sed 's/type: LoadBalancer/type: NodePort/g' | sed "s/$INGRESS_PORT/30007/g" | sed "s/$SECURE_INGRESS_PORT/30008/g" > gateway.yaml

# 4. 刪除原本istio service
kubectl delete -f gateway.yaml

# 5. 使用自己調整過的service
kubectl apply -f gateway.yaml

# 6. 刪除 gateway.yaml
rm gateway.yaml
