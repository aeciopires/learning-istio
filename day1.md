# Descomplicando Istio

<!-- TOC -->

- [Descomplicando Istio](#descomplicando-istio)
- [Day 1](#day-1)

<!-- TOC -->

# Day 1

Comandos executados em sequencia durante o treinamento.

Crie um cluster k8s usando uma das opções abaixo:

* [kubeadm](kubeadm.md);
* [kind](kind.md);


Instale o Istio com os seguintes comandos:

```bash
#------- Specifics (master)
# References:
# https://istio.io/latest/docs/setup/getting-started/
# Allow all these ports: https://istio.io/latest/docs/ops/deployment/requirements/
curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.11.4 sh -
export ISTIO_DIR_BASE="/home/ubuntu/istio-1.11.4"
cd $ISTIO_DIR_BASE
export PATH="$PATH:$ISTIO_DIR_BASE/bin"

# Complementary files
git clone https://github.com/aeciopires/learning-istio
export COMPLEMENTARY_FILES=/home/ubuntu/learning-istio/files

istioctl x precheck

# Install profile of demonstration
istioctl install --set profile=demo -y

# Install Addons (grafana, prometheus, kiali, jaeger, zipkin)
kubectl apply -f $ISTIO_DIR_BASE/samples/addons
kubectl rollout status deployment/kiali -n istio-system
kubectl rollout status deployment/grafana -n istio-system
kubectl rollout status deployment/prometheus -n istio-system
kubectl rollout status deployment/jaeger -n istio-system

# List api-resources and crd
kubectl api-resources | grep istio
kubectl get crd | grep istio
kubectl get svc -n istio-system
kubectl get pods -n istio-system

# First Deployment
kubectl label namespace default istio-injection=enabled

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/platform/kube/bookinfo.yaml
kubectl get services
kubectl get pods

# Verify everything is working correctly up to this point. Run this command to see if the app is running inside the cluster and serving HTML pages by checking for the page title in the response:
kubectl exec "$(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl -sS productpage:9080/productpage | grep -o "<title>.*</title>"

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/bookinfo-gateway.yaml
kubectl get gateway
kubectl get virtualservices -o yaml
kubectl get services
kubectl get pods
kubectl get virtualservices.networking.istio.io bookinfo -o yaml
kubectl get virtualservices.networking.istio.io
kubectl describe virtualservices.networking.istio.io bookinfo
kubectl get gateways.networking.istio.io
kubectl describe gateways.networking.istio.io bookinfo-gateway

# Port forward in background for access Productpage application
# Allow port 9080/TCP
kubectl port-forward svc/productpage 9080:9080 -n default --address=0.0.0.0 > /dev/null 2>&1 &
# Access URL: http://master:9080/productpage

# Port forward in background for access Kiali
# Allow port 20001/TCP
kubectl port-forward svc/kiali 20001:20001 -n istio-system --address=0.0.0.0 > /dev/null 2>&1 &
# Access URL: http://master:20001/kiali

# Port forward in background for access Grafana
# Allow port 3000/TCP
kubectl port-forward svc/grafana 3000:3000 -n istio-system --address=0.0.0.0 > /dev/null 2>&1 &
# Access URL: http://master:3000

# Port forward in background for access Prometheus
# Allow port 9090/TCP
kubectl port-forward svc/prometheus 9090:9090 -n istio-system --address=0.0.0.0 > /dev/null 2>&1 &
# Access URL: http://master:9090

# Port forward in background for access Jaeger
# Allow port 16685/TCP
kubectl port-forward svc/tracing 16685:80 -n istio-system --address=0.0.0.0 > /dev/null 2>&1 &
# Access URL: http://master:16685

# List jobs running in background
jobs -l

# Finish jobs running in background
kill %1 %2 %3 %4 %5
```
