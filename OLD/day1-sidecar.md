# Descomplicando Istio

<!-- TOC -->

- [Descomplicando Istio](#descomplicando-istio)
- [Day 1](#day-1)
- [Customizando logs de acesso no Istio Gateway para incluir a métrica REQUEST\_TX\_DURATION](#customizando-logs-de-acesso-no-istio-gateway-para-incluir-a-métrica-request_tx_duration)
- [Referências](#referências)

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

# Configure o repositório Helm (usando o sidecar mode):
# https://istio.io/v1.27/docs/setup/install/helm/
helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update

# Instale os componentes base do Istio
helm -n istio-system install istio-base istio/base --version $VERSION_ISTIO --set defaultRevision=default --create-namespace --wait --debug --timeout 900s

# Instale ou atualize o Kubernetes Gateway API CRDs
kubectl get crd gateways.gateway.networking.k8s.io &> /dev/null || \
  kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.3.0/standard-install.yaml

# Instale o CNI node agent. Ele é responsável por detectar os pods que pertencem  ao ambiente mesh, e configura o encaminhamento de tráfego entre os pods e o ztunnel node proxy (que será instalado mais adiante).
helm -n istio-system install istio-cni istio/cni --version $VERSION_ISTIO --wait --debug --timeout 900s

# Instale o Istiod, o componente control plane que gerencia e configura os proxies para roteamento de tráfego na mesh
helm -n istio-system install istiod istio/istiod --version $VERSION_ISTIO --wait --debug --timeout 900s
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
# Instale o Ingress gateway
helm -n istio-system install istio-ingress istio/gateway --version $VERSION_ISTIO --wait --debug --timeout 900s
```

> ATENÇÃO!!! Se durante a instalação do Istio Gateway você encontrar o seguinte erro:
>
> ``Error: INSTALLATION FAILED: values don't meet the specifications of the schema(s) in the following chart(s): gateway: - at '': additional properties '_internal_defaults_do_not_set' not allowed``
>
> Isso foi reportado neste Issue: https://github.com/istio/istio/issues/57354 e tem um fix neste pull request: https://github.com/istio/istio/pull/57426
>
> Outra opção é adicionar a flag ``--skip-schema-validation`` no comando helm acima.

Valide a instalação dos componentes do Istio com os seguintes comandos:

```bash
helm -n istio-system ls
helm -n istio-system status istio-base
helm -n istio-system status istiod
helm -n istio-system status istio-cni
helm -n istio-system status istio-ingress
kubectl -n istio-system get all --output wide
```

Faça o deploy da aplicação de exemplo chamada **Bookinfo**.

> Documentação de referência para os arquivos e comandos mostrados a seguir: https://istio.io/v1.27/docs/examples/bookinfo/

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

Por padrão, o Istio cria um serviço LoadBalancer para um gateway. Para verificar o status do gateway, execute o seguinte comando:

```bash
kubectl -n $MY_NAMESPACE get gateway
```

> ATENÇÃO!!! Se na coluna ADDRESS do gateway você ver <pending>, aguarde alguns instantes até que o MetalLB atribua um endereço IP ao serviço LoadBalancer.
> Na coluna PROGRAMMED você deve ver o valor TRUE, indicando que o gateway está pronto para receber tráfego.

Use o comando abaixo para obter o IP do gateway.

```bash
export GATEWAY_IP=$(kubectl -n $MY_NAMESPACE get gateway -o jsonpath='{.items[*].status.addresses[*].value}')
echo $GATEWAY_IP
```

Use o navegador, acesse a página http://GATEWAY_IP:80/productpage.

> ATENÇÃO!!! Substitua GATEWAY_IP pelo valor obtido no comando anterior.
> O gateway faz uso da porta 80/TCP para tráfego HTTP.
> Você deve ver a página do produto Bookinfo.

<p align="center">
  <img src="images/bookinfo.png" alt="Bookinfo productpage">
</p>

> ATENÇÃO!!! Se você estiver usando **MacOS** com o kind criado no Docker Desktop, será necessário instalar o serviço https://github.com/chipmk/docker-mac-net-connect e possívelmente configurar o fix sugerido neste issue: https://github.com/chipmk/docker-mac-net-connect/issues/62
> Comandos para reiniciar o serviço docker-mac-net-connect:

```bash
sudo brew services stop docker-mac-net-connect
sudo brew services start docker-mac-net-connect
sudo brew services info docker-mac-net-connect
```

Outra opção é criar um port-forward para acessar a página do produto Bookinfo.

```bash
kubectl -n $MY_NAMESPACE port-forward svc/bookinfo-gateway-istio 8080:80
```

Use o navegador, acesse a página http://localhost:8080/productpage.

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

Permita que o Istio gerencie as aplicações de determinado namespace (Sidecar mode https://istio.io/v1.27/docs/setup/getting-started/). O comando abaixo adiciona o label ``istio-injection=enabled`` em determinado namespace, permitindo que o Istio injete os sidecars automaticamente nos pods criados nesse namespace.

```bash
kubectl label namespace $MY_NAMESPACE istio-injection=enabled
```

Depois disse execute o comando abaixo para reiniciar os pods da aplicação Bookinfo, para que os sidecars sejam injetados.

```bash
kubectl -n $MY_NAMESPACE rollout restart deployment
kubectl -n $MY_NAMESPACE get pods
```

Envie tráfego para o aplicativo Bookinfo, para que o Kiali gere o gráfico de tráfego:

```bash
# Abordagem 1: Usando o IP do gateway
for i in $(seq 1 10000); do curl -sSI -o /dev/null http://${GATEWAY_IP}/productpage; done

# Abordagem 2: Usando o endereço da aplicação via port-forward
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

# Customizando logs de acesso no Istio Gateway para incluir a métrica REQUEST_TX_DURATION

Para incluir logs específicos do Envoy, como ``%REQUEST_TX_DURATION%`` (que mede o tempo decorrido desde o primeiro byte da solicitação recebida até o último byte enviado), você deve atualizar o ``meshConfig``.

```bash
# Configure o log do envoy para enviar logs no formato JSON
cat <<EOF > istiod-values.yaml
# Approach 1: Configure access logs in JSON format with custom metric 'REQUEST_TX_DURATION' for all Istio proxies
#
# Istio Operator values to enable JSON access logs with custom metric
#global:
#  logAsJson: true
#
# Configure access logs in JSON format with custom metric 'REQUEST_TX_DURATION'
#meshConfig:
#  accessLogFile: /dev/stdout
#  extensionProviders:
#    - name: "file-log"
#      envoyFileAccessLog:
#        path: /dev/stdout
#        logFormat:
#          # Use 'labels' to generate JSON output
#          labels:
#            timestamp: "%START_TIME%"
#            method: "%REQ(:METHOD)%"
#            path: "%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%"
#            protocol: "%PROTOCOL%"
#            response_code: "%RESPONSE_CODE%"
#            response_code_details: "%RESPONSE_CODE_DETAILS(X)%"
#            response_flags: "%RESPONSE_FLAGS%"
#            bytes_received: "%BYTES_RECEIVED%"
#            bytes_sent: "%BYTES_SENT%"
#            duration: "%DURATION%"
#            upstream_service_time: "%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%"
#            x_forwarded_for: "%REQ(X-FORWARDED-FOR)%"
#            user_agent: "%REQ(USER-AGENT)%"
#            request_id: "%REQ(X-REQUEST-ID)%"
#            authority: "%REQ(:AUTHORITY)%"
#            upstream_host: "%UPSTREAM_HOST%"
#            # The custom metric you want
#            request_tx_duration: "%REQUEST_TX_DURATION%"
#  defaultProviders:
#    accessLogging:
#      - "file-log"

#------------------------------------------------------------------------------
# # Approach 2: Configure access logs in JSON format with custom metric 'REQUEST_TX_DURATION' only in Gateways
# Sidecars will use standard config unless 'defaultProviders' is uncommented below
meshConfig:
  accessLogFile: /dev/stdout
  extensionProviders:
    - name: "json-gateway-log"
      envoyFileAccessLog:
        path: /dev/stdout
        logFormat:
          labels:
            timestamp: "%START_TIME%"
            method: "%REQ(:METHOD)%"
            path: "%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%"
            protocol: "%PROTOCOL%"
            response_code: "%RESPONSE_CODE%"
            response_code_details: "%RESPONSE_CODE_DETAILS(X)%"
            response_flags: "%RESPONSE_FLAGS%"
            bytes_received: "%BYTES_RECEIVED%"
            bytes_sent: "%BYTES_SENT%"
            duration: "%DURATION%"
            upstream_service_time: "%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%"
            x_forwarded_for: "%REQ(X-FORWARDED-FOR)%"
            user_agent: "%REQ(USER-AGENT)%"
            request_id: "%REQ(X-REQUEST-ID)%"
            authority: "%REQ(:AUTHORITY)%"
            upstream_host: "%UPSTREAM_HOST%"
            # The custom metric you want
            request_tx_duration: "%REQUEST_TX_DURATION%"
  # Remove 'defaultProviders' so sidecars revert to the standard config above.
  # defaultProviders: 
  #   accessLogging: []

#------------------------------------------------------------------------------
# Approach 3: Configure access logs in Text format with custom metric 'REQUEST_TX_DURATION'
# For Datadog its recommended to use JSON format for better parsing
#meshConfig:
#  accessLogFile: /dev/stdout
#  # Define the provider explicitly
#  extensionProviders:
#    - name: "file-log"
#      envoyFileAccessLog:
#        path: /dev/stdout
#        logFormat:
#          text: |
#            [%START_TIME%] "%REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% %PROTOCOL%" %RESPONSE_CODE% %RESPONSE_CODE_DETAILS(X)% %RESPONSE_FLAGS% %BYTES_RECEIVED% %BYTES_SENT% %DURATION% %RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)% "%REQ(X-FORWARDED-FOR)%" "%REQ(USER-AGENT)%" "%REQ(X-REQUEST-ID)%" "%REQ(:AUTHORITY)%" "%UPSTREAM_HOST%" %REQUEST_TX_DURATION%
#  # Tell Istio to use this provider for all traffic
#  defaultProviders:
#    accessLogging:
#      - "file-log"
EOF

# Atualize o Istiod com o arquivo de configuração anterior
helm -n istio-system upgrade --install istiod istio/istiod --version $VERSION_ISTIO -f istiod-values.yaml --wait --debug --timeout 900s
```

Isso será usado para criar um configmap istio com a configuração do log atualizado

Usando a **abordagem 2**, os logs personalizados ainda não serão aplicados em todos os containers do envoy. Para aplicar essa configuração apenas no Istio Gateway (e não em todos os sidecars), você pode criar um recurso Telemetry específico para o gateway, conforme mostrado abaixo:

```bash
cat <<EOF > telemetry-log-gateway.yaml
apiVersion: telemetry.istio.io/v1
kind: Telemetry
metadata:
  name: ingress-json-logging
  namespace: myapp  # typically where ingress-gateway runs
spec:
  # 1. SELECTOR: This targets only the Ingress Gateway pods
  selector:
    matchLabels:
      gateway.istio.io/managed: istio.io-gateway-controller # Check your specific label if different

  # 2. OVERRIDE: Use the specific provider defined in MeshConfig
  accessLogging:
    - providers:
      - name: json-gateway-log
EOF

kubectl apply -f telemetry-log-gateway.yaml
```

> ATENÇÃO!!! O campo `selector` no objeto Telemetry é usado para especificar a quais pods a configuração de log deve ser aplicada. Neste caso, estamos direcionando apenas os pods do Ingress Gateway, correspondendo aos seus labels. Você pode executar o seguinte comando para inspecionar os labels dos pods do Ingress Gateway de determinado namespace:

```bash
kubectl -n NAMESPACE get pods --show-labels
```

**Validando a configuração do log do Envoy**:

Não verifique o arquivo ``/etc/envoy/proxy/envoy_rev.json`` do container ``istio-proxy``. Em vez disso, consulte o proxy Envoy para saber qual é a sua configuração ativa atual usando o istioctl. A configuração do log pode ser substituída dinamicamente, portanto, o arquivo no sistema de arquivos pode não refletir a configuração real em uso.

```bash
istioctl proxy-config listener pod/POD_NAME -n myapp --port 9080 -o json | grep -C 5 "REQUEST_TX_DURATION"
```

> ATENÇÃO!!! Substitua 9080 pela porta de serviço do seu aplicativo.

Se você encontrar `%REQUEST_TX_DURATION%` aqui, significa que está configurado.

Gere tráfego e verifique os logs: Envie uma solicitação para seu aplicativo e observe os logs do proxy Istio para verificar se o campo `%REQUEST_TX_DURATION%` está sendo registrado corretamente.

Neste exemplo, o log do Envoy deve se parecer com o seguinte:

```bash
# [2025-12-29T...] "GET /productpage ..." 200 - 0 123 5 4 "-" "curl/7.8" ... "10.244.0.5:9080" 
```

O último número (por exemplo, 5) é a duração da sua transação (TX) em milisegundos.

# Referências

- https://www.envoyproxy.io/docs/envoy/latest/configuration/advanced/substitution_formatter#config-advanced-substitution-operators
- https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/
- https://gateway.envoyproxy.io/v1.5/tasks/observability/proxy-accesslog/
- https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage
- https://istio.io/latest/docs/tasks/observability/logs/access-log/
- https://dev.to/aws-builders/understanding-istio-access-logs-2k5o
