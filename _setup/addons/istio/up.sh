# 安裝Istio
istioctl install --set profile=demo

CYELLOW='\033[1;33m'
CRE='\033[0m'
CTEAL='\033[0;34m'
CRED='\033[0;31m'

echo "

若希望把Istio的Service從${CRED}LoadBalancer${CRE}改為${CRED}NodePort${CRE}，可以在終端機複製貼上以下指令：
${CTEAL}
# 1. 找出 Istio 的Ingress Port
INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
# 2. 把 istio 的 ingress service 匯出來
kubectl get svc istio-ingressgateway -n istio-system -o yaml > gateway.yaml
# 3. 調整Yaml
cat gateway.yaml | sed 's/type: LoadBalancer/type: NodePort/g' | sed \"s/\$INGRESS_PORT/30007/g\" | sed \"s/$\SECURE_INGRESS_PORT/30008/g\" > gateway.yaml
# 4. 刪除原本istio service
kubectl delete -f gateway.yaml
# 5. 使用自己調整過的service
kubectl apply -f gateway.yaml
# 6. 刪除 gateway.yaml
rm gateway.yaml
${CRE}
詳細步驟或說明，請見 => `pwd`/_setup/addons/istio/istio-svc-ingress.sh
---
${CYELLOW}
啟動Istio功能在自己的服務，可以下指令
> kubectl label namespace default istio-injection=enabled --overwrite

關閉Istio功能
> kubectl label namespace default istio-injection-
${CRE}"

# 調整istio-system的ingres service，從LoadBalaner改為NodePort
# sh ./_setup/addons/istio/istio-svc-ingress.sh
