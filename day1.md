# Descomplicando Istio

<!-- TOC -->

- [Descomplicando Istio](#descomplicando-istio)
- [Day 1](#day-1)

<!-- TOC -->

# Day 1

Comandos executados em sequência durante o treinamento.

Instale o **Docker** com as instruções da página: [docker](docker.md).

Crie um **cluster Kubernetes** usando uma das opções abaixo:

* [kubeadm](kubeadm.md);
* [kind](kind.md);

Instale o **Helm** com as instruções da página: [helm](helm.md).

Instale o **Istio** com os seguintes comandos:

```bash
#------- Specifics (master)
# References:
# https://istio.io/latest/docs/setup/getting-started/
# https://istio.io/latest/docs/setup/install/helm/
# Allow all these ports: https://istio.io/latest/docs/ops/deployment/requirements/

# Configure the Helm repository:
helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update

# Create the namespace, istio-system, for the Istio components:
kubectl create namespace istio-system

# Install the Istio base chart which contains cluster-wide Custom Resource Definitions (CRDs) which must be installed prior to the deployment of the Istio control plane.
helm install istio-base istio/base -n istio-system --set defaultRevision=default

# Validate the CRD installation with the helm ls command:
helm ls -n istio-system

# Install the Istio discovery chart which deploys the istiod service:
helm install istiod istio/istiod -n istio-system --wait --timeout=900s

# Check istiod service is successfully installed and its pods are running:
kubectl get deployments -n istio-system --output wide

# Install an ingress gateway:
kubectl create namespace istio-ingress
helm install istio-ingress istio/gateway -n istio-ingress --wait --timeout=900s

# To uninstall istio follow the instructions of the page
# https://istio.io/latest/docs/setup/install/helm/#uninstall
```

Baixe este repositório para obter os arquivos complementares

```bash
git clone https://github.com/aeciopires/learning-istio
export COMPLEMENTARY_FILES=$HOME/learning-istio/files
```

Faça o deploy da aplicação de exemplo chamada **Bookinfo**:

```bash
MY_NAMESPACE='myapp'
kubectl create namespace $MY_NAMESPACE

kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.20/samples/bookinfo/platform/kube/bookinfo.yaml -n $MY_NAMESPACE
```

Visualize os objetos/recursos da aplicação:

```bash
kubectl get all -n $MY_NAMESPACE
```

Teste o acesso a aplicação com o seguinte comando:

```bash
kubectl -n $MY_NAMESPACE exec "$(kubectl -n $MY_NAMESPACE get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl -sS productpage:9080/productpage | grep -o "<title>.*</title>"
```

Crie um ingressGateway e um VirtualService para a aplicação ser exposta fora do cluster usando o Istio e o loadbalancer.

```bash
kubectl -n $MY_NAMESPACE apply -f $COMPLEMENTARY_FILES/bookinfo-gateway.yaml
```

Obtenha as informações a cerca do IP e Porta do IngressGateway usado pela aplicação:

```bash
INGRESS_HOST=$(kubectl get svc istio-ingress -n istio-ingress -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

INGRESS_PORT=$(kubectl get svc istio-ingress -n istio-ingress -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')

SECURE_INGRESS_PORT=$(kubectl get svc istio-ingress -n istio-ingress -o jsonpath='{.spec.ports[?(@.name=="https")].port}')

GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT

echo "$GATEWAY_URL"

curl "http://$GATEWAY_URL/productpage"
```

Visualize os objetos da aplicação recém Bookinfo.

```bash
kubectl -n $MY_NAMESPACE get gateway
kubectl -n $MY_NAMESPACE get virtualservices -o yaml
kubectl -n $MY_NAMESPACE get services
kubectl -n $MY_NAMESPACE get pods
kubectl -n $MY_NAMESPACE get virtualservices.networking.istio.io bookinfo -o yaml
kubectl -n $MY_NAMESPACE get virtualservices.networking.istio.io
kubectl -n $MY_NAMESPACE describe virtualservices.networking.istio.io bookinfo
kubectl -n $MY_NAMESPACE get gateways.networking.istio.io
kubectl -n $MY_NAMESPACE describe gateways.networking.istio.io bookinfo-gateway
```

Crie um port-forward para acessar a aplicação:

```bash
# Port forward in background for access Productpage application
# Allow port 9080/TCP
kubectl -n $MY_NAMESPACE port-forward svc/productpage 9080:9080 --address=0.0.0.0
# Access URL: http://localhost:9080/productpage
```

Instale os seguintes addons para o Istio: [Grafana](https://grafana.com), [Prometheus](https://prometheus.io), [Kiali](https://kiali.io), [Jaeger](https://www.jaegertracing.io), [zipkin](https://zipkin.io/).

```bash
git clone https://github.com/istio/istio

kubectl apply -f $HOME/istio/samples/addons

kubectl rollout status deployment/kiali -n istio-system
kubectl rollout status deployment/grafana -n istio-system
kubectl rollout status deployment/prometheus -n istio-system
kubectl rollout status deployment/jaeger -n istio-system
```

Inicie um port-forward para cada addon.

```bash
kubectl -n istio-system port-forward service/kiali 20001:20001
kubectl -n istio-system port-forward service/grafana 3000:3000
kubectl -n istio-system port-forward service/prometheus 9090:9090
kubectl -n istio-system port-forward service/tracing 8080:80
```

Acesse cada addon nos seguintes endereços:

* Kiali: http://localhost:20001
* Grafana: http://localhost:3000
* Prometheus: http://localhost:9090
* Jaeger: http://localhost:8080

Liste os objetos do Istio e addons com os seguintes comandos:

```bash
kubectl api-resources | grep istio
kubectl get crd | grep istio
kubectl get all -n istio-system
kubectl get all -n istio-ingress
```

Permita que o Istio gerencie as aplicações de determinado namespace:

```bash
kubectl label namespace $MY_NAMESPACE istio-injection=enabled
```
