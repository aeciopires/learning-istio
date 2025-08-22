# Descomplicando Istio

<!-- TOC -->

- [Descomplicando Istio](#descomplicando-istio)
- [Day 1](#day-1)

<!-- TOC -->

# Day 1

Comandos executados em sequência durante o treinamento.

Instale o **Docker** com as instruções da página: [INSTALL_REQUIREMENTS.md#docker](INSTALL_REQUIREMENTS.md#docker).

Crie um **cluster Kubernetes** usando o kind conforme as instruções da página: [INSTALL_REQUIREMENTS.md#kind](INSTALL_REQUIREMENTS.md#kind);

Instale o **Helm** com as instruções da página: [INSTALL_REQUIREMENTS.md#helm](INSTALL_REQUIREMENTS.md#helm).

Instale o **MetalLB** com as instruções da página: [INSTALL_REQUIREMENTS.md#metallb](INSTALL_REQUIREMENTS.md#metallb).

Crie variáveis de ambiente uteis para baixar os arquivos complementares

```bash
export ISTIO_RELEASE=1.27
export VERSION_ISTIO="${ISTIO_RELEASE}.0"
export ISTIO_BASE_URL="https://raw.githubusercontent.com/istio/istio/release-$ISTIO_RELEASE/samples/"
export ISTIO_BOOKINFO_URL="$ISTIO_BASE_URL/bookinfo/"
export ISTIO_ADDONS_URL="$ISTIO_BASE_URL/addons"
export MY_NAMESPACE='myapp'
export ISTIO_HTTPBIN_URL="$ISTIO_BASE_URL/httpbin/"
```

Instale o **Istio** com os seguintes comandos:

```bash
#-------
# Referencias:
# https://istio.io/v1.27/docs/setup/getting-started/
# https://istio.io/v1.27/docs/setup/install/helm/
# Liberar todas essas portas no firewall: https://istio.io/latest/docs/ops/deployment/application-requirements/

# Configure o repositório Helm (usando o ambient mode):
# https://istio.io/v1.27/docs/ambient/install/helm/
helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update

# Instale os componentes base do Istio
helm -n istio-system install istio-base istio/base --version $VERSION_ISTIO --create-namespace --wait --debug --timeout 900s

# Instale ou atualize o Kubernetes Gateway API CRDs
kubectl get crd gateways.gateway.networking.k8s.io &> /dev/null || \
  kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.3.0/standard-install.yaml

# Instale o Istiod, o componente control plane que gerencia e configura os proxies para roteamento de tráfego na mesh
helm -n istio-system install istiod istio/istiod --version $VERSION_ISTIO --set profile=ambient --wait --debug --timeout 900s

# Instale o CNI node agent. Ele é responsável por detectar os pods que pertencem  ao ambiente mesh, e configura o encaminhamento de tráfego entre os pods e o ztunnel node proxy (que será instalado mais adiante).
helm -n istio-system install istio-cni istio/cni --version $VERSION_ISTIO --set profile=ambient --wait --debug --timeout 900s
```

> ATENÇÃO!!! Se você ver o erro abaixo enquanto verifica o log do istio-cni: 
``failed to create fsnotify watcher: too many open files``, corrija com o seguintes comandos:

```bash
sudo sysctl -w fs.inotify.max_user_watches=2099999999
sudo sysctl -w fs.inotify.max_user_instances=2099999999
sudo sysctl -w fs.inotify.max_queued_events=2099999999
```

Referência: https://serverfault.com/questions/1137211/failed-to-create-fsnotify-watcher-too-many-open-files

Continuação da instalação dos componentes do Istio com os seguintes comandos:

```bash
# Instale o ztunnel DaemonSet, que é o componente daemonset do node proxy do Istio no ambient mode.
helm -n istio-system install ztunnel istio/ztunnel --version $VERSION_ISTIO --wait --debug --timeout 900s

# Instale o Ingress gateway
helm -n istio-system install istio-ingress istio/gateway --version $VERSION_ISTIO --wait --debug --timeout 900s
```

> ATENÇÃO!!! Se durante a instalação do Istio Gateway você encontrar o seguinte erro:
>
> ``Error: INSTALLATION FAILED: values don't meet the specifications of the schema(s) in the following chart(s): gateway: - at '': additional properties '_internal_defaults_do_not_set' not allowed``
>
> Isso foi reportado neste Issue: https://github.com/istio/istio/issues/57354 e tem um fix neste pull request: https://github.com/istio/istio/pull/57426

Valide a instalação dos componentes do Istio com os seguintes comandos:

```bash
helm -n istio-system ls
helm -n istio-system status istio-base
helm -n istio-system status istiod
helm -n istio-system status istio-cni
helm -n istio-system status ztunnel
helm -n istio-system status istio-ingress
kubectl -n istio-system get all --output wide
```

Faça o deploy da aplicação de exemplo chamada **Bookinfo**.

> Documentação de referência para os arquivos e comandos mostrados a seguir: https://istio.io/v1.27/docs/ambient/getting-started/deploy-sample-app

Instale a aplicação de exemplo:

```bash
kubectl create namespace $MY_NAMESPACE

kubectl -n $MY_NAMESPACE apply -f "$ISTIO_BOOKINFO_URL/platform/kube/bookinfo.yaml"
kubectl -n $MY_NAMESPACE apply -f "$ISTIO_BOOKINFO_URL/platform/kube/bookinfo-versions.yaml"
```

Visualize os objetos/recursos da aplicação:

```bash
kubectl -n $MY_NAMESPACE get all 
```

Teste o acesso a aplicação com o seguinte comando:

```bash
kubectl -n $MY_NAMESPACE exec "$(kubectl -n $MY_NAMESPACE get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl -sS productpage:9080/productpage | grep -o "<title>.*</title>"
```

Crie um ingressGateway e um HTTPRoute para a aplicação ser exposta fora do cluster usando o Istio e o loadbalancer.

```bash
kubectl -n $MY_NAMESPACE apply -f "$ISTIO_BOOKINFO_URL/gateway-api/bookinfo-gateway.yaml"
```

Por padrão, o Istio cria um serviço LoadBalancer para um gateway. Como você acessará esse gateway por um túnel, não precisa de um balanceador de carga. Altere o tipo de serviço para ClusterIP na anotation do gateway:

```bash
kubectl -n $MY_NAMESPACE annotate gateway bookinfo-gateway networking.istio.io/service-type=ClusterIP 
```

Para verificar o status do gateway, execute o seguinte comando:

```bash
kubectl -n $MY_NAMESPACE get gateway
```

Crie um port-forward para acessar a página do produto Bookinfo por meio do gateway que acabou de provisionar.

```bash
kubectl -n $MY_NAMESPACE port-forward svc/bookinfo-gateway-istio 8080:80
```

Usando o navegador, acesse a página http://localhost:8080/productpage.

<p align="center">
  <img src="images/bookinfo.png" alt="Bookinfo productpage">
</p>

Visualize os objetos da aplicação Bookinfo.

```bash
kubectl -n $MY_NAMESPACE get gateway,httproute,service,pods,deployments,replicaset
```

Instale os seguintes addons para o Istio:

- [Grafana](https://grafana.com)
- [Prometheus](https://prometheus.io)
- [Kiali](https://kiali.io)
- [Jaeger](https://www.jaegertracing.io)

```bash
# Instalando os addons
kubectl apply -f "$ISTIO_ADDONS_URL/grafana.yaml"
kubectl apply -f "$ISTIO_ADDONS_URL/prometheus.yaml"
kubectl apply -f "$ISTIO_ADDONS_URL/kiali.yaml"
kubectl apply -f "$ISTIO_ADDONS_URL/jaeger.yaml"

# Verificando o status da instalação dos addons
kubectl rollout status deployment/grafana -n istio-system
kubectl rollout status deployment/prometheus -n istio-system
kubectl rollout status deployment/kiali -n istio-system
kubectl rollout status deployment/jaeger -n istio-system
```

Inicie um port-forward para cada addon.

```bash
kubectl -n istio-system port-forward service/grafana 3000:3000
kubectl -n istio-system port-forward service/prometheus 9090:9090
kubectl -n istio-system port-forward service/kiali 20001:20001
kubectl -n istio-system port-forward service/tracing 8081:80
```

Acesse cada addon nos seguintes endereços:

- Grafana: http://localhost:3000 (login: admin, senha: admin)
- Prometheus: http://localhost:9090
- Kiali: http://localhost:20001
- Jaeger: http://localhost:8081

Permita que o Istio gerencie as aplicações de determinado namespace (Ambient mode https://istio.io/v1.27/docs/ambient/getting-started/secure-and-visualize/):

```bash
kubectl label namespace $MY_NAMESPACE istio.io/dataplane-mode=ambient
```

Envie tráfego para o aplicativo Bookinfo, para que o Kiali gere o gráfico de tráfego:

```bash
for i in $(seq 1 10000); do curl -sSI -o /dev/null http://localhost:8080/productpage; done
```

Veja o resultado conforme mostrado na imagem a seguir:

<p align="center">
  <img src="images/bookinfo_kiali.png" alt="Bookinfo visualized on Kiali">
</p>

Liste os objetos do Istio e addons com os seguintes comandos:

```bash
kubectl api-resources | grep istio
kubectl get crd | grep istio
kubectl get all -n istio-system
kubectl get all -n istio-ingress
```

Para desinstalar o Istio, execute os comandos da página: [UNINSTALL_ISTIO.md](UNINSTALL_ISTIO.md).
