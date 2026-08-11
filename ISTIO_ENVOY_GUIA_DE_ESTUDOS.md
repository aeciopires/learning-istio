# Guia de Estudos: Istio e Envoy

<!-- TOC -->

- [Guia de Estudos: Istio e Envoy](#guia-de-estudos-istio-e-envoy)
- [Sobre este material](#sobre-este-material)
- [1. Glossário](#1-glossário)
  - [1.1. Conceitos gerais de Service Mesh](#11-conceitos-gerais-de-service-mesh)
  - [1.2. Glossário do Istio](#12-glossário-do-istio)
  - [1.3. Glossário do Envoy](#13-glossário-do-envoy)
  - [1.4. Glossário de infraestrutura usada neste repositório](#14-glossário-de-infraestrutura-usada-neste-repositório)
- [2. Arquitetura do Istio](#2-arquitetura-do-istio)
  - [2.1. Modo Sidecar](#21-modo-sidecar)
  - [2.2. Modo Ambient](#22-modo-ambient)
  - [2.3. Sidecar vs Ambient (comparação)](#23-sidecar-vs-ambient-comparação)
  - [2.4. Portas usadas pelo Istio e pelo Envoy](#24-portas-usadas-pelo-istio-e-pelo-envoy)
  - [2.5. Componentes de software do Istio](#25-componentes-de-software-do-istio)
  - [2.6. Topologias de implantação da mesh (multicluster, multi-network, multi-mesh)](#26-topologias-de-implantação-da-mesh-multicluster-multi-network-multi-mesh)
- [3. Envoy Proxy](#3-envoy-proxy)
  - [3.1. Visão geral e histórico](#31-visão-geral-e-histórico)
  - [3.2. Arquitetura interna (listener subsystem x cluster subsystem)](#32-arquitetura-interna-listener-subsystem-x-cluster-subsystem)
  - [3.3. xDS: como o Istio configura o Envoy](#33-xds-como-o-istio-configura-o-envoy)
  - [3.4. Life of a Request (o ciclo de vida de uma requisição no Envoy)](#34-life-of-a-request-o-ciclo-de-vida-de-uma-requisição-no-envoy)
  - [3.5. Vídeos recomendados de Matt Klein (criador do Envoy)](#35-vídeos-recomendados-de-matt-klein-criador-do-envoy)
  - [3.6. Componentes de software do Envoy: threading model e hot restart](#36-componentes-de-software-do-envoy-threading-model-e-hot-restart)
  - [3.7. A cadeia de filtros HTTP em detalhe: tipos, ordem e controle de fluxo](#37-a-cadeia-de-filtros-http-em-detalhe-tipos-ordem-e-controle-de-fluxo)
  - [3.8. Comandos de troubleshooting da cadeia de filtros](#38-comandos-de-troubleshooting-da-cadeia-de-filtros)
- [4. Comandos práticos (kind + Istio deste repositório)](#4-comandos-práticos-kind--istio-deste-repositório)
  - [4.1. Pré-requisitos](#41-pré-requisitos)
  - [4.2. istioctl: comandos essenciais](#42-istioctl-comandos-essenciais)
  - [4.3. Admin API do Envoy (porta 15000)](#43-admin-api-do-envoy-porta-15000)
  - [4.4. Comandos específicos do modo Ambient (ztunnel e waypoint)](#44-comandos-específicos-do-modo-ambient-ztunnel-e-waypoint)
  - [4.5. Comandos avançados de troubleshooting do Envoy](#45-comandos-avançados-de-troubleshooting-do-envoy)
- [5. Observabilidade: métricas com Prometheus, Grafana e Datadog](#5-observabilidade-métricas-com-prometheus-grafana-e-datadog)
  - [5.1. Métricas padrão do Istio](#51-métricas-padrão-do-istio)
  - [5.2. Estatísticas internas do Envoy](#52-estatísticas-internas-do-envoy)
  - [5.3. Métricas do control plane (istiod) e tempo de convergência](#53-métricas-do-control-plane-istiod-e-tempo-de-convergência)
  - [5.4. Response Flags do Envoy](#54-response-flags-do-envoy)
  - [5.5. Prometheus + Grafana](#55-prometheus--grafana)
  - [5.6. Datadog](#56-datadog)
  - [5.7. Como detectar erros comuns através das métricas](#57-como-detectar-erros-comuns-através-das-métricas)
- [6. Logs do Istio e do Envoy](#6-logs-do-istio-e-do-envoy)
  - [6.1. Nível de log do Envoy/istio-proxy (componentes e debug)](#61-nível-de-log-do-envoyistio-proxy-componentes-e-debug)
  - [6.2. Nível de log do istiod](#62-nível-de-log-do-istiod)
  - [6.3. Access log: habilitar e customizar campos](#63-access-log-habilitar-e-customizar-campos)
  - [6.4. Latência e outras métricas visíveis apenas no log](#64-latência-e-outras-métricas-visíveis-apenas-no-log)
  - [6.5. Comandos de troubleshooting usando logs](#65-comandos-de-troubleshooting-usando-logs)
- [7. EnvoyFilter](#7-envoyfilter)
  - [7.1. O que é e quando usar](#71-o-que-é-e-quando-usar)
  - [7.2. Estrutura do objeto: workloadSelector, configPatches e match](#72-estrutura-do-objeto-workloadselector-configpatches-e-match)
  - [7.3. `applyTo`: em qual camada da configuração o patch age](#73-applyto-em-qual-camada-da-configuração-o-patch-age)
  - [7.4. `patch.operation`: como a requisição/configuração é manipulada](#74-patchoperation-como-a-requisiçãoconfiguração-é-manipulada)
  - [7.5. Camadas de processamento e onde o EnvoyFilter se encaixa](#75-camadas-de-processamento-e-onde-o-envoyfilter-se-encaixa)
  - [7.6. Exemplo prático: habilitando compressão gzip via EnvoyFilter](#76-exemplo-prático-habilitando-compressão-gzip-via-envoyfilter)
  - [7.7. Comandos de troubleshooting de EnvoyFilter](#77-comandos-de-troubleshooting-de-envoyfilter)
- [8. Configurações avançadas do Istio](#8-configurações-avançadas-do-istio)
  - [8.1. Telemetry: habilitando/customizando métricas por workload](#81-telemetry-habilitandocustomizando-métricas-por-workload)
  - [8.2. MeshConfig: configurações globais da malha](#82-meshconfig-configurações-globais-da-malha)
  - [8.3. Certificados mTLS: plugando uma CA própria e rotação](#83-certificados-mtls-plugando-uma-ca-própria-e-rotação)
  - [8.4. Sidecar: restringindo com quais endereços um workload se relaciona](#84-sidecar-restringindo-com-quais-endereços-um-workload-se-relaciona)
  - [8.5. AuthorizationPolicy: autorizando apenas aplicações específicas](#85-authorizationpolicy-autorizando-apenas-aplicações-específicas)
  - [8.6. Autorizando por JWT e/ou mTLS (RequestAuthentication + AuthorizationPolicy)](#86-autorizando-por-jwt-eou-mtls-requestauthentication--authorizationpolicy)
- [9. Exercícios de troubleshooting](#9-exercícios-de-troubleshooting)
  - [9.0. Preparando o ambiente dos laboratórios](#90-preparando-o-ambiente-dos-laboratórios)
  - [Lab 1 - Pod fora da mesh (ambient)](#lab-1---pod-fora-da-mesh-ambient)
  - [Lab 2 - mTLS STRICT quebrando tráfego (sidecar)](#lab-2---mtls-strict-quebrando-tráfego-sidecar)
  - [Lab 3 - VirtualService/DestinationRule com host errado (503 / no healthy upstream)](#lab-3---virtualservicedestinationrule-com-host-errado-503--no-healthy-upstream)
  - [Lab 4 - istiod fora de sincronia (STALE / NOT SENT)](#lab-4---istiod-fora-de-sincronia-stale--not-sent)
  - [Lab 5 - Inspecionando o Envoy diretamente pela Admin API](#lab-5---inspecionando-o-envoy-diretamente-pela-admin-api)
  - [Lab 6 - Política L7 ignorada em ambient mode sem waypoint](#lab-6---política-l7-ignorada-em-ambient-mode-sem-waypoint)
  - [Lab 7 - Diagnosticando o ztunnel (ambient)](#lab-7---diagnosticando-o-ztunnel-ambient)
  - [Lab 8 - `istioctl analyze` encontrando erros de configuração](#lab-8---istioctl-analyze-encontrando-erros-de-configuração)
  - [Lab 9 - Fault injection e timeout (retomando o Day 2)](#lab-9---fault-injection-e-timeout-retomando-o-day-2)
  - [Lab 10 - Circuit breaking e overflow detectados via métricas](#lab-10---circuit-breaking-e-overflow-detectados-via-métricas)
  - [Lab 11 - Medindo o tempo de convergência do istiod](#lab-11---medindo-o-tempo-de-convergência-do-istiod)
  - [Lab 12 - Compressão gzip via EnvoyFilter não aplicada](#lab-12---compressão-gzip-via-envoyfilter-não-aplicada)
  - [Lab 13 - AuthorizationPolicy negando (ou liberando) por engano](#lab-13---authorizationpolicy-negando-ou-liberando-por-engano)
  - [Lab 14 - JWT inválido/ausente sendo rejeitado (ou não)](#lab-14---jwt-inválidoausente-sendo-rejeitado-ou-não)
- [10. Referências](#10-referências)

<!-- TOC -->

# Sobre este material

Este documento é um **material de estudo** (glossário + diagramas + comandos + exercícios de troubleshooting) sobre **Istio** e **Envoy**, complementar aos arquivos [day1.md](day1.md) e [day2.md](day2.md) deste repositório.

Todo o conteúdo técnico aqui apresentado foi extraído/conferido nas fontes oficiais listadas na seção [10. Referências](#10-referências) (documentação oficial do Istio, documentação oficial do Envoy e conteúdos publicados por Matt Klein, criador do Envoy). Nenhum comportamento interno foi inventado: onde a documentação oficial não detalha algo (ex.: código-fonte específico), isso é indicado explicitamente.

**Versão de referência:** Istio **1.30** (LTS, versão estável mais recente em 2026-08-10) e Envoy na versão empacotada por ela. A versão 1.31 do Istio está em status `alpha` (`1.31.0-alpha.0`) na data deste material, portanto **não é recomendada** como base de estudo/produção — por isso este guia segue o padrão já usado no restante do repositório (`ISTIO_RELEASE=1.30`), conforme confirmado em [Istio / Supported Releases](https://istio.io/latest/docs/releases/supported-releases/) e nas releases do [github.com/istio/istio](https://github.com/istio/istio/releases/).

O ambiente prático usado neste guia é o mesmo do restante do repositório: cluster **kind** + **MetalLB** + **Istio** (modo `ambient` por padrão, com o modo `sidecar` como alternativa documentada em [OLD/day1-sidecar.md](OLD/day1-sidecar.md)). Veja os pré-requisitos completos em [INSTALL_REQUIREMENTS.md](INSTALL_REQUIREMENTS.md).

---

# 1. Glossário

## 1.1. Conceitos gerais de Service Mesh

| Termo | Definição |
|---|---|
| **Service Mesh** | Camada de infraestrutura dedicada que torna a comunicação entre serviços segura, rápida e confiável, tipicamente implementada com proxies leves implantados junto às aplicações, sem que a aplicação tenha conhecimento disso. |
| **Data plane** | Conjunto de proxies (no Istio, o Envoy) que efetivamente intermediam/transportam o tráfego de rede entre as cargas de trabalho (workloads) e coletam telemetria. |
| **Control plane** | Componentes responsáveis por gerenciar e configurar os proxies do data plane, transformando a configuração declarada pelo usuário (CRDs) em configuração de proxy (xDS). |
| **Sidecar** | Padrão de implantação em que um container auxiliar (o proxy) roda dentro do mesmo Pod da aplicação, compartilhando o ciclo de vida e a rede (`localhost`) com o container principal. |
| **mTLS (mutual TLS)** | Autenticação TLS em que **ambas** as partes (cliente e servidor) apresentam certificados, usada pelo Istio para autenticar e criptografar o tráfego entre workloads na mesh. |

Fonte: [Istio / What is Istio?](https://istio.io/v1.30/docs/overview/what-is-istio/), [Istio / Glossary](https://istio.io/v1.30/docs/reference/glossary/).

## 1.2. Glossário do Istio

Termos extraídos do glossário oficial ([istio.io/v1.30/docs/reference/glossary/](https://istio.io/v1.30/docs/reference/glossary/)):

| Termo | Definição (oficial, resumida) |
|---|---|
| **Istiod** | Binário de control plane consolidado que **encapsula** as antigas funções separadas **Pilot** (descoberta de serviço/configuração do Envoy), **Galley** (validação/distribuição de configuração) e **Citadel** (Autoridade Certificadora/CA), desde o Istio 1.5. |
| **Pilot** | Subcomponente do Istiod responsável por traduzir as regras de roteamento do usuário (VirtualService, DestinationRule, Gateway, etc.) em configuração específica do Envoy (xDS) e distribuí-la aos proxies. |
| **Ambient** | Modo de data plane **sem sidecar**: um componente por nó (ztunnel) e, opcionalmente, por namespace/serviço (waypoint), em vez de um proxy por Pod. |
| **Data Plane Mode** | A forma de implantação do proxy: `sidecar`, `ambient` ou `proxyless` (gRPC sem proxy, apenas com biblioteca xDS). |
| **Secure L4 Overlay** | A camada L4 (mTLS, autenticação, autorização, telemetria básica) fornecida pelo `ztunnel` no modo ambient. |
| **Waypoint** | Proxy Envoy **opcional**, L7, implantado como Deployment (não como sidecar) para fornecer o conjunto completo de recursos do Istio (roteamento HTTP, retries, fault injection, políticas L7) a workloads no modo ambient. |
| **ztunnel** ("zero-trust tunnel") | Componente node proxy do modo ambient; escrito para ser leve, escopo restrito a L3/L4 (mTLS, autenticação, autorização L4 e telemetria básica). |
| **HBONE** | "HTTP-Based Overlay Network Environment": protocolo de tunelamento seguro (HTTP/2 CONNECT sobre mTLS) usado entre componentes do Istio no modo ambient (ex.: ztunnel-a-ztunnel, ztunnel-a-waypoint). |
| **Gateway** | Proxy standalone (Envoy) na borda da mesh, roteando tráfego de entrada (ingress) ou saída (egress). |
| **Gateway API** | API padrão do Kubernetes (`gateway.networking.k8s.io`) para configuração de roteamento de tráfego, suportada nativamente pelo Istio como alternativa aos CRDs `Gateway`/`VirtualService` clássicos do Istio. |
| **Sidecar Injection** | Webhook mutante do Kubernetes que insere automaticamente o container `istio-proxy` (Envoy) nos Pods de um namespace com o label `istio-injection=enabled` (ou por revisão). |
| **CNI (Container Network Interface)** | Padrão usado pelo `istio-cni` para configurar o redirecionamento de tráfego dos Pods sem precisar de um container `initContainer` privilegiado. |
| **Identity** | Identidade criptográfica (certificado x509 SPIFFE) atribuída a cada workload, base da autenticação mTLS. |
| **Trust Domain** | Raiz de confiança de um sistema; usada para formar as identidades de workload dentro da mesh (ex.: `cluster.local`). |
| **Proxyless** | Modo de data plane em que a aplicação (ex.: gRPC) fala xDS diretamente com o Istiod, sem passar por um proxy Envoy. |
| **Multicluster** | Modelo de implantação da mesh com mais de um cluster Kubernetes (cluster primário e clusters remotos). |
| **VirtualService** *(CRD, não está no glossário citado acima, mas é usado nos exemplos do repositório)* | Define regras de roteamento de tráfego (L7) para um ou mais hosts. |
| **DestinationRule** *(idem)* | Define políticas aplicadas ao tráfego **após** o roteamento (ex.: subsets/versões, política de TLS, circuit breaking, load balancing). |
| **PeerAuthentication** *(idem)* | Define a política de mTLS (`STRICT`, `PERMISSIVE`, `DISABLE`) para um namespace ou workload. |
| **AuthorizationPolicy** *(idem)* | Define regras de autorização (`ALLOW`/`DENY`) para tráfego entre workloads. |

Fonte: [Istio / Glossary](https://istio.io/v1.30/docs/reference/glossary/), [Istio / Architecture](https://istio.io/v1.30/docs/ops/deployment/architecture/), [Istio / Ambient Architecture](https://istio.io/v1.30/docs/ambient/architecture/), [Istio / Traffic Redirection](https://istio.io/v1.30/docs/ambient/architecture/traffic-redirection/).

## 1.3. Glossário do Envoy

Termos extraídos da documentação oficial ([envoyproxy.io/docs/envoy/latest/intro/life_of_a_request](https://www.envoyproxy.io/docs/envoy/latest/intro/life_of_a_request) e [envoyproxy.io/docs/envoy/latest/intro/arch_overview](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/arch_overview)):

| Termo | Definição (oficial, resumida) |
|---|---|
| **Listener** | "Módulo do Envoy responsável por vincular-se a um IP/porta, aceitar novas conexões TCP (ou datagramas UDP) e orquestrar os aspectos de processamento voltados ao downstream." |
| **Listener Filter** | Filtro que atua **antes** da conexão ser totalmente aceita/roteada — ex.: `TLS Inspector` (detecta SNI/ALPN para escolher a filter chain correta). |
| **Filter Chain Match** | Mecanismo que seleciona qual cadeia de filtros (transport socket + network filters) tratará a conexão, com base em critérios como SNI, porta de destino ou intervalo de IP de origem. |
| **Transport Socket** | Abstração de baixo nível responsável por ler/escrever bytes brutos na conexão (ex.: `TLS transport socket` faz o handshake e decripta/criptografa). |
| **Network Filter** | Filtro que atua no nível de conexão (L3/L4). Existem três tipos: `ReadFilter`, `WriteFilter` e `Filter` (ambos). O **HTTP Connection Manager (HCM)** é normalmente o último network filter da cadeia. |
| **HTTP Connection Manager (HCM)** | Network filter que faz o parsing do protocolo HTTP (HTTP/1.1, HTTP/2, HTTP/3), demultiplexa a conexão em *streams* e gerencia a cadeia de HTTP filters e o roteamento. |
| **HTTP Filter** | Filtro que atua no nível de request/response HTTP (ex.: CORS, JWT auth, RBAC, `router`). Existem `StreamDecoderFilter`, `StreamEncoderFilter` e `StreamFilter`. |
| **Router Filter** | HTTP filter especial, tipicamente o último da cadeia, que seleciona a rota (via tabela de rotas/RDS) e encaminha a requisição ao **cluster** de destino, conectando o *listener subsystem* ao *cluster subsystem*. |
| **Route / Route Table** | Configuração (obtida via RDS) que mapeia atributos da requisição (host, path, headers) a um cluster upstream. |
| **Cluster** | "Um serviço lógico com um conjunto de endpoints para os quais o Envoy encaminha requisições." Equivalente ao `DestinationRule`/serviço de destino no Istio. |
| **Endpoint** | Nó de rede (IP:porta) que implementa um serviço lógico; endpoints são agrupados em clusters e ficam "upstream" do Envoy. |
| **Cluster Manager** | Subsistema que gerencia todos os clusters conhecidos, seus endpoints, health checks, circuit breakers e connection pools. |
| **Load Balancer** | Componente do cluster manager que escolhe qual endpoint de um cluster receberá a próxima requisição (ex.: round robin, least request, ring hash, Maglev, random). |
| **Connection Pool** | Conjunto de conexões upstream reutilizáveis mantidas pelo Envoy para evitar o custo de abrir uma nova conexão por requisição. |
| **Circuit Breaker** | Mecanismo que limita conexões/requisições/retries pendentes por cluster para evitar sobrecarga em cascata. |
| **Outlier Detection** | Mecanismo de *health checking passivo* que remove temporariamente do load balancing endpoints que respondem com erros. |
| **Downstream** | Entidade que se conecta ao Envoy (ex.: a aplicação local, no caso de um listener de saída/*egress*, ou um outro nó da mesh, no caso de um listener de entrada/*ingress*). |
| **Upstream** | Endpoint para o qual o Envoy encaminha a requisição. |
| **xDS** | Família de APIs de descoberta usadas para configurar o Envoy dinamicamente (ver seção [3.3](#33-xds-como-o-istio-configura-o-envoy)). |
| **Admin interface** | Servidor HTTP embutido do Envoy (padrão porta configurável, `15000` no Istio) que expõe endpoints como `/config_dump`, `/clusters`, `/listeners`, `/stats`, `/logging`, `/ready`. |

Fonte: [Envoy / Life of a Request](https://www.envoyproxy.io/docs/envoy/latest/intro/life_of_a_request), [Envoy / Architecture Overview](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/arch_overview), [Envoy / Admin interface](https://www.envoyproxy.io/docs/envoy/latest/operations/admin).

## 1.4. Glossário de infraestrutura usada neste repositório

| Termo | Definição |
|---|---|
| **kind** | "Kubernetes IN Docker": ferramenta que roda um cluster Kubernetes usando containers Docker como "nós", usada neste repositório para criar o cluster de estudo/laboratório — ver [INSTALL_REQUIREMENTS.md#kind](INSTALL_REQUIREMENTS.md#kind). |
| **MetalLB** | Implementação de LoadBalancer para clusters bare-metal/kind, usada para atribuir IP externo ao `Service` do tipo `LoadBalancer` criado pelo Istio Ingress Gateway — ver [INSTALL_REQUIREMENTS.md#metallb](INSTALL_REQUIREMENTS.md#metallb). |
| **Helm** | Gerenciador de pacotes do Kubernetes, usado neste repositório para instalar todos os componentes do Istio (`istio-base`, `istiod`, `istio-cni`, `ztunnel`, `istio-ingress`) — ver [day1.md](day1.md). |
| **Bookinfo** | Aplicação de exemplo multi-serviço (productpage, details, reviews v1/v2/v3, ratings) mantida pelo projeto Istio, usada nos laboratórios deste repositório e deste guia. |

---

# 2. Arquitetura do Istio

## 2.1. Modo Sidecar

No modo sidecar (o modo "clássico" do Istio), cada Pod da aplicação recebe um container extra, o **Envoy** (`istio-proxy`), injetado automaticamente. Todo o tráfego de/para o Pod passa por esse Envoy (via regras `iptables` configuradas pelo `istio-cni` ou por um `initContainer`), que é configurado dinamicamente pelo **istiod** via **xDS**.

```mermaid
flowchart TB
    Client(["Cliente externo"])

    subgraph CP["Control plane - namespace istio-system"]
        ISTIOD["istiod<br/>(Pilot + Galley + Citadel)"]
    end

    subgraph IGW["istio-ingressgateway"]
        GwEnvoy["Envoy Gateway"]
    end

    subgraph NS["Namespace da aplicação (ex.: myapp)"]
        subgraph PodA["Pod productpage"]
            direction TB
            AppA["Container productpage"]
            EnvoyA["Envoy sidecar (istio-proxy)"]
        end
        subgraph PodB["Pod reviews"]
            direction TB
            AppB["Container reviews"]
            EnvoyB["Envoy sidecar (istio-proxy)"]
        end
    end

    Client -->|"HTTP/HTTPS :80/:443"| GwEnvoy
    GwEnvoy -->|"mTLS - inbound :15006"| EnvoyA
    AppA <-->|"localhost"| EnvoyA
    EnvoyA -->|"outbound :15001 -> mTLS -> inbound :15006"| EnvoyB
    EnvoyB <-->|"localhost"| AppB

    ISTIOD -->|"xDS gRPC :15012 (CDS/EDS/LDS/RDS/SDS)"| EnvoyA
    ISTIOD -->|"xDS gRPC :15012"| EnvoyB
    ISTIOD -->|"xDS gRPC :15012"| GwEnvoy
```

Pontos-chave (fonte: [Istio / Architecture](https://istio.io/v1.30/docs/ops/deployment/architecture/), [Istio / Dataplane modes](https://istio.io/v1.30/docs/overview/dataplane-modes/)):

- O **data plane** é o conjunto de proxies Envoy implantados como sidecars, responsáveis por descoberta de serviço, load balancing, terminação TLS, HTTP/2 e proxy gRPC, circuit breaking, health checks, injeção de falhas e coleta de métricas.
- O **control plane** (**istiod**) fornece descoberta de serviço, gerenciamento de configuração e gerenciamento de certificados; converte as regras de roteamento (CRDs) em configuração específica do Envoy e atua como Autoridade Certificadora (CA) para o mTLS.
- O `istiod` **encapsula** (desde a versão 1.5) as antigas funções separadas **Pilot**, **Galley** e **Citadel** em um único binário — ver [Istio / Glossary — Istiod](https://istio.io/v1.30/docs/reference/glossary/).

## 2.2. Modo Ambient

No modo **ambient** (usado por padrão nos comandos de instalação deste repositório — ver [day1.md](day1.md)) não há injeção de sidecar. A malha é dividida em duas camadas:

1. **Secure L4 overlay**, provida pelo **ztunnel** (um DaemonSet, um por nó) — cuida de mTLS, autenticação, autorização L4 e telemetria básica, **sem** terminar/interpretar HTTP.
2. **Camada L7 opcional**, provida por **waypoint proxies** (Envoy completo, implantado como Deployment, por namespace/service account/serviço) — necessária para roteamento HTTP, retries, fault injection e políticas L7 (path, método, header).

O **istio-cni** configura, no namespace de rede de cada Pod, o redirecionamento de tráfego para o `ztunnel` do nó (hoje via `iptables`, com suporte a `eBPF` dependendo da plataforma).

```mermaid
flowchart TB
    subgraph CP["Control plane - istio-system"]
        ISTIOD["istiod"]
    end

    subgraph Node1["Node Kubernetes 1"]
        direction TB
        CNI1["istio-cni node agent"]
        ZT1["ztunnel (DaemonSet)"]
        subgraph PodA["Pod productpage (sem sidecar)"]
            AppA["Container productpage"]
        end
    end

    subgraph Node2["Node Kubernetes 2"]
        direction TB
        CNI2["istio-cni node agent"]
        ZT2["ztunnel (DaemonSet)"]
        subgraph WP["waypoint (Deployment, opcional - L7)"]
            WPEnvoy["Envoy"]
        end
        subgraph PodB["Pod reviews (sem sidecar)"]
            AppB["Container reviews"]
        end
    end

    AppA -->|"redirecionado via iptables/eBPF"| ZT1
    ZT1 -->|"HBONE (HTTP/2 CONNECT + mTLS) :15008"| ZT2
    ZT2 -->|"se o serviço/pod tiver waypoint associado"| WPEnvoy
    WPEnvoy --> AppB
    ZT2 -.->|"sem policy L7: encaminha direto"| AppB

    ISTIOD -->|"xDS gRPC :15012"| ZT1
    ISTIOD -->|"xDS gRPC :15012"| ZT2
    ISTIOD -->|"xDS gRPC :15012"| WPEnvoy
    CNI1 -.->|"programa redirecionamento de tráfego"| ZT1
    CNI2 -.->|"programa redirecionamento de tráfego"| ZT2
```

Pontos-chave (fonte: [Istio / Ambient Overview](https://istio.io/latest/docs/ambient/overview/), [Istio / Ambient Architecture — Traffic Redirection](https://istio.io/v1.30/docs/ambient/architecture/traffic-redirection/), [Istio / Glossary](https://istio.io/v1.30/docs/reference/glossary/)):

- O **ztunnel** é escrito em Rust e é **intencionalmente restrito** a funções L3/L4: mTLS, autenticação, autorização L4 e telemetria — ele **não** termina tráfego HTTP nem interpreta headers HTTP da aplicação.
- **HBONE** ("HTTP-Based Overlay Network Environment") é o protocolo de tunelamento seguro usado entre ztunnels e entre ztunnel e waypoint; usa a porta **15008**.
- Redirecionamento de tráfego: entrada em texto plano é redirecionada para a porta **15006**; entrada via HBONE usa a porta **15008**; saída é redirecionada para a porta **15001**, usando `iptables` (tabelas `mangle`/`nat`, `TPROXY` e marcação de conexão).
- Políticas **L4** (identidade de origem, porta) podem ser aplicadas pelo próprio `ztunnel`; políticas **L7** (método HTTP, path, headers) exigem um **waypoint**.

## 2.3. Sidecar vs Ambient (comparação)

| Aspecto | Sidecar | Ambient |
|---|---|---|
| Proxy por Pod? | Sim (`istio-proxy` injetado em cada Pod) | Não (nenhum container extra no Pod da app) |
| Reinício de Pods ao entrar na mesh? | Sim (para injetar o sidecar) | Não (basta o label de namespace/label do Pod) |
| Componente L4 (mTLS, AuthZ L4) | Envoy sidecar | `ztunnel` (DaemonSet por nó) |
| Componente L7 (rotas HTTP, fault injection, retries) | Envoy sidecar | `waypoint` (opcional, por namespace/service account) |
| Linguagem do proxy de nó | N/A (não existe) | Rust (ztunnel) |
| Overhead de recursos por Pod | Um Envoy completo por Pod | Nenhum por padrão; custo concentrado no `ztunnel`/`waypoint` compartilhados |
| Uso de `iptables`/CNI | Sim | Sim |

Fonte: [Istio / Dataplane modes](https://istio.io/v1.30/docs/overview/dataplane-modes/), [Istio / Ambient Overview](https://istio.io/latest/docs/ambient/overview/).

## 2.4. Portas usadas pelo Istio e pelo Envoy

| Porta | Protocolo | Uso |
|---|---|---|
| 15000 | TCP | Porta de administração do Envoy (comandos/diagnóstico) |
| 15001 | TCP | Envoy outbound (interceptação de tráfego de saída) |
| 15002 | TCP | Porta de escuta para detecção de falhas |
| 15004 | HTTP | Porta de debug |
| 15006 | TCP | Envoy inbound (interceptação de tráfego de entrada) |
| 15008 | HTTP/2 | Porta do túnel mTLS **HBONE** |
| 15010 | gRPC (plaintext) | xDS do istiod em texto plano (não deve ser usada em produção) |
| 15012 | gRPC (mTLS) | xDS do istiod com TLS/mTLS — porta usada pelos proxies para obter configuração |
| 15014 | HTTP | Porta de controle/monitoramento do istiod (métricas para Prometheus e endpoints de debug) |
| 15020 | HTTP | Telemetria Prometheus consolidada (Istio agent + Envoy + aplicação) |
| 15021 | HTTP | Health checks |
| 15053 | DNS | Porta de DNS, quando o DNS capture está habilitado |
| 15090 | HTTP | Telemetria Prometheus do Envoy |

Fonte: [Istio / Application Requirements — Ports Used by Istio](https://istio.io/v1.30/docs/ops/deployment/application-requirements/), [Istio / Debug Endpoints](https://istio.io/latest/docs/ops/integrations/integration-guide/debug-endpoints/), [Istio / Ambient Architecture — Traffic Redirection](https://istio.io/v1.30/docs/ambient/architecture/traffic-redirection/).

## 2.5. Componentes de software do Istio

Esta seção detalha o **software** por trás de cada caixa dos diagramas das seções [2.1](#21-modo-sidecar) e [2.2](#22-modo-ambient).

| Componente | Onde roda | Função |
|---|---|---|
| **istiod** | 1+ Pods no namespace `istio-system` | Binário único de control plane. Fornece descoberta de serviço, gerenciamento/validação de configuração e gerenciamento de certificados (CA); converte os CRDs (VirtualService, DestinationRule, etc.) em configuração xDS e a distribui a todos os proxies (sidecars, gateways, ztunnel, waypoints). |
| **istio-agent** (binário `pilot-agent`) | Container `istio-proxy`, junto com **cada** Envoy sidecar/gateway | Processo auxiliar que roda ao lado do Envoy: gera a configuração de **bootstrap** do Envoy, atua como cliente/servidor **SDS** local (recebe certificados do istiod e os repassa ao Envoy) e monitora a expiração do certificado para solicitar sua renovação. |
| **Envoy sidecar / gateway** | Container `istio-proxy` (sidecar) ou Pod dedicado (gateway) | O proxy de dados propriamente dito — ver seção [3](#3-envoy-proxy). |
| **istio-cni** | DaemonSet `istio-cni-node`, `kube-system`/`istio-system` | Agente CNI que configura o redirecionamento de tráfego (`iptables`/`eBPF`) dos Pods para o Envoy sidecar (modo sidecar) ou para o `ztunnel` (modo ambient), evitando o uso de um `initContainer` privilegiado por Pod. |
| **ztunnel** | DaemonSet, um por nó (apenas ambient) | Ver seção [2.2](#22-modo-ambient). |
| **waypoint** | Deployment, opcional, por namespace/service account (apenas ambient) | Ver seção [2.2](#22-modo-ambient). |

O fluxo de emissão/rotação de certificados entre `istio-agent` e `istiod` — a base do mTLS automático do Istio — é descrito literalmente assim pela documentação oficial:

> "When started, the Istio agent creates the private key and CSR, and then sends the CSR with its credentials to `istiod` for signing." (...) "When a workload is started, Envoy requests the certificate and key from the Istio agent in the same container via the Envoy secret discovery service (SDS) API." (...) "The Istio agent monitors the expiration of the workload certificate. The above process repeats periodically for certificate and key rotation."

```mermaid
flowchart TB
    subgraph CP["istiod"]
        direction TB
        DISC["Discovery / xDS server<br/>(função Pilot)"]
        CFG["Ingestão e validação de config<br/>(função Galley)"]
        CA["Autoridade Certificadora<br/>(função Citadel)"]
        WH["Webhook de injeção e<br/>validação (Kubernetes)"]
    end

    K8S["Kubernetes API Server<br/>(Service, Endpoints, CRDs Istio)"] --> CFG
    CFG --> DISC

    subgraph PodApp["Pod da aplicacao"]
        direction TB
        APP["Container da aplicacao"]
        AGENT["istio-agent (pilot-agent)"]
        ENVOY["Envoy (istio-proxy)"]
    end

    DISC -->|"xDS gRPC :15012 (CDS/EDS/LDS/RDS)"| ENVOY
    AGENT -->|"1. Gera CSR e chave privada"| CA
    CA -->|"2. Assina e devolve o certificado"| AGENT
    AGENT -->|"3. Entrega certificado via SDS (Unix socket local)"| ENVOY
    APP <-->|"localhost"| ENVOY
    WH -.->|"injeta o container istio-proxy no Pod"| PodApp
```

Fonte: [Istio / Glossary — Istiod](https://istio.io/v1.30/docs/reference/glossary/), [Istio / Security Concepts](https://istio.io/v1.30/docs/concepts/security/), [Istio / Ambient Architecture — Traffic Redirection](https://istio.io/v1.30/docs/ambient/architecture/traffic-redirection/).

## 2.6. Topologias de implantação da mesh (multicluster, multi-network, multi-mesh)

Tudo visto até aqui (seções [2.1](#21-modo-sidecar)/[2.2](#22-modo-ambient)) assume implicitamente **um único cluster**. A documentação oficial do Istio ([Deployment Models](https://istio.io/v1.30/docs/ops/deployment/deployment-models/)) organiza as decisões de topologia de instalação em **dimensões independentes**, que podem ser combinadas: modelo de **cluster** (único ou múltiplos), modelo de **rede** (única ou múltiplas), modelo de **control plane** (único/compartilhado, múltiplos, ou externo), modelo de **mesh** (única ou múltiplas — federação) e modelo de **tenancy** (por namespace, por cluster ou por mesh).

> Fonte complementar (com exemplos práticos passo a passo): **POSTA, Christian E.; MALOKU, Rinor. *Istio in Action*. Manning Publications, 2022, ISBN 9781617295829 — Capítulo 12, "Scaling Istio in your organization"** (o mesmo Christian Posta do blog [blog.christianposta.com](https://blog.christianposta.com/tags/istio/) citado nas referências deste guia).

### 2.6.1. Modelos de cluster e de rede

Um **cluster único** (o que este repositório usa, com o **kind**) é a implantação mais simples, mas sem isolamento de falha entre times/regiões. Um **multicluster** dentro da mesma mesh dá isolamento de falha, roteamento por localidade e permite espalhar workloads por regiões/nuvens diferentes — ao custo de exigir 3 pré-requisitos, segundo o livro *Istio in Action*: **cross-cluster workload discovery** (os control planes precisam enxergar o Kubernetes API server dos clusters pares), **cross-cluster workload connectivity** (os workloads precisam conseguir se conectar de fato) e **common trust** (identidades mTLS mutuamente confiáveis entre clusters).

Dentro de um multicluster, o modelo de **rede** define como essa conectividade acontece:

- **Rede única (single/flat network)**: todos os clusters compartilham uma rede totalmente roteável (ex.: mesma VPC, VPCs pareadas, ou interligadas por um hub como o AWS Transit Gateway — seção [2.6.5](#265-exemplo-prático-rede-única-multi-conta-na-aws-com-aws-transit-gateway)). Os Pods se enxergam diretamente por IP, sem exigir um gateway do Istio para o tráfego cross-cluster; a limitação é que os IPs/CIDRs de serviço e Pod **não podem se sobrepor** entre os clusters.
- **Multi-network**: os clusters estão em redes isoladas (CIDRs podem até se sobrepor). Workloads de redes diferentes **só conseguem se comunicar através de um ou mais gateways do Istio** — chamados de **east-west gateway** — que fazem *SNI passthrough* (o gateway não decripta o mTLS; apenas repassa a conexão TLS até o proxy de destino, que é quem de fato termina o mTLS). O Istio usa *partitioned service discovery*: cada cliente recebe uma visão de endpoints adequada à sua rede (endpoint local direto, ou o IP do east-west gateway remoto).

```mermaid
flowchart TB
    subgraph SN["Modelo de rede UNICA - flat network"]
        direction LR
        A1["Workload - cluster 1"] <-->|"IP roteavel diretamente, sem gateway Istio"| A2["Workload - cluster 2"]
    end

    subgraph MN["Modelo MULTI-NETWORK"]
        direction LR
        B1["Workload - west-cluster<br/>rede: west-network"] --> EWG1["East-West Gateway<br/>west-cluster"]
        EWG1 -->|"mTLS - SNI passthrough"| EWG2["East-West Gateway<br/>east-cluster"]
        EWG2 --> B2["Workload - east-cluster<br/>rede: east-network"]
    end
```

### 2.6.2. Modelos de control plane

Independentemente da rede, o Istio distingue **cluster primário** (onde o control plane roda) de **cluster remoto** (sem control plane próprio, gerenciado por um primário). A partir daí, há três modelos:

| Modelo | Também chamado de | Descrição |
|---|---|---|
| **Primary-remote** | *single/shared control plane* | Um único control plane gerencia todos os clusters. Usa menos recursos, mas uma indisponibilidade do cluster primário afeta a mesh inteira (baixa disponibilidade). |
| **Primary-primary** | *replicated control plane* | Cada cluster roda seu próprio control plane; eles se descobrem mutuamente. Maior disponibilidade (uma indisponibilidade fica restrita ao cluster afetado), ao custo de mais recursos. |
| **External control plane** | — | Todos os clusters são "remotos"; o control plane roda fora deles (ex.: em um cluster dedicado, ou fornecido como serviço gerenciado por um cloud provider). |

```mermaid
flowchart TB
    subgraph PR["Primary-Remote (control plane unico/compartilhado)"]
        direction LR
        P1["Cluster Primario<br/>Istio control plane"] -->|Discovery| R1["Cluster Remoto<br/>sem control plane proprio"]
    end

    subgraph PP["Primary-Primary (control plane replicado)"]
        direction LR
        CP1["Cluster A<br/>Istio control plane"] <-->|Discovery mutua| CP2["Cluster B<br/>Istio control plane"]
    end

    subgraph EX["External control plane"]
        direction LR
        EXT["Control plane externo<br/>ex: gerenciado pelo cloud provider"] -->|Discovery| REM1["Cluster remoto 1"]
        EXT -->|Discovery| REM2["Cluster remoto 2"]
    end
```

Para que o discovery cross-cluster funcione, cada control plane precisa se autenticar no API server dos clusters pares — o Istio faz isso com uma conta de serviço (`istio-reader-service-account`) e um **remote secret**, criado com `istioctl x create-remote-secret` e aplicado no(s) cluster(s) que devem enxergar aquele cluster remoto. A documentação oficial também destaca um requisito de infraestrutura: **"the control plane in a primary cluster must be accessible via a stable IP"** — em clusters de redes diferentes, isso é obtido expondo o control plane por um gateway do Istio.

### 2.6.3. Modelos de mesh: mesh única vs. múltiplas meshes (federação)

- **Single mesh**: nomes de serviço são únicos globalmente dentro da mesh, mesmo que ela se estenda por múltiplos clusters/redes; namespaces do Kubernetes com o mesmo nome em clusters diferentes são tratados, por padrão, como o **mesmo** namespace lógico (os endpoints são somados automaticamente).
- **Multiple meshes (mesh federation / multi-mesh)**: duas ou mais meshes **independentes** (cada uma com seu próprio control plane e, tipicamente, sua própria CA) expõem um subconjunto de serviços/identidades uma para a outra através de gateways — sem compartilhar o mesmo control plane nem exigir que os clusters confiem plenamente uns nos outros. É a opção recomendada quando as meshes são operadas por **times diferentes** ou têm requisitos fortes de isolamento de segurança; a documentação registra explicitamente que **"Istio does not provide any tooling to exchange trust bundles across meshes"** — a troca de confiança entre meshes é manual (ou via SPIFFE Trust Domain Federation).

```mermaid
flowchart TB
    subgraph SM["Single mesh (multicluster)"]
        direction LR
        M1["Cluster 1<br/>mesh: usmesh"] <-.->|"mesmo control plane / meshID"| M2["Cluster 2<br/>mesh: usmesh"]
    end

    subgraph MM["Multiple meshes (mesh federation)"]
        direction LR
        subgraph MA["Mesh do time A - CA propria"]
            MA1["Cluster(s) do time A"]
        end
        subgraph MB["Mesh do time B - CA propria"]
            MB1["Cluster(s) do time B"]
        end
        MA1 -->|"Gateway + troca manual de trust bundle"| MB1
    end
```

### 2.6.4. Exemplo prático completo: multicluster + multi-network + multi-control-plane

O capítulo 12 do livro *Istio in Action* constrói passo a passo um exemplo real combinando os três modelos acima: dois clusters (`west-cluster` e `east-cluster`), cada um em uma **rede diferente** (`west-network`/`east-network`, em regiões diferentes de um cloud provider), cada um com **seu próprio control plane** (modelo primary-primary, escolhido pela alta disponibilidade), conectados por **east-west gateways**, com **confiança comum via CA plugada** (intermediate CAs próprias em cada cluster, assinadas pela mesma root CA):

```mermaid
flowchart TB
    ROOT["Root CA comum"] -->|assina| ICA1["Intermediate CA - west-cluster"]
    ROOT -->|assina| ICA2["Intermediate CA - east-cluster"]

    subgraph WC["Cluster west-cluster - rede west-network"]
        IP1["Istio Pilot - west"]
        EWG_W["East-West Gateway"]
        WEBAPP["webapp"]
        ICA1 -.->|"secret cacerts"| IP1
    end

    subgraph EC["Cluster east-cluster - rede east-network"]
        IP2["Istio Pilot - east"]
        EWG_E["East-West Gateway"]
        CATALOG["catalog"]
        ICA2 -.->|"secret cacerts"| IP2
    end

    IP1 <-->|"Discovery mutua via remote secret"| IP2
    WEBAPP --> EWG_W -->|"mTLS - SNI passthrough"| EWG_E --> CATALOG
```

Resumo do fluxo (comandos completos e arquivos estão no capítulo 12 do livro): (1) provisionar 2 clusters em redes/regiões separadas; (2) gerar uma root CA e duas intermediate CAs, aplicando cada intermediate como o Secret `cacerts` no namespace `istio-system` do respectivo cluster (seção [8.3](#83-certificados-mtls-plugando-uma-ca-própria-e-rotação) deste guia); (3) rotular o namespace `istio-system` de cada cluster com `topology.istio.io/network=<nome-da-rede>`; (4) instalar o control plane em cada cluster via `IstioOperator`, com `spec.values.global.meshID` **igual** nos dois (para formarem a mesma mesh) e `spec.values.global.multiCluster.clusterName`/`spec.values.global.network` distintos por cluster; (5) gerar e aplicar o *remote secret* de cada cluster no outro com `istioctl x create-remote-secret --name=<cluster>`; (6) confirmar nos logs do `istiod` a linha `Adding cluster_id=<cluster> from secret=...`.

### 2.6.5. Exemplo prático: rede única multi-conta na AWS com AWS Transit Gateway

O modelo de **rede única** (seção [2.6.1](#261-modelos-de-cluster-e-de-rede)) não exige que todos os clusters estejam na mesma conta/VPC — só exige que a rede seja **totalmente roteável** e sem CIDRs sobrepostos. Uma forma comum de conseguir isso entre **múltiplas contas AWS** (um padrão organizacional frequente: uma conta por time/ambiente) é o **AWS Transit Gateway**, descrito oficialmente como um serviço que **"connects your Amazon Virtual Private Clouds (VPCs) and on-premises networks through a central hub"**, funcionando como **"a highly scalable cloud router"** em topologia hub-and-spoke — cada VPC (spoke) se conecta uma única vez ao hub, evitando o N×N de conexões de peering diretas entre VPCs.

```mermaid
flowchart TB
    TGW["AWS Transit Gateway<br/>(hub)"]

    subgraph ACC1["Conta AWS 1 - VPC A"]
        EKS1["Cluster Kubernetes 1<br/>Istio control plane (primary-primary)"]
    end
    subgraph ACC2["Conta AWS 2 - VPC B"]
        EKS2["Cluster Kubernetes 2<br/>Istio control plane (primary-primary)"]
    end
    subgraph ACC3["Conta AWS 3 - VPC C"]
        EKS3["Cluster Kubernetes 3<br/>Istio control plane (primary-primary)"]
    end

    EKS1 <--> TGW
    EKS2 <--> TGW
    EKS3 <--> TGW
```

Nesse desenho, como as três VPCs formam uma **única rede roteável** (via Transit Gateway compartilhado entre as contas, usando AWS Resource Access Manager para o compartilhamento cross-account), a mesh usa o modelo de **rede única** da seção [2.6.1](#261-modelos-de-cluster-e-de-rede) — **não** é necessário um east-west gateway para a conectividade cross-cluster (diferente do exemplo da seção [2.6.4](#264-exemplo-prático-completo-multicluster--multi-network--multi-control-plane), que precisa dele por estar em redes isoladas). O modelo de control plane (primary-remote ou primary-primary) e o de confiança (CA própria ou externa) continuam sendo escolhas **independentes** da topologia de rede, exatamente como descrito na seção [2.6.2](#262-modelos-de-control-plane).

> Este é um padrão de infraestrutura observado na prática (a documentação oficial do Istio não cita nomes de produtos de nuvem específicos); os CIDRs de Pod/Service de cada cluster ainda precisam ser planejados para não se sobrepor entre si, um requisito explícito do modelo de rede única do Istio.

Fonte: [Istio / Deployment Models](https://istio.io/v1.30/docs/ops/deployment/deployment-models/), *Istio in Action* (Posta & Maloku, Manning 2022) — capítulo 12, [AWS / Transit Gateway](https://aws.amazon.com/transit-gateway/).

---

# 3. Envoy Proxy

## 3.1. Visão geral e histórico

O **Envoy** é um proxy de rede L4/L7 de alto desempenho, escrito em C++, originalmente criado por **Matt Klein** na **Lyft** para substituir uma infraestrutura de comunicação de serviços difícil de operar por um "sidecar" de comunicação rápido e observável. Envoy foi o terceiro projeto a se graduar na **CNCF**, depois de Kubernetes e Prometheus. O Istio usa o Envoy como seu data plane (tanto no modo sidecar quanto, de forma adaptada, no `ztunnel`/`waypoint` do modo ambient).

Fonte: [Envoy / Introduction](https://www.envoyproxy.io/docs/envoy/latest/intro/intro), [The New Stack — Matt Klein on the Success of Envoy](https://thenewstack.io/matt-klein-on-the-success-of-envoy-and-the-future-of-the-service-mesh/), [Kubernetes Podcast from Google — Episode 221: Creating Envoy, with Matt Klein](https://kubernetespodcast.com/episode/221-creating-envoy/).

## 3.2. Arquitetura interna (listener subsystem x cluster subsystem)

A documentação oficial do Envoy descreve a arquitetura interna como dividida em dois grandes subsistemas, conectados pelo **router filter**:

- **Listener subsystem**: processa as requisições vindas do downstream, cuida do ciclo de vida da conexão e do caminho de resposta.
- **Cluster subsystem**: seleciona/configura a conexão upstream, gerencia health checking, load balancing e connection pooling.

```mermaid
flowchart LR
    subgraph Downstream["Downstream"]
        C["Conexão TCP/UDP do cliente"]
    end

    subgraph EnvoyProc["Processo Envoy"]
        direction TB
        LST["Listener<br/>(bind IP:porta)"]
        LSTF["Listener Filters<br/>(ex.: TLS Inspector)"]
        FCM["Filter Chain Match<br/>(SNI / ALPN / porta)"]
        TSOC["Transport Socket<br/>(TLS termination)"]
        NFC["Network Filter Chain<br/>(último filtro: HTTP Connection Manager)"]
        HFC["HTTP Filter Chain<br/>(ex.: CORS, JWT, RBAC)"]
        RTF["Router Filter<br/>(usa a Route Table / RDS)"]
        CLM["Cluster Manager<br/>(CDS)"]
        LB["Load Balancer"]
        CP["Connection Pool"]
    end

    subgraph Upstream["Upstream"]
        EP["Endpoints (EDS)"]
    end

    C --> LST --> LSTF --> FCM --> TSOC --> NFC --> HFC --> RTF --> CLM --> LB --> CP --> EP
```

Fonte: [Envoy / Life of a Request](https://www.envoyproxy.io/docs/envoy/latest/intro/life_of_a_request), [Envoy / Architecture Overview](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/arch_overview).

## 3.3. xDS: como o Istio configura o Envoy

O Envoy não lê um arquivo estático de configuração de produção; ele é configurado **dinamicamente** por um *management server* (no Istio, o **istiod**) através da família de APIs **xDS** ("`x` Discovery Service"):

| API | Nome completo | Função |
|---|---|---|
| **LDS** | Listener Discovery Service | Descobre quais **listeners** (IP:porta) o Envoy deve abrir |
| **RDS** | Route Discovery Service | Descobre as **tabelas de rotas** usadas pelo HTTP Connection Manager |
| **CDS** | Cluster Discovery Service | Descobre quais **clusters** (serviços upstream) existem |
| **EDS** | Endpoint Discovery Service | Descobre os **endpoints** (IP:porta de cada réplica) de cada cluster |
| **SDS** | Secret Discovery Service | Distribui **certificados/chaves** (usado pelo Istio para os certificados mTLS) |
| **ADS** | Aggregated Discovery Service | Agrega **todas** as APIs acima em um único stream gRPC bidirecional, permitindo sequenciamento correto das atualizações |

Segundo a documentação oficial do Envoy, para evitar perda de tráfego durante atualizações, a ordem recomendada de aplicação é: **CDS → EDS → LDS → RDS** (clusters e endpoints antes dos listeners, e rotas relacionadas a novos listeners só depois de CDS/EDS/LDS).

```mermaid
sequenceDiagram
    participant K as Kubernetes API Server
    participant Istiod as istiod (Pilot)
    participant Envoy as Envoy (sidecar / ztunnel / waypoint)

    K-->>Istiod: Watch de Service, Endpoints, VirtualService, DestinationRule, PeerAuthentication...
    Envoy->>Istiod: Conecta stream ADS (gRPC mTLS porta 15012)
    Istiod-->>Envoy: CDS (Clusters)
    Istiod-->>Envoy: EDS (Endpoints dos clusters)
    Istiod-->>Envoy: LDS (Listeners)
    Istiod-->>Envoy: RDS (Rotas dos listeners)
    Istiod-->>Envoy: SDS (Certificados mTLS)
    Envoy-->>Istiod: ACK (versionInfo + nonce) ou NACK em caso de erro
```

Fonte: [Envoy / xDS REST and gRPC protocol](https://www.envoyproxy.io/docs/envoy/latest/api-docs/xds_protocol), [Istio / Architecture](https://istio.io/v1.30/docs/ops/deployment/architecture/), [Istio / Debugging Envoy and Istiod](https://istio.io/v1.30/docs/ops/diagnostic-tools/proxy-cmd/).

O comando `istioctl proxy-status` mostra o estado de sincronização de cada proxy da malha em relação ao istiod:

```mermaid
stateDiagram-v2
    [*] --> NOT_SENT
    NOT_SENT --> SYNCED: istiod envia configuração e o Envoy confirma - ACK
    SYNCED --> STALE: Envoy não confirma uma nova config a tempo - ex. falha de rede ou CPU
    STALE --> SYNCED: Reconexão e novo ACK bem-sucedido
    SYNCED --> [*]
```

Fonte: [Istio / Debugging Envoy and Istiod — proxy-status](https://istio.io/v1.30/docs/ops/diagnostic-tools/proxy-cmd/).

## 3.4. Life of a Request (o ciclo de vida de uma requisição no Envoy)

Resumo fiel (traduzido) do fluxo descrito na documentação oficial *"Life of a Request"* do Envoy, para uma requisição HTTP/2 chegando em um listener com TLS:

1. O **Listener Manager** instancia listeners vinculados a um IP/porta; cada *worker thread* mantém sua própria instância de listener, aceitando novas conexões TCP (callback `onAccept()`).
2. Os **Listener Filters** rodam (ex.: `TLS Inspector`, que extrai o SNI do handshake). Essas informações são usadas no **filter chain match**, que decide qual `transport socket` e quais `network filters` serão aplicados à conexão.
3. O **Transport Socket** TLS realiza o handshake e entrega um fluxo de bytes já decriptado ao gerenciador de filtros (`FilterManagerImpl`), responsável pela cadeia de **network filters**.
4. Os **Network Filters** são instanciados por conexão. O **HTTP Connection Manager (HCM)** é tipicamente o último network filter, responsável pelo codec HTTP e pela cadeia de HTTP filters.
5. O **codec HTTP/2** decodifica os frames e demultiplexa a conexão em múltiplos **streams** independentes — cada stream corresponde a uma requisição/resposta.
6. Para cada stream, o HCM instancia uma cadeia de **HTTP filters** (`StreamDecoderFilter`, `StreamEncoderFilter` ou `StreamFilter`). Filtros customizados (ex.: autenticação, CORS) processam a requisição primeiro; o **router filter** roda por último, selecionando a rota e o cluster de destino.
7. O **Load Balancer** do cluster escolhe um **endpoint** (algoritmos: round robin ponderado, Maglev, least-request, etc.). O **connection pool** reutiliza uma conexão existente ou abre uma nova, respeitando os **circuit breakers** configurados.
8. O **codec HTTP/2 upstream** multiplexa o stream de requisição junto com outros streams sobre a mesma conexão TCP.
9. O **Transport Socket** TLS upstream (se aplicável) criptografa os bytes e escreve no socket TCP da conexão upstream.
10. A **resposta** propaga-se de volta pelo caminho inverso, passando pelos filtros HTTP em **ordem inversa**; a *outlier detection* é atualizada com o resultado.
11. Ao final: estatísticas são atualizadas, o **access log** é escrito (arquivo, statsd, etc.) e os **spans de tracing** são finalizados.

```mermaid
sequenceDiagram
    participant D as Downstream (cliente)
    participant L as Listener
    participant LF as Listener Filters (ex.: TLS Inspector)
    participant TS as Transport Socket (TLS)
    participant NF as Network Filter Chain
    participant HCM as HTTP Connection Manager
    participant HF as HTTP Filter Chain
    participant R as Router Filter
    participant CM as Cluster Manager / Load Balancer
    participant U as Upstream (endpoint)

    D->>L: Nova conexão TCP
    L->>LF: onAccept() aceita a conexão
    LF->>TS: SNI/ALPN detectados -> escolhe a filter chain (match)
    TS->>NF: Handshake TLS concluído, decripta o fluxo de bytes
    NF->>HCM: Último network filter decodifica HTTP/2 (ou HTTP/1.1)
    HCM->>HF: Cria o stream e a cadeia de HTTP filters
    HF->>R: Filtros customizados processam a requisição (auth, cors...)
    R->>CM: Router filter seleciona a rota e o cluster
    CM->>CM: Load balancing (round robin, least-request, Maglev...)
    CM->>U: Reusa conexão do pool ou cria uma nova (respeita circuit breaker)
    U-->>CM: Resposta HTTP
    CM-->>R: Resposta propaga de volta
    R-->>HF: Filtros processam a resposta (ordem inversa)
    HF-->>HCM: Resposta finalizada
    HCM-->>NF: Codifica em HTTP/2 (ou HTTP/1.1)
    NF-->>TS: Criptografa (se TLS)
    TS-->>D: Envia os bytes de resposta
    Note over D,U: Ao final - stats, access log e trace span são finalizados
```

Fonte (fiel/tradução): [Envoy / Life of a Request](https://www.envoyproxy.io/docs/envoy/latest/intro/life_of_a_request).

## 3.5. Vídeos recomendados de Matt Klein (criador do Envoy)

Lista curada a partir da página oficial de palestras de Matt Klein ([mattklein123.dev/appearances](https://mattklein123.dev/appearances/)), priorizando os conteúdos mais relevantes para entender a arquitetura do Envoy e a origem do conceito de service mesh:

| Palestra | Evento / Ano | Link |
|---|---|---|
| Envoy Internals Deep Dive | KubeCon EU, 2018 (nível avançado: threading model, stats, hot restart, xDS) | [youtube.com/watch?v=gQF23Vw0keg](https://www.youtube.com/watch?v=gQF23Vw0keg) |
| Envoy - Embracing a Service Mesh | QCon NY, 2018 | [youtube.com/watch?v=55yi4MMVBi4](https://www.youtube.com/watch?v=55yi4MMVBi4) |
| The Mechanics of Deploying Envoy at Lyft | KubeCon NA, 2017 | [youtube.com/watch?v=IeJDjq-COjk](https://www.youtube.com/watch?v=IeJDjq-COjk) |
| Lyft's Envoy: Experiences Operating a Large Service Mesh | SREcon Americas, 2017 | [youtube.com/watch?v=fZBf6MLjNqM](https://www.youtube.com/watch?v=fZBf6MLjNqM) e [usenix.org (resumo)](https://www.usenix.org/conference/srecon17americas/program/presentation/klein) |
| Service Mesh Day Keynote | Service Mesh Day, 2019 | [youtube.com/watch?v=IC0eduR56DA](https://www.youtube.com/watch?v=IC0eduR56DA) |
| Creating Envoy, with Matt Klein (podcast) | Kubernetes Podcast from Google, ep. 221 | [kubernetespodcast.com/episode/221-creating-envoy](https://kubernetespodcast.com/episode/221-creating-envoy/) |

Artigo escrito (não vídeo) complementar, também de Matt Klein, muito citado como introdução ao tema: [Service mesh data plane vs. control plane](https://blog.envoyproxy.io/service-mesh-data-plane-vs-control-plane-2774e720f7fc) e [Introduction to modern network load balancing and proxying](https://blog.envoyproxy.io/introduction-to-modern-network-load-balancing-and-proxying-a57f6ff80236).

## 3.6. Componentes de software do Envoy: threading model e hot restart

O Envoy é um único processo ("single process, multiple threads") com três tipos de thread, segundo a documentação oficial:

- **Main thread**: não trata tráfego de alto volume diretamente. Processa as atualizações de configuração (xDS), o flush periódico de estatísticas (para StatsD/Prometheus), hospeda o **admin listener** (seção [4.3](#43-admin-api-do-envoy-porta-15000)) e trata sinais do sistema operacional (`SIGTERM`, `SIGHUP`).
- **Worker threads**: número configurável (flag `--concurrency`) que fazem o trabalho pesado — escuta, filtragem e encaminhamento do tráfego. Quando uma conexão é aceita por um listener, ela fica vinculada a **uma única** worker thread durante toda a sua vida útil, o que mantém o caminho de dados ("hot path") altamente paralelizado e evita locks complexos.
- **File flusher thread**: dedicada a gravar os access logs em disco, evitando bloquear o caminho principal de processamento.

Cada thread (main e worker) executa um laço de eventos não-bloqueante em torno de um **`Event::Dispatcher`** (wrapper sobre a biblioteca `libevent`), que trata file descriptors, timers e sinais via callbacks em vez de I/O bloqueante.

Para evitar contenção de locks, o Envoy usa **Thread Local Storage (TLS)**: o código aloca um "slot" na main thread, que funciona como índice em um vetor mantido em cada thread; ao atualizar configuração, a main thread posta um closure para cada worker thread, que cria/atualiza sua própria cópia local dos dados — resultando em leitura O(1) sem lock no caminho de dados.

```mermaid
flowchart TB
    subgraph Main["Main thread"]
        XDS["Processamento xDS<br/>(recebe config do istiod)"]
        STATS["Flush periodico de stats"]
        ADMIN["Admin listener :15000"]
        SIG["Tratamento de sinais do SO"]
    end

    subgraph W1["Worker thread 1"]
        D1["Event::Dispatcher"]
        C1["Conexoes atribuidas a este worker"]
    end
    subgraph W2["Worker thread 2"]
        D2["Event::Dispatcher"]
        C2["Conexoes atribuidas a este worker"]
    end
    subgraph WN["Worker thread N (--concurrency)"]
        DN["Event::Dispatcher"]
        CN["Conexoes atribuidas a este worker"]
    end

    FF["File flusher thread<br/>(grava access logs em disco)"]

    XDS -->|"posta closure via TLS: atualiza config"| D1
    XDS -->|"posta closure via TLS: atualiza config"| D2
    XDS -->|"posta closure via TLS: atualiza config"| DN
    C1 -.-> FF
    C2 -.-> FF
    CN -.-> FF
```

**Hot restart**: é o mecanismo (independente do Envoy usado como binário standalone) que permite trocar de binário/configuração **sem derrubar conexões existentes**: um novo processo Envoy é iniciado, herda os sockets de escuta do processo antigo via Unix domain sockets, e só então o processo antigo entra em modo de *draining* (encerrando gradualmente as conexões restantes, com agressividade configurável por `--drain-time-s`) até finalizar (`--parent-shutdown-time-s`). No Istio, a maior parte das mudanças de configuração **não** exige hot restart — elas chegam dinamicamente via xDS (seção [3.3](#33-xds-como-o-istio-configura-o-envoy)) sem reiniciar o processo Envoy; o hot restart importa principalmente quando o **binário** do Envoy muda de versão (ex.: upgrade da imagem do `istio-proxy`), o que no Istio tipicamente é feito recriando o Pod, não via hot restart em tempo de execução.

Fonte: [Envoy / Threading Model](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/intro/threading_model), [Envoy / Hot Restart](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/hot_restart).

## 3.7. A cadeia de filtros HTTP em detalhe: tipos, ordem e controle de fluxo

Esta seção aprofunda o que a seção [3.4](#34-life-of-a-request-o-ciclo-de-vida-de-uma-requisição-no-envoy) já introduziu sobre a **HTTP filter chain**: os tipos exatos de filtro, a ordem de execução na ida e na volta, e como um filtro pode pausar ou encerrar a cadeia antes do `router`.

### Os três tipos de filtro HTTP

A documentação oficial do Envoy define três interfaces:

| Interface | Callbacks | Atua em |
|---|---|---|
| `StreamDecoderFilter` | `decodeHeaders()`, `decodeData()`, `decodeTrailers()` | **Requisição** (caminho de entrada/downstream → upstream) |
| `StreamEncoderFilter` | `encodeHeaders()`, `encodeData()`, `encodeTrailers()` | **Resposta** (caminho de saída/upstream → downstream) |
| `StreamFilter` | Ambos os conjuntos acima | Requisição **e** resposta |

### Ordem de execução: ida direta, volta invertida

A regra oficial, citada literalmente: **"Decoder and encoder-decoder filters are executed on the request path. Encoder and encoder-decoder filters are executed on the response path, in reverse direction."** Ou seja: se a cadeia HTTP configurada for `A → B → C → router` (nessa ordem, do primeiro ao último filtro), a **requisição** passa por `A → B → C → router`, mas a **resposta** volta na ordem **inversa**: `router → C → B → A → cliente`. Isso é o que permite, por exemplo, um filtro de compressão (seção [7.6](#76-exemplo-prático-habilitando-compressão-gzip-via-envoyfilter)) processar a resposta depois que ela já foi gerada pelo `router`, mesmo tendo sido "inserido" antes dele na configuração.

```mermaid
flowchart LR
    subgraph REQ["Requisicao - ordem direta (decode)"]
        direction LR
        D["Downstream<br/>cliente"] --> MX["istio.metadata_exchange<br/>(StreamFilter)"]
        MX --> JWT["envoy.filters.http.jwt_authn<br/>(StreamDecoderFilter)"]
        JWT --> RBAC["envoy.filters.http.rbac<br/>(StreamDecoderFilter)"]
        RBAC --> STATS["istio.stats<br/>(StreamFilter)"]
        STATS --> RT["envoy.filters.http.router<br/>(ultimo da cadeia)"]
    end
```

```mermaid
flowchart RL
    subgraph RESP["Resposta - ordem INVERSA (encode)"]
        direction RL
        RT2["envoy.filters.http.router"] --> STATS2["istio.stats<br/>(StreamFilter)"]
        STATS2 --> RBAC2["envoy.filters.http.rbac<br/>NAO participa - e so decoder"]
        RBAC2 -.->|"pulado: nao tem encodeHeaders"| MX2["istio.metadata_exchange<br/>(StreamFilter)"]
        MX2 --> D2["Downstream<br/>cliente"]
    end
```

> A composição real da cadeia (quais filtros existem) **depende de quais recursos Istio estão configurados** para aquele workload: `envoy.filters.http.jwt_authn` só aparece se existir uma `RequestAuthentication` (seção [8.6](#86-autorizando-por-jwt-eou-mtls-requestauthentication--authorizationpolicy)); `envoy.filters.http.rbac` só aparece se existir uma `AuthorizationPolicy` (seção [8.5](#85-authorizationpolicy-autorizando-apenas-aplicações-específicas)); `envoy.filters.http.fault`/`cors` só aparecem se configurados no `VirtualService`. Os filtros `istio.metadata_exchange` (troca metadados de identidade entre proxies) e `istio.stats` (telemetria — seção [5.1](#51-métricas-padrão-do-istio)) são inseridos pelo próprio Istio via `EnvoyFilter` internos. **Sempre confirme a ordem real** com os comandos da seção [3.8](#38-comandos-de-troubleshooting-da-cadeia-de-filtros) em vez de assumir esta lista — a ordem já teve variações entre versões do Istio.

### Controle de fluxo: `StopIteration`/`Continue` e "local reply"

Cada callback (`decodeHeaders`, etc.) retorna um `FilterStatus` que decide se a cadeia continua imediatamente (`Continue`) ou pausa (`StopIteration`) — citado oficialmente: **"the `FilterStatus` allows filters to pause execution of the filter chain."** Isso é o que permite um filtro fazer uma chamada **assíncrona** (ex.: um serviço externo de autorização) sem bloquear a *worker thread* (seção [3.6](#36-componentes-de-software-do-envoy-threading-model-e-hot-restart)): o filtro retorna `StopIteration`, e quando a resposta assíncrona chega, ele mesmo chama `continueDecoding()`/`continueEncoding()` para retomar a cadeia de onde parou.

Um filtro também pode **encerrar a cadeia mais cedo**, sem nunca chegar ao `router` — chamado de **"local reply"**: por exemplo, um filtro de rate limiting retornando `429` diretamente. Comportamento oficial: se os headers de resposta do upstream **ainda não** foram enviados, o Envoy gera uma resposta internamente; se já tiverem sido enviados, o Envoy reseta o stream.

```mermaid
sequenceDiagram
    participant C as Cliente
    participant F1 as Filtro A (ex.: jwt_authn)
    participant EXT as Servico externo (ex.: ext_authz)
    participant F2 as Filtro B (ex.: rbac)
    participant R as Router

    C->>F1: decodeHeaders()
    F1->>EXT: chamada assincrona
    F1-->>F1: retorna StopIteration (pausa a cadeia)
    Note over F1,EXT: worker thread NAO fica bloqueada esperando
    EXT-->>F1: resposta assincrona recebida
    F1->>F1: continueDecoding() - retoma a cadeia
    F1->>F2: decodeHeaders()
    alt F2 nega a requisicao
        F2-->>C: local reply 403 (RBAC: access denied) - router NUNCA e chamado
    else F2 permite
        F2->>R: decodeHeaders()
        R->>R: seleciona rota e cluster, consulta o ClusterManager
    end
```

Fonte: [Envoy / Life of a Request](https://www.envoyproxy.io/docs/envoy/latest/intro/life_of_a_request), [Christian Posta — Understanding Istio Telemetry v2](https://blog.christianposta.com/understanding-istio-telemetry-v2/) (exemplo real do filtro `istio.stats` inserido com `INSERT_BEFORE` o `router`, com a observação literal do autor: *"the router should be the last filter in this chain"*), [Envoy / RBAC filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/rbac_filter).

## 3.8. Comandos de troubleshooting da cadeia de filtros

```bash
export PP_POD=$(kubectl -n $MY_NAMESPACE get pod -l app=productpage -o jsonpath='{.items[0].metadata.name}')

# Lista, NA ORDEM REAL configurada, os filtros de rede do listener inbound (porta 9080)
istioctl proxy-config listener $PP_POD.$MY_NAMESPACE --port 9080 -o json \
  | jq '.[0].filterChains[0].filters[].name'

# Lista, NA ORDEM REAL, os filtros HTTP dentro do HTTP Connection Manager
# (esta é a forma mais confiável de responder "qual filtro roda antes de qual" - seção 3.7)
istioctl proxy-config listener $PP_POD.$MY_NAMESPACE --port 9080 -o json \
  | jq -r '.[0].filterChains[0].filters[]
      | select(.name=="envoy.filters.network.http_connection_manager")
      | .typedConfig.httpFilters[].name'

# Confirme que o router (envoy.filters.http.router) é sempre o ÚLTIMO da lista acima
# Compare a cadeia INBOUND (porta da app) com a OUTBOUND (chamadas que o Pod faz) - costumam ser diferentes
istioctl proxy-config listener $PP_POD.$MY_NAMESPACE --port 15001 -o json \
  | jq -r '.[0].filterChains[0].filters[]
      | select(.name=="envoy.filters.network.http_connection_manager")
      | .typedConfig.httpFilters[].name'

# Mesma informação via Admin API / config_dump (seção 4.3), sem depender do istioctl
kubectl -n $MY_NAMESPACE port-forward $PP_POD 15000:15000 &
curl -s http://localhost:15000/config_dump | jq -r '.. | .httpFilters? // empty | .[].name'
kill %1
```

**Descobrir qual filtro gerou uma resposta 4xx/5xx sem passar pelo `router`** (uma "local reply" — seção [3.7](#37-a-cadeia-de-filtros-http-em-detalhe-tipos-ordem-e-controle-de-fluxo)): o Envoy registra a causa exata no campo `%RESPONSE_CODE_DETAILS%` do access log (já presente no formato padrão do Istio — seção [6.3](#63-access-log-habilitar-e-customizar-campos), logo após os `RESPONSE_FLAGS`). Valores típicos e o que cada um indica:

| `RESPONSE_CODE_DETAILS` | Indica |
|---|---|
| `rbac_access_denied_matched_policy[<nome>]` | O filtro `envoy.filters.http.rbac` negou a requisição — a `AuthorizationPolicy` (seção [8.5](#85-authorizationpolicy-autorizando-apenas-aplicações-específicas)) funcionou como configurado (`<nome>` vem vazio/`none` se nenhuma policy `ALLOW` casou) |
| `route_not_found` | O HTTP Connection Manager não achou nenhuma rota (RDS) para o `Host`/path da requisição — response flag `NR` (seção [5.4](#54-response-flags-do-envoy)) |
| `filter_chain_not_found` | Nenhuma filter chain do listener casou com a conexão (ex.: SNI inesperado) |
| `no_healthy_upstream` | O cluster existe, mas nenhum endpoint está saudável — response flag `UH` |
| `direct_response` | Uma resposta direta configurada (ex.: em uma rota de redirect) foi usada, sem chegar ao upstream |
| `via_upstream` | Fluxo normal — a resposta veio mesmo do serviço upstream, passando pelo `router` |

```bash
# Ver o RESPONSE_CODE_DETAILS de requisições recentes (já presente no access log padrão do Istio)
kubectl -n $MY_NAMESPACE logs deploy/reviews-v1 -c istio-proxy | grep -oE '"[a-z_]+(\[[^]]*\])?" "[^"]*"' | head
```

**Causas mais comuns de "a ordem dos filtros não é a que eu esperava":** um `EnvoyFilter` (seção [7](#7-envoyfilter)) com `patch.operation: INSERT_FIRST` competindo com outro `EnvoyFilter` também com `INSERT_FIRST` (o resultado depende da ordem de aplicação entre `EnvoyFilter`s — seção [7.2](#72-estrutura-do-objeto-workloadselector-configpatches-e-match)); confundir a cadeia `SIDECAR_INBOUND` com a `SIDECAR_OUTBOUND` (Lab 12); ou assumir que `jwt_authn`/`rbac` sempre existem quando, na verdade, só existem se houver `RequestAuthentication`/`AuthorizationPolicy` aplicável àquele workload.

Fonte: [Envoy / Response Code Details](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/response_code_details), [Envoy / RBAC filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/rbac_filter), [Istio / Debugging Envoy and Istiod](https://istio.io/v1.30/docs/ops/diagnostic-tools/proxy-cmd/).

---

# 4. Comandos práticos (kind + Istio deste repositório)

## 4.1. Pré-requisitos

Siga [INSTALL_REQUIREMENTS.md](INSTALL_REQUIREMENTS.md) e [day1.md](day1.md) para: instalar Docker, criar o cluster `kind`, instalar Helm e MetalLB, e instalar o Istio no modo `ambient` (ou [OLD/day1-sidecar.md](OLD/day1-sidecar.md) para o modo `sidecar`). As variáveis de ambiente usadas abaixo seguem o padrão já definido nesses arquivos:

```bash
export ISTIO_RELEASE=1.30
export VERSION_ISTIO="${ISTIO_RELEASE}.3"
export MY_NAMESPACE='myapp'
```

## 4.2. istioctl: comandos essenciais

Fonte: [Istio / Debugging Envoy and Istiod](https://istio.io/v1.30/docs/ops/diagnostic-tools/proxy-cmd/), [Istio / istioctl reference](https://istio.io/v1.30/docs/reference/commands/istioctl/).

```bash
# Visão geral de sincronização de TODOS os proxies da malha (sidecars, gateways, ztunnel, waypoints)
istioctl proxy-status

# Diff entre a config que o Envoy tem e a que o istiod enviaria (detecta drift)
istioctl proxy-status <pod>.<namespace>

# Clusters conhecidos por um proxy específico
istioctl proxy-config cluster -n $MY_NAMESPACE <pod>

# Listeners configurados em um proxy
istioctl proxy-config listener <pod>.<namespace> --port 15006 -o json

# Tabela de rotas de um proxy
istioctl proxy-config route <pod>.<namespace> --name 8080 -o json

# Endpoints saudáveis conhecidos por um cluster específico
istioctl proxy-config endpoint <pod>.<namespace> \
  --cluster "outbound|9080||reviews.$MY_NAMESPACE.svc.cluster.local"

# Segredos/certificados mTLS carregados por um proxy (SDS)
istioctl proxy-config secret <pod>.<namespace>

# Configuração de bootstrap enviada ao Envoy na inicialização
istioctl proxy-config bootstrap -n istio-system <pod-gateway>

# Ajusta o nível de log de um Envoy em tempo real (ex.: para debug)
istioctl proxy-config log <pod>.<namespace> --level debug

# Valida a configuração de todo o cluster (ou de arquivos locais) em busca de erros comuns
istioctl analyze --all-namespaces
istioctl analyze files/bookinfo-gateway.yaml

# Verifica pré-requisitos antes de instalar/atualizar o Istio
istioctl x precheck

# Abre o dashboard de admin do Envoy (via port-forward automático para a porta 15000)
istioctl dashboard envoy <pod>.<namespace>
```

## 4.3. Admin API do Envoy (porta 15000)

Fonte: [Envoy / Administration interface](https://www.envoyproxy.io/docs/envoy/latest/operations/admin).

```bash
# Port-forward manual até a porta de administração do Envoy dentro do Pod
kubectl -n $MY_NAMESPACE port-forward pod/<pod> 15000:15000

# Em outro terminal:

# Dump completo da configuração ativa (bootstrap + xDS já recebido)
curl -s http://localhost:15000/config_dump | jq .

# Estado dos clusters (endpoints descobertos, health, circuit breaker, outlier detection)
curl -s http://localhost:15000/clusters

# Listeners abertos e a porta real (útil quando configurado como porta 0)
curl -s http://localhost:15000/listeners

# Todas as estatísticas (também disponível em formato prometheus)
curl -s http://localhost:15000/stats
curl -s http://localhost:15000/stats/prometheus

# Informações do processo (versão do Envoy, uptime, estado)
curl -s http://localhost:15000/server_info

# Liveness (200 = LIVE, 503 = não pronto)
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:15000/ready

# Ajustar nível de log em tempo real (ex.: habilitar debug no componente "router")
curl -s -X POST "http://localhost:15000/logging?router=debug"
```

## 4.4. Comandos específicos do modo Ambient (ztunnel e waypoint)

Fonte: [Istio / Troubleshoot connectivity issues with ztunnel](https://istio.io/latest/docs/ambient/usage/troubleshoot-ztunnel/), [Istio / Troubleshoot issues with waypoints](https://istio.io/latest/docs/ambient/usage/troubleshoot-waypoint/).

```bash
# Lista os workloads que o ztunnel do nó conhece/gerencia (e o waypoint associado a cada um, se houver)
istioctl ztunnel-config workload

# Mostra qual waypoint (se algum) está associado a cada serviço
istioctl ztunnel-config service

# Certificados carregados pelo ztunnel
export ZTUNNEL=$(kubectl get pods -n istio-system -l app=ztunnel -o jsonpath='{.items[0].metadata.name}')
istioctl ztunnel-config certificate "$ZTUNNEL".istio-system

# Dump bruto da configuração do ztunnel (via admin API, porta 15000, dentro do próprio Pod)
kubectl debug -it "$ZTUNNEL" -n istio-system --image=curlimages/curl -- \
  curl -s localhost:15000/config_dump

# Logs do ztunnel, filtrando eventos de tráfego inbound/outbound
kubectl -n istio-system logs -l app=ztunnel | grep -E "inbound|outbound"

# Criar um waypoint para um namespace/serviço (recurso Gateway com gatewayClassName istio-waypoint)
istioctl waypoint apply -n $MY_NAMESPACE --enroll-namespace

# Status/listagem dos waypoints existentes
istioctl waypoint list -n $MY_NAMESPACE

# Toda a config (clusters/listeners/routes/endpoints/secrets) de um waypoint
istioctl proxy-config all deploy/waypoint -n $MY_NAMESPACE
```

## 4.5. Comandos avançados de troubleshooting do Envoy

Fonte: [Envoy / Administration interface](https://www.envoyproxy.io/docs/envoy/latest/operations/admin), [Istio / Debugging Envoy and Istiod](https://istio.io/v1.30/docs/ops/diagnostic-tools/proxy-cmd/), [Discuss Istio — How do I take an Istio mTLS tcpdump?](https://discuss.istio.io/t/how-do-i-take-istio-mtls-tcpdump/9795).

```bash
export PP_POD=$(kubectl -n $MY_NAMESPACE get pod -l app=productpage -o jsonpath='{.items[0].metadata.name}')
kubectl -n $MY_NAMESPACE port-forward pod/$PP_POD 15000:15000 &

# Filtra estatísticas por regex (ex.: só as métricas do circuit breaker/retry do cluster reviews)
curl -s "http://localhost:15000/stats?filter=reviews.*(circuit_breakers|upstream_rq_retry)" 

# Mesmo filtro, em JSON (útil para script/jq)
curl -s "http://localhost:15000/stats?filter=upstream_rq_pending_overflow&format=json" | jq .

# Estado completo dos clusters em JSON (endpoints, health, circuit breaker, outlier detection)
curl -s "http://localhost:15000/clusters?format=json" | jq .

# Ajusta o nível de log de componentes específicos do Envoy em tempo real (ex.: "router" e "http")
curl -s -X POST "http://localhost:15000/logging?router=debug&http=debug"

# Mesmo resultado, via istioctl (sem precisar do port-forward manual)
istioctl proxy-config log $PP_POD.$MY_NAMESPACE --level router:debug,http:debug

kill %1   # encerra o port-forward em background
```

Captura de pacotes (`tcpdump`) dentro do container `istio-proxy`, útil para confirmar se o handshake TLS/mTLS está de fato ocorrendo entre dois workloads:

```bash
# Tráfego de saída da aplicação sendo interceptado pelo Envoy (outbound, porta 15001)
kubectl -n $MY_NAMESPACE exec $PP_POD -c istio-proxy -- tcpdump -i any -n port 15001

# Tráfego de entrada já decriptado, indo do Envoy para a aplicação (porta da app, ex.: 9080)
kubectl -n $MY_NAMESPACE exec $PP_POD -c istio-proxy -- tcpdump -i any -n port 9080
```

> ATENÇÃO!!! Em muitos ambientes o container `istio-proxy` **não tem** a capability `NET_ADMIN`/`NET_RAW` nem o binário `tcpdump` por padrão, o que faz o comando falhar com erro de permissão. Nesses casos, use um container de debug efêmero compartilhando o namespace de rede do Pod (`kubectl debug -it $PP_POD --image=nicolaka/netshoot --target=istio-proxy`) ou uma ferramenta como o `ksniff` (plugin do `kubectl`).

---

# 5. Observabilidade: métricas com Prometheus, Grafana e Datadog

Esta seção documenta as métricas que o Istio e o Envoy expõem nativamente, como consultá-las via **Prometheus**/**Grafana** (ferramentas já instaladas nos addons do [day1.md](day1.md)) e via **Datadog**, e como usar essas métricas para **detectar** os mesmos tipos de erro reproduzidos nos laboratórios da seção [9](#9-exercícios-de-troubleshooting).

## 5.1. Métricas padrão do Istio

O Istio expõe, por padrão, um conjunto de métricas Prometheus **padronizadas** (chamadas de *standard metrics*), geradas pelo proxy Envoy (sidecar, gateway, ztunnel ou waypoint) para todo tráfego HTTP/HTTP2/gRPC e TCP que passa por ele:

| Métrica | Tipo | Descrição (oficial) |
|---|---|---|
| `istio_requests_total` | COUNTER | "incremented for every request handled by an Istio proxy" — total de requisições HTTP/gRPC |
| `istio_request_duration_milliseconds` | DISTRIBUTION | "measures the duration of requests" — latência das requisições |
| `istio_request_bytes` | DISTRIBUTION | Tamanho do corpo das requisições HTTP |
| `istio_response_bytes` | DISTRIBUTION | Tamanho do corpo das respostas HTTP |
| `istio_tcp_connections_opened_total` | COUNTER | "incremented for every opened connection" — conexões TCP abertas |
| `istio_tcp_connections_closed_total` | COUNTER | "incremented for every closed connection" — conexões TCP fechadas |
| `istio_tcp_sent_bytes_total` | COUNTER | Total de bytes enviados em conexões TCP |
| `istio_tcp_received_bytes_total` | COUNTER | Total de bytes recebidos em conexões TCP |

Labels (tags) documentadas para essas métricas incluem: `response_code`, `grpc_response_status`, `source_workload`, `source_version`, `destination_workload`, `destination_version`, `destination_service`, `connection_security_policy` (indica se o tráfego usou mTLS), `response_flags` (ver [5.4](#54-response-flags-do-envoy)) e `reporter` (indica se quem reportou a métrica foi o proxy de origem — `source` — ou o de destino — `destination`).

> As métricas expostas (e suas tags) podem ser customizadas/reduzidas pelo recurso `Telemetry` (campo `metricsOverrides`/`tagOverride`).

Fonte: [Istio / Standard Metrics Reference](https://istio.io/v1.30/docs/reference/config/metrics/).

## 5.2. Estatísticas internas do Envoy

Além das métricas "padrão" do Istio, o **Envoy** expõe centenas de estatísticas internas de baixo nível (via `/stats` na Admin API, seção [4.3](#43-admin-api-do-envoy-porta-15000)), organizadas por cluster, listener e servidor. As mais úteis para troubleshooting:

| Métrica (por cluster, prefixo `cluster.<nome>.`) | Significado |
|---|---|
| `upstream_cx_total` | Total de conexões abertas para o cluster |
| `upstream_cx_active` | Conexões atualmente ativas |
| `upstream_cx_connect_fail` | Total de falhas ao conectar no upstream |
| `upstream_cx_connect_timeout` | Total de timeouts ao tentar conectar |
| `upstream_cx_overflow` | Incrementado quando o **circuit breaker de conexões** é excedido |
| `upstream_rq_total` | Total de requisições enviadas ao cluster |
| `upstream_rq_pending_overflow` | Incrementado quando o **circuit breaker de requisições pendentes** é excedido (fila cheia) |
| `upstream_rq_timeout` | Total de requisições que sofreram timeout aguardando resposta |
| `upstream_rq_retry` | Total de tentativas de retry |
| `upstream_rq_<2xx\|3xx\|4xx\|5xx>` | Respostas agregadas por família de código HTTP |
| `circuit_breakers.<priority>.cx_open` / `.rq_open` | `1` quando o circuit breaker de conexões/requisições está **aberto** para essa prioridade |
| `circuit_breakers.<priority>.remaining_cx` / `.remaining_pending` / `.remaining_rq` / `.remaining_retries` | Quantos recursos ainda restam antes do circuit breaker abrir (requer `track_remaining: true`) |
| `outlier_detection.ejections_active` | Número de endpoints atualmente **ejetados** do load balancing por outlier detection |
| `outlier_detection.ejections_enforced_total` | Total de ejeções aplicadas (por tipo de anomalia) |

Fonte: [Envoy / Cluster statistics](https://github.com/envoyproxy/envoy/blob/main/docs/root/configuration/upstream/cluster_manager/cluster_stats.rst), [Envoy / Circuit breaking](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/circuit_breaking), [Istio / Envoy Statistics](https://istio.io/latest/docs/ops/configuration/telemetry/envoy-stats/).

> Por padrão, o Istio configura o Envoy para registrar **apenas um subconjunto mínimo** de estatísticas (`cluster_manager`, `listener_manager`, `server`, `cluster.xds-grpc`), para reduzir o overhead de CPU/memória dos proxies. Para habilitar métricas adicionais (ex.: todas as de `outlier_detection` e `upstream_cx_`/`upstream_rq_retry`), use `proxyStatsMatcher.inclusionRegexps`/`inclusionSuffixes` no `IstioOperator`, ou a anotação `proxy.istio.io/config` por Pod.

## 5.3. Métricas do control plane (istiod) e tempo de convergência

O **istiod** expõe suas próprias métricas Prometheus na porta **15014** (`/metrics`). As mais relevantes para saúde do control plane:

| Métrica | Significado |
|---|---|
| `pilot_proxy_convergence_time` (histograma) | Tempo (em segundos) que o istiod leva para propagar uma nova configuração para os proxies — o **tempo de convergência** da malha. Valores altos no p99 indicam gargalo no control plane. |
| `pilot_xds_pushes` | Contador de mensagens xDS enviadas (com label de tipo/erro) — mostra a frequência de atualizações de configuração e falhas de push. |
| `pilot_xds` | Gauge com o número de conexões (proxies) atualmente conectadas ao servidor xDS do istiod. |

Uma consulta PromQL comum para acompanhar a "cauda" do tempo de convergência (p99, em milissegundos) é:

```promql
histogram_quantile(0.99, sum by (le) (rate(pilot_proxy_convergence_time_bucket[5m]))) * 1000
```

Segundo a documentação oficial, o consumo de CPU/memória do istiod escala com três fatores: **"the rate of deployment changes"**, **"the rate of configuration changes"** e **"the number of proxies connecting to Istiod"** — ou seja, picos de convergência lenta tendem a coincidir com deploys em massa, mudanças grandes de configuração (ex.: um `VirtualService` aplicado em um namespace com muitos serviços) ou crescimento no número de proxies conectados.

O livro *Istio in Action* (capítulo 11, "Performance-tuning the control plane") detalha esse modelo com o **método dos "quatro sinais dourados" do control plane** (latência, saturação, tráfego e erros) e métricas adicionais úteis para diagnosticar **qual** dos três fatores acima está causando o gargalo:

| Métrica/variável | Categoria | Significado |
|---|---|---|
| `pilot_proxy_queue_time` | Latência | Tempo que um push fica **na fila** antes de começar a ser processado (parte do tempo total de convergência que `pilot_proxy_convergence_time` não separa) |
| `container_cpu_usage_seconds_total` / `process_cpu_seconds_total` | Saturação | CPU consumida pelo `istiod` — quando ultrapassa ~90%, as atualizações começam a enfileirar e atrasar |
| `pilot_inbound_updates` | Tráfego (entrada) | Contagem de atualizações de configuração recebidas pelo istiod |
| `pilot_push_triggers` | Tráfego (entrada) | Total de eventos que dispararam um push, por tipo (`service`, `endpoint`, `config`) |
| `pilot_services` | Tráfego (entrada) | Número de serviços conhecidos pelo pilot — quanto maior, mais processamento por evento |
| `pilot_total_xds_rejects` (+ `pilot_xds_eds_reject`/`_lds_reject`/`_rds_reject`/`_cds_reject`) | Erros | Total de pushes de configuração **rejeitados** pelos proxies (e o subconjunto por tipo de API xDS) |
| `pilot_xds_write_timeout` | Erros | Timeouts ao tentar enviar um push |
| `pilot_xds_push_context_errors` | Erros | Erros do istiod ao **gerar** a configuração do Envoy (geralmente bugs no Pilot) |

Para **reduzir** o tempo de convergência quando o gargalo é confirmado, o livro recomenda, em ordem de custo/benefício: (1) definir recursos `Sidecar` (seção [8.4](#84-sidecar-restringindo-com-quais-endereços-um-workload-se-relaciona)) para reduzir o tamanho da configuração gerada por workload; (2) usar `discoverySelectors` no `meshConfig` para o istiod **ignorar** eventos de namespaces irrelevantes (por padrão, o control plane observa recursos em **todos** os namespaces do cluster); (3) ajustar o *debounce*/*throttle* de push via variáveis de ambiente do `istiod` — `PILOT_DEBOUNCE_AFTER` (padrão 100 ms — tempo para agrupar eventos que chegam em sequência antes de gerar um push), `PILOT_DEBOUNCE_MAX` (padrão 10 s — teto para esse agrupamento), `PILOT_ENABLE_EDS_DEBOUNCE` (padrão `true`) e `PILOT_PUSH_THROTTLE` (padrão 100 pushes concorrentes); (4) escalar o `istiod` **out** (mais réplicas, quando o gargalo é volume de proxies) ou **up** (mais CPU/memória por réplica, quando o gargalo é volume de configuração a processar).

```yaml
# Exemplo (IstioOperator) citado no livro para reduzir o escopo de descoberta do control plane
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  namespace: istio-system
spec:
  meshConfig:
    discoverySelectors:
      - matchLabels:
          istio-discovery: enabled
```

Fonte: [Istio / Performance and Scalability](https://istio.io/v1.30/docs/ops/deployment/performance-and-scalability/), [Datadog — How to Monitor Istiod](https://www.datadoghq.com/blog/how-to-monitor-istiod/), POSTA & MALOKU, *Istio in Action* (Manning, 2022), capítulo 11.

## 5.4. Response Flags do Envoy

O Envoy anota cada entrada do **access log** com "response flags" (`%RESPONSE_FLAGS%`), curtas siglas que explicam a causa de uma resposta de erro. O Istio também expõe essa informação como a tag `response_flags` da métrica `istio_requests_total` (seção [5.1](#51-métricas-padrão-do-istio)), o que permite filtrar erros por causa diretamente no Prometheus/Grafana. As mais comuns:

| Flag | Nome | Significado |
|---|---|---|
| `UH` | NoHealthyUpstream | Nenhum host saudável no cluster upstream (além do código 503) |
| `UF` | UpstreamConnectionFailure | Falha ao conectar no upstream (além do 503) |
| `UO` | UpstreamOverflow | **Circuit breaker** do cluster foi excedido (além do 503) |
| `NR` | NoRouteFound | Nenhuma rota configurada, ou nenhuma filter chain correspondente |
| `URX` | UpstreamRetryLimitExceeded | Limite de retries/tentativas de conexão upstream esgotado |
| `UC` | UpstreamConnectionTermination | Conexão upstream terminada (além do 503) |
| `UT` | UpstreamRequestTimeout | Requisição upstream sofreu timeout (além do 504) |
| `LR` | LocalReset | Reset local da conexão (além do 503) |
| `DC` | DownstreamConnectionTermination | Conexão downstream terminada |
| `DI` | DelayInjected | Requisição atrasada por *fault injection* |
| `FI` | FaultInjected | Requisição abortada por *fault injection* |
| `RL` | RateLimited | Requisição limitada por rate limiting local (além do 429) |

Múltiplas flags podem aparecer juntas separadas por vírgula (ex.: `UC,URX` = a conexão upstream foi terminada **e** os retries se esgotaram).

Fonte: [Envoy / Substitution Formatter — Response Flags](https://www.envoyproxy.io/docs/envoy/latest/configuration/advanced/substitution_formatter), [Istio / Standard Metrics Reference](https://istio.io/v1.30/docs/reference/config/metrics/).

## 5.5. Prometheus + Grafana

O [day1.md](day1.md) já instala o Prometheus e o Grafana via `samples/addons` do próprio Istio. O Istio disponibiliza dashboards Grafana **oficiais e pré-configurados**, identificados pelo ID no grafana.com:

| Dashboard | ID | Função (oficial) |
|---|---|---|
| Mesh Dashboard | 7639 | "Provides an overview of all services in the mesh" |
| Service Dashboard | 7636 | "Provides a detailed breakdown of metrics for a service" |
| Workload Dashboard | 7630 | "Provides a detailed breakdown of metrics for a workload" |
| Performance Dashboard | 11829 | "Monitors the resource usage of the mesh" |
| Control Plane Dashboard | 7645 | "Monitors the health and performance of the control plane" |
| WASM Extension Dashboard | 13277 | Visão geral de extensões WebAssembly carregadas na malha |
| Ztunnel Dashboard | 21306 | Monitora o componente `ztunnel` do modo ambient |

```bash
# Alternativa de instalação via kubectl (equivalente ao que já é feito em day1.md)
kubectl apply -f "https://raw.githubusercontent.com/istio/istio/release-1.30/samples/addons/grafana.yaml"

# Acesso local (já documentado em day1.md)
kubectl -n istio-system port-forward service/grafana 3000:3000
kubectl -n istio-system port-forward service/prometheus 9090:9090
```

Exemplos de consultas **PromQL** úteis (usar em Grafana ou diretamente em `http://localhost:9090`):

```promql
# Taxa de erros (5xx) do serviço "reviews", por segundo, nos últimos 5 minutos
sum(rate(istio_requests_total{destination_service_name="reviews", response_code=~"5.*"}[5m]))

# Latência p99 (em ms) das requisições para "reviews", vindas do productpage
histogram_quantile(0.99,
  sum(rate(istio_request_duration_milliseconds_bucket{destination_service_name="reviews"}[5m])) by (le)
)

# Requisições rejeitadas por falta de upstream saudável (UH) ou circuit breaker (UO)
sum(rate(istio_requests_total{response_flags=~"UH|UO"}[5m])) by (destination_service_name, response_flags)

# Conexões atualmente ejetadas por outlier detection, por cluster do Envoy
envoy_cluster_outlier_detection_ejections_active
```

Fonte: [Istio / Grafana](https://istio.io/v1.30/docs/ops/integrations/grafana/), [Istio / Visualizing Metrics with Grafana](https://istio.io/latest/docs/tasks/observability/metrics/using-istio-dashboard/).

## 5.6. Datadog

O **Datadog Agent** tem uma integração oficial com Istio que coleta métricas de dois pontos diferentes, via *OpenMetrics* (Prometheus):

- **Data plane (Envoy)**: coleta do endpoint `http://<pod-ip>:15020/stats/prometheus` de cada sidecar `istio-proxy`, gerando métricas com prefixo **`istio.mesh.*`** (latência, tamanho de request/response, conexões TCP — os equivalentes às métricas padrão da seção [5.1](#51-métricas-padrão-do-istio)).
- **Control plane (istiod)**: coleta do endpoint `http://<istiod-pod-ip>:15014/metrics`, gerando métricas com prefixo **`istio.pilot.*`** (convergência de proxies, conflitos de listener, eventos xDS — equivalentes à seção [5.3](#53-métricas-do-control-plane-istiod-e-tempo-de-convergência)) e, em instalações legadas, `istio.mixer.*`/`istio.galley.*`/`istio.citadel.*`.

Exemplo mínimo de configuração do check `istio` do Datadog Agent (`istio.d/conf.yaml`) para coletar as métricas de data plane:

```yaml
ad_identifiers:
  - proxyv2
instances:
  - use_openmetrics: true
    istio_mesh_endpoint: http://%%host%%:15020/stats/prometheus
```

Para o control plane, o Datadog Agent descobre o `istiod` via *autodiscovery annotation* no próprio Deployment, apontando para `http://%%host%%:15014/metrics`. Em instalações no modo **ambient** (Istio 1.24+), a integração também aceita endpoints específicos para `ztunnel_endpoint` e `waypoint_endpoint`.

Também é possível monitorar cada proxy **individualmente** (não apenas as métricas agregadas do mesh) habilitando a [integração Envoy do Datadog](https://github.com/DataDog/integrations-core/tree/master/envoy#istio) no mesmo Agent.

Fonte: [Datadog / Istio Integration](https://docs.datadoghq.com/integrations/istio/), [Datadog — How to Monitor Istiod](https://www.datadoghq.com/blog/how-to-monitor-istiod/).

```mermaid
flowchart LR
    subgraph PodApp["Pod da aplicacao"]
        ENVOY["Envoy sidecar<br/>:15090/stats/prometheus (Envoy)<br/>:15020/stats/prometheus (merged)"]
    end

    subgraph CP["istio-system"]
        ISTIOD["istiod<br/>:15014/metrics"]
    end

    PROM["Prometheus<br/>(scrape :15020 e :15014)"]
    GRAFANA["Grafana<br/>(dashboards oficiais do Istio)"]
    DDAGENT["Datadog Agent<br/>(openmetrics check istio)"]
    DD["Datadog<br/>(istio.mesh.* / istio.pilot.*)"]

    ENVOY --> PROM
    ISTIOD --> PROM
    PROM --> GRAFANA

    ENVOY --> DDAGENT
    ISTIOD --> DDAGENT
    DDAGENT --> DD
```

## 5.7. Como detectar erros comuns através das métricas

Tabela de correlação entre **sintoma observável** e **métrica/consulta** a inspecionar, cruzando com os laboratórios da seção [9](#9-exercícios-de-troubleshooting):

| Sintoma | Métrica/consulta a checar | Causa provável |
|---|---|---|
| Latência alta (cauda/p99) em um serviço | `histogram_quantile(0.99, ...istio_request_duration_milliseconds_bucket...)` (seção [5.5](#55-prometheus--grafana)) | Fault injection/delay (Lab 9), upstream lento, ou sobrecarga do worker thread do Envoy |
| Picos de erro 503 | `sum(rate(istio_requests_total{response_code="503"}[5m])) by (response_flags)` | Combinar com `response_flags` (seção [5.4](#54-response-flags-do-envoy)) para achar a causa exata: `UH`, `UF`, `UO`, `UC`, `NR` |
| `response_flags="UH"` (No Healthy Upstream) | `istio_requests_total{response_flags="UH"}` | Host/serviço de destino incorreto ou sem endpoints saudáveis (Lab 3) |
| `response_flags="UO"` ou `cluster.<x>.upstream_cx_overflow` crescendo | `envoy_cluster_upstream_cx_overflow` / `envoy_cluster_upstream_rq_pending_overflow` (seção [5.2](#52-estatísticas-internas-do-envoy)) | **Circuit breaker** atingido — limite de conexões/requisições pendentes configurado baixo demais para o tráfego real (Lab 10) |
| Requisições saindo com `URX` | `istio_requests_total{response_flags="URX"}` | Retries configurados esgotados, ou `outlier_detection` ejetando o(s) único(s) endpoint(s) saudável(is) |
| `outlier_detection.ejections_active` > 0 | `envoy_cluster_outlier_detection_ejections_active` | Um ou mais endpoints do cluster estão respondendo com erro e foram ejetados temporariamente do load balancing |
| Configuração aplicada mas Envoy não reflete a mudança | `pilot_proxy_convergence_time` (p99) + `istioctl proxy-status` (seção [5.3](#53-métricas-do-control-plane-istiod-e-tempo-de-convergência), Lab 4, Lab 11) | Tempo de convergência do control plane alto (deploy em massa, istiod sobrecarregado ou `NetworkPolicy` bloqueando a porta 15012) |
| Tráfego chegando sem mTLS quando deveria ser `STRICT` | `istio_requests_total{connection_security_policy!="mutual_tls", reporter="destination"}` (valores possíveis: `mutual_tls`, `none`, `unknown`) | Cliente fora da mesh, ou `PeerAuthentication` mal configurada (Lab 2) |
| Conexões TCP abrindo e fechando rapidamente sem dados | `istio_tcp_connections_opened_total` vs `istio_tcp_connections_closed_total` (razão próxima de 1 com pouca duração) | Falha de handshake mTLS ou health check derrubando a conexão |

Fonte: seções [5.1](#51-métricas-padrão-do-istio) a [5.6](#56-datadog) deste guia (todas com suas próprias referências oficiais).

---

# 6. Logs do Istio e do Envoy

O Istio tem **dois mecanismos de log completamente distintos**: o log de **componentes** (mensagens internas de diagnóstico do processo, tanto do Envoy quanto do istiod) e o **access log** (uma linha por requisição/conexão processada pelo Envoy, usado nas seções [5.4](#54-response-flags-do-envoy) e nos Labs). Esta seção cobre os dois.

## 6.1. Nível de log do Envoy/istio-proxy (componentes e debug)

O Envoy organiza suas mensagens de log internas por **componente** (chamados de *loggers*), cada um com seu próprio nível. Os níveis válidos são, do mais verboso ao mais silencioso: **`trace, debug, info, warning, error, critical, off`**. Alguns dos componentes disponíveis: `admin`, `assert`, `backtrace`, `client`, `config`, `connection`, `conn_handler`, `filter`, `forward_proxy`, `grpc`, `http`, `http2`, `init`, `io`, `jwt`, `main`, `router`, `upstream`, entre outros.

```bash
export PP_POD=$(kubectl -n $MY_NAMESPACE get pod -l app=productpage -o jsonpath='{.items[0].metadata.name}')

# Lista o nível atual de TODOS os loggers do proxy
istioctl proxy-config log $PP_POD.$MY_NAMESPACE

# Habilita "modo debug" em TODOS os componentes de uma vez (bem verboso!)
istioctl proxy-config log $PP_POD.$MY_NAMESPACE --level debug

# Habilita debug só nos componentes mais úteis para troubleshooting de roteamento/mTLS
istioctl proxy-config log $PP_POD.$MY_NAMESPACE --level http:debug,router:debug,connection:debug,upstream:debug

# Volta ao nível padrão (desabilita o modo debug)
istioctl proxy-config log $PP_POD.$MY_NAMESPACE --level info
```

> ATENÇÃO!!! Essa mudança é **temporária e não persistente**: ela é feita chamando a Admin API do proxy em tempo real (o mesmo endpoint `POST /logging` da seção [4.5](#45-comandos-avançados-de-troubleshooting-do-envoy)) e volta ao nível configurado na inicialização caso o Pod seja reiniciado.

## 6.2. Nível de log do istiod

O **istiod** (um binário Go) organiza suas próprias mensagens de log em **escopos** (*scopes*) — exemplos citados na documentação oficial: `default`, `ads` (Aggregated Discovery Service — o que processa e envia xDS), `model`, `authorization`, `ca`, `cache`. Os níveis válidos para istiod são **`none, error, warn, info, debug`**, com **`info`** como padrão de todos os escopos.

```bash
# Ver/alterar o nível de log do istiod é feito pela flag --log_output_level na inicialização,
# por exemplo editando o Deployment ou (se instalado via istioctl/IstioOperator) o profile de instalação:
#   --log_output_level=default:info,ads:debug,authorization:debug,model:warn

# Alternativa em tempo real, sem reiniciar o istiod, via interface ControlZ:
export ISTIOD=$(kubectl get pods -n istio-system -l app=istiod -o jsonpath='{.items[0].metadata.name}')
kubectl -n istio-system port-forward "$ISTIOD" 9876:9876
# Acesse http://localhost:9876 no navegador para ver/ajustar os escopos de log em tempo real
```

Fonte: [Istio / Component logging](https://istio.io/v1.30/docs/ops/diagnostic-tools/component-logging/).

## 6.3. Access log: habilitar e customizar campos

O **access log** é habilitado/desabilitado e customizado de duas formas (a Telemetry API, seção [8.1](#81-telemetry-habilitandocustomizando-métricas-por-workload), é a forma recomendada e granular; o `meshConfig`, seção [8.2](#82-meshconfig-configurações-globais-da-malha), é a forma global):

```yaml
# Forma recomendada: Telemetry API (pode ser aplicada por namespace/workload, não só mesh-wide)
apiVersion: telemetry.istio.io/v1
kind: Telemetry
metadata:
  name: mesh-default
  namespace: istio-system
spec:
  accessLogging:
    - providers:
        - name: envoy
```

```bash
# Forma global via meshConfig (ver seção 8.2) — helm upgrade, já que este repo instala via Helm:
helm -n istio-system upgrade istiod istio/istiod --version $VERSION_ISTIO --reuse-values \
  --set meshConfig.accessLogFile=/dev/stdout \
  --set meshConfig.accessLogEncoding=JSON
```

O **formato padrão (TEXT)** do access log do Istio é, literalmente:

```
[%START_TIME%] "%REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% %PROTOCOL%" %RESPONSE_CODE% %RESPONSE_FLAGS% %RESPONSE_CODE_DETAILS% %CONNECTION_TERMINATION_DETAILS%
"%UPSTREAM_TRANSPORT_FAILURE_REASON%" %BYTES_RECEIVED% %BYTES_SENT% %DURATION% %RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)% "%REQ(X-FORWARDED-FOR)%" "%REQ(USER-AGENT)%" "%REQ(X-REQUEST-ID)%"
"%REQ(:AUTHORITY)%" "%UPSTREAM_HOST%" %UPSTREAM_CLUSTER_RAW% %UPSTREAM_LOCAL_ADDRESS% %DOWNSTREAM_LOCAL_ADDRESS% %DOWNSTREAM_REMOTE_ADDRESS% %REQUESTED_SERVER_NAME% %ROUTE_NAME%
```

Que produz uma linha como esta (exemplo oficial):

```text
[2020-11-25T21:26:18.409Z] "GET /status/418 HTTP/1.1" 418 - via_upstream - "-" 0 135 4 4 "-" "curl/7.73.0-DEV" "84961386-6d84-929d-98bd-c5aee93b5c88" "httpbin:8000" "10.44.1.27:80" outbound|8000||httpbin.foo.svc.cluster.local 10.44.1.23:37652 10.0.45.184:8000 10.44.1.23:46520 - default
```

**Adicionando informação ao log**: para customizar o formato (adicionar/remover campos), defina `providers[].name: envoy` no `Telemetry` e um `extensionProviders` do tipo `envoyFileAccessLog` no `meshConfig` (seção [8.2](#82-meshconfig-configurações-globais-da-malha)) com um `logFormat` próprio, reutilizando quaisquer *command operators* do Envoy (`%...%`), incluindo cabeçalhos de requisição/resposta arbitrários via `%REQ(<nome-do-header>)%`/`%RESP(<nome-do-header>)%` — por exemplo, para logar um header de tracing customizado ou um `X-Tenant-Id` de aplicação multi-tenant.

Fonte: [Istio / Envoy Access Logs](https://istio.io/v1.30/docs/tasks/observability/logs/access-log/), [Envoy / Access logging](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage).

## 6.4. Latência e outras métricas visíveis apenas no log

Diferente das métricas agregadas do Prometheus (seção [5](#5-observabilidade-métricas-com-prometheus-grafana-e-datadog)), o access log registra a **latência de uma requisição individual**, o que é essencial para depurar um caso específico (ex.: um `X-Request-Id` reportado por um usuário). Os campos de tempo do Envoy, definidos oficialmente como:

| Campo | Significado (oficial) |
|---|---|
| `%START_TIME%` | Timestamp de início da requisição (t0) |
| `%REQUEST_DURATION%` | "Total duration in milliseconds of the request from the start time to the last byte of the request received from the downstream" — quanto tempo levou para receber a requisição inteira |
| `%RESPONSE_DURATION%` | "Total duration in milliseconds of the request from the start time to the first byte read from the upstream host" — tempo até o primeiro byte de resposta do upstream (*time to first byte*) |
| `%RESPONSE_TX_DURATION%` | "Total duration in milliseconds of the request from the first byte read from the upstream host to the last byte sent downstream" — tempo gasto **transmitindo** a resposta de volta |
| `%DURATION%` | "Total duration in milliseconds of the request from the start time to the last byte out" — duração **total**, o campo usado por padrão no log do Istio |
| `%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%` | Tempo gasto **apenas dentro do serviço upstream** (extraído do header de resposta que o próprio Envoy upstream adiciona) — é o campo que aparece por padrão no access log do Istio, logo após `%DURATION%` |

```mermaid
flowchart LR
    T0["t0: START_TIME"] -->|REQUEST_DURATION| T1["t1: ultimo byte da requisicao recebido"]
    T0 -->|RESPONSE_DURATION - TTFB| T2["t2: primeiro byte de resposta do upstream"]
    T2 -->|RESPONSE_TX_DURATION| T3["t3: ultimo byte de resposta enviado ao downstream"]
    T0 -->|"DURATION (total)"| T3
```

Na prática, para investigar **por que uma requisição específica foi lenta**, compare `%DURATION%` (tempo total no Envoy) com `%RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)%` (tempo só no serviço upstream): se a diferença for grande, o tempo foi gasto no próprio Envoy/rede (fila de conexão, TLS handshake, retries) e não na aplicação de destino.

Fonte: [Envoy / Substitution Formatter](https://www.envoyproxy.io/docs/envoy/latest/configuration/advanced/substitution_formatter), [Istio / Envoy Access Logs](https://istio.io/v1.30/docs/tasks/observability/logs/access-log/).

## 6.5. Comandos de troubleshooting usando logs

```bash
# Acompanhar o access log de um workload em tempo real
kubectl -n $MY_NAMESPACE logs -f deploy/productpage-v1 -c istio-proxy

# Filtrar apenas erros (resposta >= 400) usando grep no formato TEXT padrão
kubectl -n $MY_NAMESPACE logs deploy/productpage-v1 -c istio-proxy | grep -E ' [45][0-9]{2} '

# Se accessLogEncoding=JSON (seção 6.3), usar jq para filtrar/formatar
kubectl -n $MY_NAMESPACE logs deploy/productpage-v1 -c istio-proxy | jq -r 'select(.response_code >= 400)'

# Correlacionar uma requisição específica pelo X-Request-Id em TODOS os proxies do caminho
# (productpage -> reviews -> ratings), útil para reconstruir o trace manualmente sem Jaeger
REQ_ID="84961386-6d84-929d-98bd-c5aee93b5c88"
for d in productpage-v1 reviews-v1 ratings-v1; do
  echo "== $d =="
  kubectl -n $MY_NAMESPACE logs deploy/$d -c istio-proxy | grep "$REQ_ID"
done

# Aumentar o log level (seção 6.1) só durante a reprodução do problema, depois voltar ao normal
istioctl proxy-config log $PP_POD.$MY_NAMESPACE --level debug
# ... reproduzir o problema ...
istioctl proxy-config log $PP_POD.$MY_NAMESPACE --level info
```

---

# 7. EnvoyFilter

## 7.1. O que é e quando usar

`EnvoyFilter` é o recurso do Istio que permite **customizar diretamente a configuração do Envoy gerada pelo istiod** — ele opera em um nível mais baixo que `VirtualService`/`DestinationRule`, manipulando a estrutura interna de listeners, filter chains, filtros HTTP/rede, rotas e clusters. A própria documentação recomenda cautela: é uma API de baixo nível, acoplada a detalhes internos do Envoy que podem mudar entre versões, e deve ser usada **apenas quando não existe outro recurso Istio de mais alto nível** (`VirtualService`, `DestinationRule`, `Telemetry`, etc.) que resolva a necessidade.

## 7.2. Estrutura do objeto: workloadSelector, configPatches e match

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: <nome>
  namespace: <namespace>
spec:
  workloadSelector:        # opcional; se omitido, aplica a todos os workloads do namespace
    labels:
      app: productpage
  configPatches:
    - applyTo: HTTP_FILTER              # ver seção 7.3
      match:                            # condições para a patch ser aplicada
        context: SIDECAR_OUTBOUND       # ANY | SIDECAR_INBOUND | SIDECAR_OUTBOUND | GATEWAY
        listener:
          portNumber: 9080
          filterChain:
            filter:
              name: "envoy.filters.network.http_connection_manager"
              subFilter:
                name: "envoy.filters.http.router"
      patch:                            # ver seção 7.4
        operation: INSERT_BEFORE
        value: { ... }                  # configuração JSON/proto do objeto Envoy
```

- **`workloadSelector`**: "criteria used to select the specific set of pods/VMs on which this patch configuration should be applied" — sem ele, a patch atinge todos os workloads do namespace; se o `EnvoyFilter` estiver no **namespace raiz** (o namespace de instalação do Istio, ex.: `istio-system`), ele se estende a **todos os namespaces** da malha.
- **`configPatches[].match`** (`EnvoyConfigObjectMatch`): define **onde** procurar no Envoy antes de aplicar a patch — `context` (o Envoy sendo gerado é de um sidecar inbound, outbound ou de um gateway), `listener` (porta, nome, filter chain por SNI/protocolo de transporte/aplicação), `filter`/`subFilter` (filtro de rede/HTTP específico pelo nome, ex.: `envoy.filters.network.http_connection_manager`), `routeConfiguration` e `cluster` (porta, serviço, subset).
- **Ordem de aplicação**: quando há múltiplos `EnvoyFilter`, a ordem é — primeiro todos os `EnvoyFilter` do namespace raiz, depois todos os do namespace do workload, respeitando a ordem de criação dentro de cada grupo.

## 7.3. `applyTo`: em qual camada da configuração o patch age

| `applyTo` | Camada afetada |
|---|---|
| `LISTENER` | O listener inteiro (ex.: adicionar um novo listener) |
| `FILTER_CHAIN` | Uma filter chain dentro de um listener |
| `NETWORK_FILTER` | Um filtro de rede (L3/L4) dentro de uma filter chain — "modify an existing filter or add a new filter" |
| `HTTP_FILTER` | Um filtro HTTP (L7) dentro do HTTP Connection Manager — mesmo texto oficial de "modify an existing filter or add a new filter" |
| `ROUTE_CONFIGURATION` | A tabela de rotas usada pelo HTTP Connection Manager (só aceita operação `MERGE`) |
| `VIRTUAL_HOST` | Um virtual host dentro de uma `ROUTE_CONFIGURATION` |
| `HTTP_ROUTE` | Uma rota HTTP dentro do virtual host casado |
| `CLUSTER` | Um cluster (saída do CDS) — também usado para **adicionar** novos clusters |
| `EXTENSION_CONFIG` | Configuração de extensão via ECDS (só suportado por HTTP filters) |

```mermaid
flowchart TB
    subgraph L["Listener  -  applyTo: LISTENER"]
        LF["Listener Filters"]
        subgraph FC["Filter Chain  -  applyTo: FILTER_CHAIN"]
            subgraph NF["Network Filters  -  applyTo: NETWORK_FILTER"]
                HCM["HTTP Connection Manager"]
                subgraph HF["HTTP Filter Chain  -  applyTo: HTTP_FILTER"]
                    CUSTOM["ex: envoy.filters.http.compressor (inserido via EnvoyFilter)"]
                    ROUTERF["envoy.filters.http.router"]
                end
            end
        end
    end

    subgraph RC["ROUTE_CONFIGURATION / VIRTUAL_HOST / HTTP_ROUTE"]
        RT["Route table usada pelo HCM (RDS)"]
    end

    subgraph CL["CLUSTER"]
        C["Definicoes de cluster (CDS)"]
    end

    HCM -.->|usa| RT
    ROUTERF -.->|encaminha para| C
```

## 7.4. `patch.operation`: como a requisição/configuração é manipulada

A **requisição em si** é manipulada indiretamente: o `EnvoyFilter` não intercepta a requisição diretamente, ele **insere/substitui/remove filtros** na cadeia que o Envoy executa para cada requisição (a mesma cadeia da seção [3.4](#34-life-of-a-request-o-ciclo-de-vida-de-uma-requisição-no-envoy)) — um filtro inserido via `INSERT_BEFORE` do `router`, por exemplo, processa a requisição **antes** dela ser roteada ao cluster de destino.

| `operation` | Comportamento (oficial) |
|---|---|
| `MERGE` | "Merge the provided config with the generated config using proto merge semantics" |
| `ADD` | "Add the provided config to an existing list (of listeners, clusters, virtual hosts, network filters, or http filters)" |
| `REMOVE` | "Remove the selected object from the list. Does not require a value to be specified" |
| `INSERT_BEFORE` | "Insert before the selected filter or sub filter. If no filter is selected, the specified filter will be inserted at the front of the list" |
| `INSERT_AFTER` | "Insert after the selected filter or sub filter. If no filter is selected, the specified filter will be inserted at the end of the list" |
| `INSERT_FIRST` | "Insert first in the list based on the presence of selected filter or not" |
| `REPLACE` | "Replace contents of a named filter with new contents" — válido apenas para `HTTP_FILTER` e `NETWORK_FILTER` |

## 7.5. Camadas de processamento e onde o EnvoyFilter se encaixa

O diagrama da seção [7.3](#73-applyto-em-qual-camada-da-configuração-o-patch-age) já mostra as camadas; a regra prática mais importante é: **a ordem dos filtros HTTP importa**. Um filtro de autenticação/compressão precisa ser inserido **antes** do `envoy.filters.http.router` (que é sempre o último da cadeia — ver seção [3.2](#32-arquitetura-interna-listener-subsystem-x-cluster-subsystem)), senão a requisição já terá sido roteada ao upstream antes do seu filtro rodar. Por isso o padrão mais comum é `operation: INSERT_BEFORE` combinado com `match.listener.filterChain.filter.subFilter.name: envoy.filters.http.router`.

## 7.6. Exemplo prático: habilitando compressão gzip via EnvoyFilter

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: enable-gzip-compression
  namespace: myapp
spec:
  workloadSelector:
    labels:
      app: productpage
  configPatches:
    - applyTo: HTTP_FILTER
      match:
        context: SIDECAR_INBOUND
        listener:
          filterChain:
            filter:
              name: "envoy.filters.network.http_connection_manager"
              subFilter:
                name: "envoy.filters.http.router"
      patch:
        operation: INSERT_BEFORE
        value:
          name: envoy.filters.http.compressor
          typed_config:
            "@type": type.googleapis.com/envoy.extensions.filters.http.compressor.v3.Compressor
            response_direction_config:
              common_config:
                min_content_length: 100
                content_type:
                  - text/html
                  - application/json
            compressor_library:
              name: text_optimized
              typed_config:
                "@type": type.googleapis.com/envoy.extensions.compression.gzip.compressor.v3.Gzip
```

Isso insere o filtro `envoy.filters.http.compressor` (com a biblioteca `gzip`) imediatamente antes do `router`, na cadeia **inbound** do sidecar do `productpage`: quando o cliente envia `Accept-Encoding: gzip`, o Envoy comprime a resposta antes de enviá-la, removendo `content-length` e adicionando `transfer-encoding: chunked`.

Fonte: [Istio / EnvoyFilter reference](https://istio.io/v1.30/docs/reference/config/networking/envoy-filter/), [Envoy / Compressor filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/compressor_filter).

## 7.7. Comandos de troubleshooting de EnvoyFilter

```bash
# Lista os EnvoyFilters aplicados no cluster/namespace
kubectl get envoyfilter -A

# istioctl analyze não valida profundamente o conteúdo do patch, mas sinaliza conflitos comuns
istioctl analyze -n $MY_NAMESPACE

# A forma MAIS CONFIÁVEL de confirmar que a patch foi aplicada: procurar o filtro
# diretamente na configuração final recebida pelo Envoy (não no YAML aplicado)
istioctl proxy-config listener $PP_POD.$MY_NAMESPACE --port 9080 -o json | grep -A5 compressor

# Ou via Admin API diretamente (seção 4.3), olhando o config_dump completo
kubectl -n $MY_NAMESPACE port-forward $PP_POD 15000:15000 &
curl -s http://localhost:15000/config_dump | grep -B5 -A30 "envoy.filters.http.compressor"
kill %1
```

**Causas mais comuns de um `EnvoyFilter` "não fazer nada":**
- `workloadSelector` não casa com nenhum label do Pod (ou está no namespace errado).
- `match.context` errado (ex.: patch feito para `SIDECAR_INBOUND` mas o tráfego é `SIDECAR_OUTBOUND`, ou vice-versa).
- `match.listener.portNumber`/`filterChain` não corresponde exatamente ao listener real (confira com `istioctl proxy-config listener`, seção [4.2](#42-istioctl-comandos-essenciais)).
- Tentar `INSERT_AFTER`/`REPLACE` referenciando um `subFilter` que não existe com esse nome exato na versão do Envoy empacotada pelo Istio em uso.

Fonte: [Istio / EnvoyFilter reference](https://istio.io/v1.30/docs/reference/config/networking/envoy-filter/), [Istio / Debugging Envoy and Istiod](https://istio.io/v1.30/docs/ops/diagnostic-tools/proxy-cmd/).

---

# 8. Configurações avançadas do Istio

## 8.1. Telemetry: habilitando/customizando métricas por workload

O recurso `Telemetry` (`telemetry.istio.io`) controla **métricas**, **access log** (seção [6.3](#63-access-log-habilitar-e-customizar-campos)) e **tracing** para toda a malha, um namespace ou um workload específico, em ordem de precedência: configuração do workload → do namespace → a global (`mesh-default`, no namespace raiz).

```yaml
# Habilita access log (envoy) e reduz cardinalidade de métricas em todo o mesh
apiVersion: telemetry.istio.io/v1
kind: Telemetry
metadata:
  name: mesh-default
  namespace: istio-system
spec:
  accessLogging:
    - providers:
        - name: envoy
  metrics:
    - providers:
        - name: prometheus
      overrides:
        - match:
            metric: REQUEST_COUNT
            mode: SERVER
          tagOverrides:
            response_code:
              operation: REMOVE           # remove a dimensão response_code (menos cardinalidade)
        - match:
            metric: ALL_METRICS
            mode: CLIENT
          disabled: true                  # desliga métricas "client-side" nesse escopo
```

Campos principais: `selector`/`targetRefs` (no máximo um dos dois, escolhe o(s) workload(s) alvo), `metrics[].providers`, `metrics[].overrides[].match.metric` (`REQUEST_COUNT`, `REQUEST_DURATION`, `REQUEST_SIZE`, `RESPONSE_SIZE`, `TCP_OPENED_CONNECTIONS`, `TCP_CLOSED_CONNECTIONS`, `TCP_SENT_BYTES`, `TCP_RECEIVED_BYTES`, `ALL_METRICS`), `tagOverrides[].operation` (`UPSERT` para adicionar/alterar uma dimensão, `REMOVE` para excluí-la), e `tracing[].randomSamplingPercentage`/`disableSpanReporting`/`customTags`.

Fonte: [Istio / Telemetry API reference](https://istio.io/v1.30/docs/reference/config/telemetry/).

## 8.2. MeshConfig: configurações globais da malha

O `meshConfig` é o "arquivo" central de configuração da malha — na prática, um objeto YAML guardado na chave `mesh:` do `ConfigMap` **`istio`**, no namespace `istio-system`, que o istiod observa (*watch*) e aplica dinamicamente.

```bash
# Ver a configuração atual da malha
kubectl get configmap istio -n istio-system -o yaml
```

> ATENÇÃO!!! É possível editar esse `ConfigMap` diretamente com `kubectl edit`, mas **não é recomendado**: como este repositório instala o Istio via **Helm** ([day1.md](day1.md)), a próxima execução de `helm upgrade` sobrescreve qualquer edição manual. A forma correta e persistente de alterar o `meshConfig` neste repositório é via `--set` no `helm upgrade` do chart `istiod`:

```bash
helm -n istio-system upgrade istiod istio/istiod --version $VERSION_ISTIO --reuse-values \
  --set meshConfig.accessLogFile=/dev/stdout \
  --set meshConfig.outboundTrafficPolicy.mode=REGISTRY_ONLY \
  --set meshConfig.defaultConfig.holdApplicationUntilProxyStarts=true
```

Alguns campos úteis de `meshConfig`:

| Campo | Função |
|---|---|
| `defaultConfig` (`ProxyConfig`) | Configuração padrão aplicada a **todo** sidecar/gateway **no momento da injeção** — mudanças aqui só valem para Pods **recriados** depois da alteração |
| `extensionProviders` | Registra provedores externos (ex.: `envoyOtelAls`/`envoyFileAccessLog` para access log, `prometheus` para métricas, provedores de tracing como Zipkin/OpenTelemetry) que depois são referenciados pelo nome no `Telemetry` (seção [8.1](#81-telemetry-habilitandocustomizando-métricas-por-workload)) |
| `outboundTrafficPolicy.mode` | `ALLOW_ANY` (padrão — tráfego de saída não registrado na malha é permitido) ou `REGISTRY_ONLY` (bloqueia qualquer destino que não seja um `ServiceEntry`/serviço conhecido) |
| `accessLogFile` / `accessLogEncoding` | Ver seção [6.3](#63-access-log-habilitar-e-customizar-campos) |
| `enableTracing` | Habilita a geração de spans/IDs de rastreamento |
| `trustDomain` | Raiz de confiança usada para montar as identidades SPIFFE dos workloads (ex.: `cluster.local`) |
| `defaultServiceExportTo` | Visibilidade padrão dos serviços entre namespaces (`*`, `.`, `~`) |

Fonte: [Istio / MeshConfig reference](https://istio.io/v1.30/docs/reference/config/istio.mesh.v1alpha1/), [howardjohn's blog — Everything you need to know about Istio installation](https://blog.howardjohn.info/posts/istio-install/).

## 8.3. Certificados mTLS: plugando uma CA própria e rotação

Por padrão, o **istiod** atua como sua própria CA (Autoridade Certificadora) autoassinada — ver seção [2.5](#25-componentes-de-software-do-istio). Em produção, é comum **plugar uma CA própria** (própria empresa ou uma CA intermediária de uma CA corporativa já existente), no modo chamado *plug-in CA*.

```bash
# 1. Gere (ou obtenha da sua CA corporativa) 4 arquivos: root-cert.pem, ca-cert.pem (intermediária),
#    ca-key.pem (chave privada da intermediária) e cert-chain.pem
mkdir -p certs && cd certs
make -f ../tools/certs/Makefile.selfsigned.mk root-ca
make -f ../tools/certs/Makefile.selfsigned.mk cluster1-cacerts

# 2. Crie o Secret "cacerts" ANTES de instalar/com o istiod já rodando
kubectl create namespace istio-system
kubectl create secret generic cacerts -n istio-system \
  --from-file=cluster1/ca-cert.pem \
  --from-file=cluster1/ca-key.pem \
  --from-file=cluster1/root-cert.pem \
  --from-file=cluster1/cert-chain.pem
```

O istiod monta esse Secret e passa a usá-lo como CA intermediária para assinar os certificados de workload (fluxo já detalhado na seção [2.5](#25-componentes-de-software-do-istio)); o certificado de cada workload tem validade padrão de **24 horas**, renovada automaticamente pelo `istio-agent` antes de expirar.

**Rotação do certificado raiz** (quando a própria CA raiz precisa trocar, não apenas o certificado de workload) é um processo **manual**: gera-se um `root-cert.pem` combinado (contendo a CA antiga **e** a nova), atualiza-se o Secret `cacerts` com a nova CA intermediária + a raiz combinada, e aguarda-se até um tempo de vida de certificado (até 24h) para que todos os workloads rotacionem naturalmente antes de remover a CA antiga do bundle combinado.

```bash
# Verificar se o cluster está em modo plug-in CA (Secret presente) e ver a raiz de confiança atual
kubectl get secret cacerts -n istio-system
kubectl get configmap istio-ca-root-cert -n istio-system -o yaml
```

Fonte: [Istio / Plug in CA Certificates](https://istio.io/v1.30/docs/tasks/security/cert-management/plugin-ca-cert/), [Istio / Security Concepts](https://istio.io/v1.30/docs/concepts/security/).

## 8.4. Sidecar: restringindo com quais endereços um workload se relaciona

O recurso `Sidecar` (`networking.istio.io`) controla **quais hosts** (`egress`) e **quais portas** (`ingress`) o proxy de um workload enxerga — por padrão, o Envoy de cada Pod carrega a configuração de **todos** os serviços da malha, o que não escala bem em meshes grandes e aumenta a superfície de tráfego "acidental" entre times.

```yaml
apiVersion: networking.istio.io/v1
kind: Sidecar
metadata:
  name: default
  namespace: myapp
spec:
  egress:
    - hosts:
        - "myapp/*"          # qualquer serviço no próprio namespace
        - "istio-system/*"   # necessário para os proxies falarem com istiod/telemetria
        - "shared-apis/api.exemplo.com"   # apenas ESSE host específico em outro namespace
  outboundTrafficPolicy:
    mode: REGISTRY_ONLY
```

Campos principais: `workloadSelector` (a quais Pods esse `Sidecar` se aplica — sem ele, aplica a todo o namespace), `ingress[].port`/`defaultEndpoint` (para onde o tráfego de entrada é encaminhado dentro do Pod, ex.: `127.0.0.1:8080`), `egress[].hosts` (lista no formato **`<namespace>/<dnsName>`**, aceitando `*` em ambas as partes), `egress[].captureMode` (`DEFAULT`, `IPTABLES` ou `NONE`) e `outboundTrafficPolicy.mode`. Um `Sidecar` sem `workloadSelector` colocado no **namespace raiz** funciona como o padrão de **toda a malha** — e existem três escopos possíveis (*mesh-wide*, *namespace-wide* e *workload-specific*, do mais genérico ao mais específico, análogos aos escopos de `PeerAuthentication`).

> O impacto prático de restringir o `egress` é grande: o livro *Istio in Action* mede, em um exemplo real, a configuração do Envoy de um workload caindo de **2 MB para 644 KB** (uma redução de ~70%) só por aplicar um `Sidecar` restringindo o egress ao essencial — o que também reduz o número de pushes que o istiod precisa fazer a cada mudança (seção [5.3](#53-métricas-do-control-plane-istiod-e-tempo-de-convergência)).

Fonte: [Istio / Sidecar reference](https://istio.io/v1.30/docs/reference/config/networking/sidecar/), POSTA & MALOKU, *Istio in Action* (Manning, 2022), capítulo 11.

## 8.5. AuthorizationPolicy: autorizando apenas aplicações específicas

`AuthorizationPolicy` (`security.istio.io`) define **quem** (identidade mTLS/JWT) pode fazer **o quê** (método, path, porta) contra um workload. Ações possíveis: `ALLOW` (padrão), `DENY`, `AUDIT` (loga sem bloquear) e `CUSTOM` (delega a um serviço externo de autorização). Ordem de avaliação quando várias políticas coexistem: **`CUSTOM` → `DENY` → `ALLOW`**; se existir ao menos uma política `ALLOW` para o workload e nenhuma casar com a requisição, ela é **negada por padrão**.

```yaml
# Só o service account "productpage" (no namespace myapp) pode chamar GET em /reviews* do serviço reviews
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: reviews-allow-productpage
  namespace: myapp
spec:
  selector:
    matchLabels:
      app: reviews
  action: ALLOW
  rules:
    - from:
        - source:
            principals: ["cluster.local/ns/myapp/sa/bookinfo-productpage"]
      to:
        - operation:
            methods: ["GET"]
            paths: ["/reviews*"]
```

```yaml
# Nega explicitamente qualquer tráfego vindo do namespace "dev" para "ratings"
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: ratings-deny-dev
  namespace: myapp
spec:
  selector:
    matchLabels:
      app: ratings
  action: DENY
  rules:
    - from:
        - source:
            namespaces: ["dev"]
```

Campos de `rules[].from[].source`: `principals` (identidade SPIFFE do peer mTLS, formato `<trust-domain>/ns/<namespace>/sa/<service-account>`), `namespaces`, `ipBlocks` (CIDR). Campos de `rules[].to[].operation`: `methods`, `paths`, `ports`, `hosts`. `rules[].when` permite condições adicionais usando atributos Istio (ex.: `request.auth.claims[iss]` — ver seção [8.6](#86-autorizando-por-jwt-eou-mtls-requestauthentication--authorizationpolicy)).

> ATENÇÃO!!! `principals`/`namespaces`/`requestPrincipals` só funcionam de forma confiável com **mTLS habilitado** (`PeerAuthentication` `STRICT` — Lab 2) — sem mTLS, o Envoy não tem como verificar criptograficamente a identidade de origem, e a documentação recomenda fortemente usar mTLS `STRICT` junto com essas regras para evitar bypass de política.

Fonte: [Istio / AuthorizationPolicy reference](https://istio.io/v1.30/docs/reference/config/security/authorization-policy/).

## 8.6. Autorizando por JWT e/ou mTLS (RequestAuthentication + AuthorizationPolicy)

A validação de **JWT** é feita pelo recurso `RequestAuthentication`, que **apenas valida** o token (assinatura, `issuer`, `audience`) e popula os atributos `request.auth.*`/`source.principal` do Envoy — ele **não bloqueia** requisições sem token por conta própria; o bloqueio é sempre feito combinando com uma `AuthorizationPolicy`.

```yaml
# 1. Valida qualquer JWT emitido por "testing@secure.istio.io" chegando no ingress gateway
apiVersion: security.istio.io/v1
kind: RequestAuthentication
metadata:
  name: jwt-example
  namespace: istio-system
spec:
  selector:
    matchLabels:
      istio: ingressgateway
  jwtRules:
    - issuer: "testing@secure.istio.io"
      jwksUri: "https://raw.githubusercontent.com/istio/istio/release-1.30/security/tools/jwt/samples/jwks.json"
      audiences: ["api-service"]
      forwardOriginalToken: true
```

```yaml
# 2. EXIGE que a requisição tenha um JWT válido (nega quem não tiver um "request principal")
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: require-jwt
  namespace: istio-system
spec:
  selector:
    matchLabels:
      istio: ingressgateway
  action: DENY
  rules:
    - from:
        - source:
            notRequestPrincipals: ["*"]
```

Campos de `jwtRules[]`: `issuer` ("identifies the issuer that issued the JWT" — valida o claim `iss`), `jwksUri`/`jwks` (onde obter as chaves públicas — URL ou o JWKS embutido), `audiences` (valida o claim `aud`), `fromHeaders`/`fromParams` (onde procurar o token na requisição, com suporte a prefixo tipo `"Bearer "`), `outputPayloadToHeader` (encaminha o payload decodificado do JWT para a aplicação em um header) e `forwardOriginalToken`.

**Combinando JWT *e* mTLS na mesma regra** (o cenário pedido: só aceitar um serviço se ele tiver identidade mTLS **e** um JWT válido de um usuário/serviço específico) — campos dentro do **mesmo** bloco `source` são combinados com **E** (AND):

```yaml
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: reviews-mtls-and-jwt
  namespace: myapp
spec:
  selector:
    matchLabels:
      app: reviews
  action: ALLOW
  rules:
    - from:
        - source:
            principals: ["cluster.local/ns/myapp/sa/bookinfo-productpage"]   # exige mTLS desse workload
            requestPrincipals: ["testing@secure.istio.io/*"]                 # E exige JWT válido desse issuer
      to:
        - operation:
            methods: ["GET"]
```

Fonte: [Istio / Authentication Policy (JWT)](https://istio.io/v1.30/docs/tasks/security/authentication/authn-policy/), [Istio / RequestAuthentication reference](https://istio.io/v1.30/docs/reference/config/security/request_authentication/), [Istio / AuthorizationPolicy reference](https://istio.io/v1.30/docs/reference/config/security/authorization-policy/).

---

# 9. Exercícios de troubleshooting

Os laboratórios a seguir usam o cluster **kind** e o **Istio** já instalados conforme [day1.md](day1.md) (modo `ambient`) — o Lab 2 usa o modo `sidecar` (ver [OLD/day1-sidecar.md](OLD/day1-sidecar.md)) porque a falha que ele demonstra (mTLS `STRICT` x `PERMISSIVE`) é conceito de `PeerAuthentication`, válido nos dois modos, mas mais didático de observar com o Envoy sidecar.

Cada laboratório segue o formato: **Objetivo → Provocar a falha → Diagnosticar → Corrigir → Referência**.

## 9.0. Preparando o ambiente dos laboratórios

```bash
export MY_NAMESPACE='myapp'
export ISTIO_BASE_URL="https://raw.githubusercontent.com/istio/istio/release-1.30/samples/"
export ISTIO_BOOKINFO_URL="${ISTIO_BASE_URL}bookinfo/"

kubectl create namespace $MY_NAMESPACE
kubectl label namespace $MY_NAMESPACE istio.io/dataplane-mode=ambient

kubectl -n $MY_NAMESPACE apply -f "${ISTIO_BOOKINFO_URL}platform/kube/bookinfo.yaml"
kubectl -n $MY_NAMESPACE apply -f "${ISTIO_BOOKINFO_URL}platform/kube/bookinfo-versions.yaml"
kubectl -n $MY_NAMESPACE rollout status deployment --timeout=300s
```

## Lab 1 - Pod fora da mesh (ambient)

**Objetivo:** entender como confirmar se um workload está de fato participando da malha e diagnosticar quando não está.

**Provocar a falha** (remova o namespace do ambient mode):

```bash
kubectl label namespace $MY_NAMESPACE istio.io/dataplane-mode-
kubectl -n $MY_NAMESPACE rollout restart deployment
```

**Diagnosticar:**

```bash
# O workload não aparecerá na lista de workloads conhecidos pelo ztunnel
istioctl ztunnel-config workload | grep $MY_NAMESPACE

# proxy-status também não vai listar o Pod (ele não existe como "proxy" da malha)
istioctl proxy-status | grep $MY_NAMESPACE

# Confirme o label do namespace
kubectl get namespace $MY_NAMESPACE --show-labels
```

**Causa raiz:** o label `istio.io/dataplane-mode=ambient` é o que instrui o `istio-cni` a programar o redirecionamento de tráfego do Pod para o `ztunnel`; sem ele, o tráfego do Pod nunca passa pela malha (sem mTLS, sem políticas, sem telemetria da mesh).

**Corrigir:**

```bash
kubectl label namespace $MY_NAMESPACE istio.io/dataplane-mode=ambient
kubectl -n $MY_NAMESPACE rollout restart deployment
istioctl ztunnel-config workload | grep $MY_NAMESPACE
```

**Referência:** [Istio / Ambient — Secure and visualize your mesh](https://istio.io/v1.30/docs/ambient/getting-started/secure-and-visualize/), [Istio / Troubleshoot ztunnel](https://istio.io/latest/docs/ambient/usage/troubleshoot-ztunnel/).

## Lab 2 - mTLS STRICT quebrando tráfego (sidecar)

**Objetivo:** observar o efeito de uma `PeerAuthentication` `STRICT` quando existe tráfego em texto plano chegando ao workload.

> Este lab pressupõe instalação em modo **sidecar** (ver [OLD/day1-sidecar.md](OLD/day1-sidecar.md)) e sidecar injection habilitada no namespace (`kubectl label namespace $MY_NAMESPACE istio-injection=enabled`).

**Provocar a falha:**

```bash
cat <<EOF | kubectl apply -n $MY_NAMESPACE -f -
apiVersion: security.istio.io/v1
kind: PeerAuthentication
metadata:
  name: default
spec:
  mtls:
    mode: STRICT
EOF

# Gera tráfego em texto plano diretamente para o Pod, sem passar por outro sidecar da mesh
kubectl run --rm -it plaintext-client --image=curlimages/curl --restart=Never -n $MY_NAMESPACE -- \
  curl -sS -o /dev/null -w "HTTP %{http_code}\n" http://productpage:9080/productpage
```

**Diagnosticar:**

```bash
# Verifique se a requisição falhou (connection reset / EOF, pois o Envoy STRICT rejeita texto plano)
kubectl -n $MY_NAMESPACE logs deploy/productpage-v1 -c istio-proxy | tail -n 20

# Confirme o modo mTLS efetivo
istioctl x describe pod "$(kubectl -n $MY_NAMESPACE get pod -l app=productpage -o jsonpath='{.items[0].metadata.name}')" -n $MY_NAMESPACE
```

**Causa raiz:** com `STRICT`, o sidecar só aceita conexões mTLS; um client sem sidecar (como o `curl` acima, fora da mesh) não apresenta certificado, então a conexão é recusada.

**Corrigir** (duas opções válidas a depender do objetivo):

```bash
# Opção 1: usar PERMISSIVE (aceita texto plano E mTLS - bom durante migração)
kubectl -n $MY_NAMESPACE patch peerauthentication default --type merge -p '{"spec":{"mtls":{"mode":"PERMISSIVE"}}}'

# Opção 2: remover o client de fora da mesh e testar com um Pod que tenha sidecar
kubectl -n $MY_NAMESPACE delete peerauthentication default
```

**Referência:** [Istio / Mutual TLS Migration](https://istio.io/v1.30/docs/tasks/security/authentication/mtls-migration/), [Istio / PeerAuthentication reference](https://istio.io/v1.30/docs/reference/config/security/peer_authentication/).

## Lab 3 - VirtualService/DestinationRule com host errado (503 / no healthy upstream)

**Objetivo:** reproduzir e diagnosticar o erro clássico "no healthy upstream"/`503 UC` causado por uma referência de host incorreta.

**Provocar a falha:**

```bash
cat <<EOF | kubectl apply -n $MY_NAMESPACE -f -
apiVersion: networking.istio.io/v1
kind: DestinationRule
metadata:
  name: reviews-quebrado
spec:
  host: reviews.nao-existe.svc.cluster.local
  subsets:
  - name: v1
    labels:
      version: v1
---
apiVersion: networking.istio.io/v1
kind: VirtualService
metadata:
  name: reviews-quebrado
spec:
  hosts:
  - reviews
  http:
  - route:
    - destination:
        host: reviews.nao-existe.svc.cluster.local
        subset: v1
EOF
```

**Diagnosticar:**

```bash
# istioctl analyze deve sinalizar o host inválido / referenced host not found
istioctl analyze -n $MY_NAMESPACE

# Verifique a existência (ou não) do cluster correspondente no Envoy do productpage
istioctl proxy-config cluster "$(kubectl -n $MY_NAMESPACE get pod -l app=productpage -o jsonpath='{.items[0].metadata.name}')" -n $MY_NAMESPACE | grep reviews

# Gere tráfego e observe o response flag UC (upstream connect error) ou NR (no route) no access log
kubectl -n $MY_NAMESPACE logs deploy/productpage-v1 -c istio-proxy -f
```

**Causa raiz:** o `host` do `DestinationRule`/`VirtualService` aponta para um serviço Kubernetes inexistente (`reviews.nao-existe...`), então o Envoy não tem endpoints (EDS vazio) para esse cluster.

**Corrigir:**

```bash
kubectl -n $MY_NAMESPACE delete virtualservice reviews-quebrado
kubectl -n $MY_NAMESPACE delete destinationrule reviews-quebrado
```

**Referência:** [Istio / Common Problems — 503 errors after setting a DestinationRule](https://istio.io/v1.30/docs/ops/common-problems/network-issues/), [Istio / istioctl analyze](https://istio.io/v1.30/docs/ops/diagnostic-tools/istioctl-analyze/).

## Lab 4 - istiod fora de sincronia (STALE / NOT SENT)

**Objetivo:** identificar proxies que não estão sincronizados com o istiod (`STALE`/`NOT SENT`) e investigar a causa.

**Provocar a falha** (simulando indisponibilidade momentânea do control plane):

```bash
kubectl -n istio-system rollout restart deployment/istiod
```

**Diagnosticar (imediatamente após o restart):**

```bash
# Observe proxies em STALE ou NOT SENT enquanto o novo istiod inicializa
istioctl proxy-status

# Detalhe a diferença entre a config atual do Envoy e a que o istiod enviaria
istioctl proxy-status <pod>.$MY_NAMESPACE

# Confirme conectividade à porta xDS (15012) a partir de um Pod da malha
kubectl -n $MY_NAMESPACE exec deploy/productpage-v1 -c istio-proxy -- \
  curl -sk https://istiod.istio-system:15012 -o /dev/null -w "%{http_code}\n"
```

**Causa raiz mais comum em ambiente real:** `NetworkPolicy` bloqueando a porta 15012, falta de recursos (CPU/memória) no istiod, ou pico de reconexão simultânea de muitos proxies após um restart do control plane.

**Corrigir:**

```bash
kubectl -n istio-system rollout status deployment/istiod --timeout=180s
istioctl proxy-status
```

**Referência:** [Istio / Debugging Envoy and Istiod — proxy-status](https://istio.io/v1.30/docs/ops/diagnostic-tools/proxy-cmd/).

## Lab 5 - Inspecionando o Envoy diretamente pela Admin API

**Objetivo:** praticar o uso da Admin API do Envoy (porta 15000) para confirmar visualmente os conceitos da seção [3.2](#32-arquitetura-interna-listener-subsystem-x-cluster-subsystem) (listeners, clusters, endpoints).

```bash
export PP_POD=$(kubectl -n $MY_NAMESPACE get pod -l app=productpage -o jsonpath='{.items[0].metadata.name}')
kubectl -n $MY_NAMESPACE port-forward pod/$PP_POD 15000:15000 &

# Liste os listeners abertos (ex.: inbound 9080, outbound para virtual services, etc.)
curl -s http://localhost:15000/listeners

# Encontre o cluster outbound do serviço "reviews" e confira seus endpoints/saúde
curl -s http://localhost:15000/clusters | grep reviews

# Veja o config_dump completo do listener outbound HTTP (porta 0.0.0.0_9080, por exemplo)
curl -s "http://localhost:15000/config_dump?resource=dynamic_listeners" | jq .

kill %1   # encerra o port-forward em background
```

**Referência:** [Envoy / Administration interface](https://www.envoyproxy.io/docs/envoy/latest/operations/admin).

## Lab 6 - Política L7 ignorada em ambient mode sem waypoint

**Objetivo:** demonstrar que políticas L7 (ex.: restringir por path HTTP) **não** têm efeito em ambient mode sem um **waypoint**, pois o `ztunnel` só enforça L4.

**Provocar o cenário:**

```bash
cat <<EOF | kubectl apply -n $MY_NAMESPACE -f -
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: bloqueia-reviews-path
spec:
  selector:
    matchLabels:
      app: reviews
  action: DENY
  rules:
  - to:
    - operation:
        paths: ["/wpcatalog*"]
EOF

# Tráfego normal ao productpage continua funcionando mesmo esperando bloqueio de path,
# porque não há waypoint para reviews neste namespace ainda
kubectl -n $MY_NAMESPACE exec deploy/productpage-v1 -- \
  curl -sS -o /dev/null -w "HTTP %{http_code}\n" http://reviews:9080/wpcatalog
```

**Diagnosticar:**

```bash
# Confirme que não existe waypoint associado ao serviço reviews
istioctl ztunnel-config service | grep reviews

istioctl analyze -n $MY_NAMESPACE   # pode alertar sobre policy L7 sem waypoint, dependendo da versão
```

**Causa raiz:** `AuthorizationPolicy` com condição de **path** é uma regra **L7**; o `ztunnel` (L3/L4) não consegue avaliá-la, então é necessário um `waypoint` na frente do workload alvo.

**Corrigir:**

```bash
istioctl waypoint apply -n $MY_NAMESPACE --enroll-namespace
istioctl ztunnel-config service | grep reviews   # agora deve mostrar o waypoint associado

kubectl -n $MY_NAMESPACE exec deploy/productpage-v1 -- \
  curl -sS -o /dev/null -w "HTTP %{http_code}\n" http://reviews:9080/wpcatalog   # esperado: 403
```

**Limpeza:**

```bash
kubectl -n $MY_NAMESPACE delete authorizationpolicy bloqueia-reviews-path
istioctl waypoint delete -n $MY_NAMESPACE --all
```

**Referência:** [Istio / Ambient — L4 vs L7 policies](https://istio.io/latest/docs/ambient/usage/), [Istio / Troubleshoot waypoints](https://istio.io/latest/docs/ambient/usage/troubleshoot-waypoint/).

## Lab 7 - Diagnosticando o ztunnel (ambient)

**Objetivo:** usar os comandos de troubleshooting oficiais do `ztunnel` para confirmar certificados e tráfego mTLS.

```bash
export ZTUNNEL=$(kubectl get pods -n istio-system -l app=ztunnel -o jsonpath='{.items[0].metadata.name}')

# Certificados x509 carregados (deve haver um por workload identity no nó)
istioctl ztunnel-config certificate "$ZTUNNEL".istio-system

# Workloads conhecidos pelo ztunnel deste nó
istioctl ztunnel-config workload

# Logs em tempo real filtrando conexões inbound/outbound (gere tráfego do Bookinfo em paralelo)
kubectl -n istio-system logs -f "$ZTUNNEL" | grep -E "inbound|outbound"
```

**Referência:** [Istio / Troubleshoot connectivity issues with ztunnel](https://istio.io/latest/docs/ambient/usage/troubleshoot-ztunnel/).

## Lab 8 - `istioctl analyze` encontrando erros de configuração

**Objetivo:** usar `istioctl analyze` (estático e contra o cluster) para achar erros antes de aplicá-los.

```bash
# Analisa arquivos locais sem tocar no cluster
istioctl analyze --use-kube=false files/bookinfo-gateway.yaml

# Analisa tudo o que já está aplicado no namespace
istioctl analyze -n $MY_NAMESPACE

# Analisa o cluster inteiro
istioctl analyze --all-namespaces
```

Tipos de achado, conforme a documentação: **erros** (ex.: `Referenced selector not found`, quando um `Gateway` referencia um selector inexistente), **avisos** (ex.: host do `DestinationRule` sem correspondência com nenhum serviço na mesh) e **informações** (ex.: namespace sem sidecar injection habilitada).

**Referência:** [Istio / istioctl analyze](https://istio.io/v1.30/docs/ops/diagnostic-tools/istioctl-analyze/).

## Lab 9 - Fault injection e timeout (retomando o Day 2)

Os exercícios de **fault injection** (delay de 7s entre `reviews:v2` e `ratings`) e **request timeout** já estão documentados com comandos completos em [day2.md#fault-injection](day2.md#fault-injection). Ao executá-los, use os comandos das seções [4.2](#42-istioctl-comandos-essenciais) e [4.3](#43-admin-api-do-envoy-porta-15000) deste guia para:

- Confirmar, via `istioctl proxy-config route`, que o `timeout`/`fault.delay` está presente na rota ativa do Envoy (não apenas no YAML aplicado).
- Observar, via `/stats` da Admin API do Envoy, os contadores `upstream_rq_timeout` e `upstream_rq_pending_overflow` durante o teste.

**Referência:** [Istio / Fault Injection](https://istio.io/v1.30/docs/tasks/traffic-management/fault-injection/), [day2.md](day2.md).

## Lab 10 - Circuit breaking e overflow detectados via métricas

**Objetivo:** provocar um `circuit breaker` de conexões no `DestinationRule` e observar o efeito diretamente nas métricas do Envoy (seção [5.2](#52-estatísticas-internas-do-envoy)), sem depender apenas do código de resposta HTTP.

**Provocar a falha** (limite de 1 conexão/1 requisição pendente para `reviews`, seguido de carga concorrente):

```bash
cat <<EOF | kubectl apply -n $MY_NAMESPACE -f -
apiVersion: networking.istio.io/v1
kind: DestinationRule
metadata:
  name: reviews-circuit-breaker
spec:
  host: reviews
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 1
      http:
        http1MaxPendingRequests: 1
        maxRequestsPerConnection: 1
    outlierDetection:
      consecutive5xxErrors: 1
      interval: 5s
      baseEjectionTime: 30s
EOF

# Gera carga concorrente (20 requisições simultâneas) para estourar o circuit breaker
kubectl -n $MY_NAMESPACE exec deploy/productpage-v1 -- sh -c \
  'for i in $(seq 1 20); do curl -sS -o /dev/null -w "%{http_code}\n" http://reviews:9080/reviews/0 & done; wait'
```

**Diagnosticar (via Admin API do Envoy, seção [4.3](#43-admin-api-do-envoy-porta-15000)/[4.5](#45-comandos-avançados-de-troubleshooting-do-envoy)):**

```bash
kubectl -n $MY_NAMESPACE port-forward deploy/productpage-v1 15000:15000 &

# upstream_rq_pending_overflow > 0 confirma requisições rejeitadas pelo circuit breaker
curl -s "http://localhost:15000/stats?filter=reviews.*(overflow|pending)"

# Estado do circuit breaker (cx_open / rq_open = 1 quando aberto)
curl -s "http://localhost:15000/stats?filter=reviews.*circuit_breakers"

kill %1
```

Se o Prometheus dos addons estiver de pé (ver [day1.md](day1.md)), a mesma condição aparece na consulta:

```promql
sum(rate(istio_requests_total{destination_service_name="reviews", response_flags=~"UO"}[1m]))
```

**Causa raiz:** `maxConnections: 1` e `http1MaxPendingRequests: 1` são limites artificialmente baixos; qualquer concorrência acima disso é imediatamente rejeitada pelo circuit breaker do Envoy (flag `UO`, métrica `upstream_rq_pending_overflow`/`upstream_cx_overflow`).

**Corrigir:**

```bash
kubectl -n $MY_NAMESPACE delete destinationrule reviews-circuit-breaker
```

**Referência:** [Envoy / Circuit breaking](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/circuit_breaking), [Istio / DestinationRule reference](https://istio.io/v1.30/docs/reference/config/networking/destination-rule/).

## Lab 11 - Medindo o tempo de convergência do istiod

**Objetivo:** observar a métrica `pilot_proxy_convergence_time` (seção [5.3](#53-métricas-do-control-plane-istiod-e-tempo-de-convergência)) refletindo, na prática, o tempo entre aplicar uma configuração e ela chegar a todos os proxies.

```bash
# Consulta direta às métricas cruas do istiod (sem precisar do Prometheus dos addons)
export ISTIOD=$(kubectl get pods -n istio-system -l app=istiod -o jsonpath='{.items[0].metadata.name}')
kubectl -n istio-system port-forward "$ISTIOD" 15014:15014 &

curl -s http://localhost:15014/metrics | grep -A3 '^pilot_proxy_convergence_time'

kill %1
```

```bash
# Aplique uma mudança de configuração (ex.: o VirtualService do Lab 3, já removido) e,
# em paralelo, observe istioctl proxy-status até TODOS os proxies voltarem a SYNCED
kubectl -n $MY_NAMESPACE apply -f "${ISTIO_BOOKINFO_URL}networking/virtual-service-all-v1.yaml"
watch -n1 istioctl proxy-status

kubectl -n $MY_NAMESPACE delete -f "${ISTIO_BOOKINFO_URL}networking/virtual-service-all-v1.yaml"
```

Se estiver usando o Prometheus dos addons, a consulta PromQL da seção [5.3](#53-métricas-do-control-plane-istiod-e-tempo-de-convergência) mostra a mesma informação de forma agregada e histórica, o que é mais útil do que ler o texto bruto da métrica em produção.

**Referência:** [Istio / Performance and Scalability](https://istio.io/v1.30/docs/ops/deployment/performance-and-scalability/), [Istio / Debugging Envoy and Istiod](https://istio.io/v1.30/docs/ops/diagnostic-tools/proxy-cmd/).

## Lab 12 - Compressão gzip via EnvoyFilter não aplicada

**Objetivo:** praticar o troubleshooting de um `EnvoyFilter` (seção [7](#7-envoyfilter)) que "não faz nada" por causa de um `match` incorreto — o erro mais comum com esse recurso.

**Provocar a falha** (contexto errado: a patch mira o tráfego **outbound** do productpage, mas o que precisamos comprimir é a resposta que ele mesmo gera, ou seja, o tráfego **inbound**):

```bash
cat <<EOF | kubectl apply -n $MY_NAMESPACE -f -
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: enable-gzip-compression
spec:
  workloadSelector:
    labels:
      app: productpage
  configPatches:
  - applyTo: HTTP_FILTER
    match:
      context: SIDECAR_OUTBOUND
      listener:
        filterChain:
          filter:
            name: "envoy.filters.network.http_connection_manager"
            subFilter:
              name: "envoy.filters.http.router"
    patch:
      operation: INSERT_BEFORE
      value:
        name: envoy.filters.http.compressor
        typed_config:
          "@type": type.googleapis.com/envoy.extensions.filters.http.compressor.v3.Compressor
          compressor_library:
            name: text_optimized
            typed_config:
              "@type": type.googleapis.com/envoy.extensions.compression.gzip.compressor.v3.Gzip
EOF

# Teste: a resposta não vem comprimida, mesmo pedindo gzip
kubectl -n $MY_NAMESPACE exec deploy/productpage-v1 -c istio-proxy -- \
  curl -sS -o /dev/null -D - -H "Accept-Encoding: gzip" http://localhost:9080/productpage | grep -i content-encoding
```

**Diagnosticar:**

```bash
# Procure o filtro na configuração REAL do listener inbound (porta 9080) do productpage
istioctl proxy-config listener deploy/productpage-v1.$MY_NAMESPACE --port 9080 -o json | grep -i compressor
# Não encontra nada -> confirma que a patch não foi aplicada onde a requisição de fato passa

istioctl analyze -n $MY_NAMESPACE
```

**Causa raiz:** `context: SIDECAR_OUTBOUND` faz a patch mirar o listener de **saída** do productpage (quando ele chama outros serviços), mas a resposta HTTP que queremos comprimir é servida no listener de **entrada** (`SIDECAR_INBOUND`, porta 9080) — exatamente o cenário descrito na seção [7.6](#76-exemplo-prático-habilitando-compressão-gzip-via-envoyfilter).

**Corrigir:**

```bash
kubectl -n $MY_NAMESPACE patch envoyfilter enable-gzip-compression --type=json \
  -p '[{"op":"replace","path":"/spec/configPatches/0/match/context","value":"SIDECAR_INBOUND"}]'

istioctl proxy-config listener deploy/productpage-v1.$MY_NAMESPACE --port 9080 -o json | grep -i compressor
```

**Limpeza:**

```bash
kubectl -n $MY_NAMESPACE delete envoyfilter enable-gzip-compression
```

**Referência:** [Istio / EnvoyFilter reference](https://istio.io/v1.30/docs/reference/config/networking/envoy-filter/), seção [7.7](#77-comandos-de-troubleshooting-de-envoyfilter) deste guia.

## Lab 13 - AuthorizationPolicy negando (ou liberando) por engano

**Objetivo:** observar o efeito de uma `AuthorizationPolicy` mal escrita — um erro comum é esquecer que, assim que existe **qualquer** política `ALLOW` para um workload, **tudo que não casar** com nenhuma regra passa a ser **negado por padrão** (seção [8.5](#85-authorizationpolicy-autorizando-apenas-aplicações-específicas)).

**Provocar a falha:**

```bash
cat <<EOF | kubectl apply -n $MY_NAMESPACE -f -
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: reviews-allow-only-productpage
spec:
  selector:
    matchLabels:
      app: reviews
  action: ALLOW
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/$MY_NAMESPACE/sa/bookinfo-productpage"]
EOF

# productpage continua funcionando, mas QUALQUER outro caller (ex.: um curl manual) passa a ser negado (RBAC: access denied)
kubectl -n $MY_NAMESPACE run --rm -it curl-test --image=curlimages/curl --restart=Never -- \
  curl -sS -o /dev/null -w "HTTP %{http_code}\n" http://reviews.$MY_NAMESPACE.svc.cluster.local:9080/reviews/0
```

**Diagnosticar:**

```bash
# O access log do Envoy de "reviews" mostra RBAC: access denied com response_code 403
kubectl -n $MY_NAMESPACE logs deploy/reviews-v1 -c istio-proxy | grep -i denied

# istioctl x describe explica quais policies afetam o Pod alvo
istioctl x describe pod "$(kubectl -n $MY_NAMESPACE get pod -l app=reviews -o jsonpath='{.items[0].metadata.name}')" -n $MY_NAMESPACE
```

**Causa raiz:** essa `AuthorizationPolicy` é o comportamento **esperado e correto** — qualquer origem que não seja o service account `bookinfo-productpage` é implicitamente negada, pois existe uma política `ALLOW` para `reviews` e nenhuma regra casa com o `curl-test`. É um erro comum interpretar esse 403 como "bug" quando na verdade é a política funcionando como projetada.

**Corrigir (se a intenção era permitir mais origens):**

```bash
kubectl -n $MY_NAMESPACE delete authorizationpolicy reviews-allow-only-productpage
```

**Referência:** [Istio / AuthorizationPolicy reference](https://istio.io/v1.30/docs/reference/config/security/authorization-policy/), seção [8.5](#85-authorizationpolicy-autorizando-apenas-aplicações-específicas) deste guia.

## Lab 14 - JWT inválido/ausente sendo rejeitado (ou não)

**Objetivo:** demonstrar a diferença entre **validar** um JWT (`RequestAuthentication`) e **exigir** um JWT válido (`AuthorizationPolicy` com `notRequestPrincipals`) — seção [8.6](#86-autorizando-por-jwt-eou-mtls-requestauthentication--authorizationpolicy).

**Provocar o cenário (só `RequestAuthentication`, sem `AuthorizationPolicy`):**

```bash
cat <<EOF | kubectl apply -n $MY_NAMESPACE -f -
apiVersion: security.istio.io/v1
kind: RequestAuthentication
metadata:
  name: jwt-example
spec:
  selector:
    matchLabels:
      app: productpage
  jwtRules:
  - issuer: "testing@secure.istio.io"
    jwksUri: "https://raw.githubusercontent.com/istio/istio/release-1.30/security/tools/jwt/samples/jwks.json"
EOF

# Requisição SEM token: ainda funciona (RequestAuthentication sozinho não bloqueia nada)
kubectl -n $MY_NAMESPACE exec deploy/productpage-v1 -c istio-proxy -- \
  curl -sS -o /dev/null -w "sem token: HTTP %{http_code}\n" http://localhost:9080/productpage

# Requisição com token INVÁLIDO: agora sim é rejeitada com 401 (a assinatura não bate)
kubectl -n $MY_NAMESPACE exec deploy/productpage-v1 -c istio-proxy -- \
  curl -sS -o /dev/null -w "token invalido: HTTP %{http_code}\n" \
  -H "Authorization: Bearer token-invalido-de-propósito" http://localhost:9080/productpage
```

**Diagnosticar:** o comportamento acima é o **esperado**: `RequestAuthentication` valida um token **se ele existir**, mas não torna o token obrigatório.

**Corrigir (tornando o JWT obrigatório):**

```bash
cat <<EOF | kubectl apply -n $MY_NAMESPACE -f -
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: require-jwt-productpage
spec:
  selector:
    matchLabels:
      app: productpage
  action: DENY
  rules:
  - from:
    - source:
        notRequestPrincipals: ["*"]
EOF

# Agora, sem token, a requisição é negada com 403 (não mais 200)
kubectl -n $MY_NAMESPACE exec deploy/productpage-v1 -c istio-proxy -- \
  curl -sS -o /dev/null -w "sem token apos fix: HTTP %{http_code}\n" http://localhost:9080/productpage
```

**Limpeza:**

```bash
kubectl -n $MY_NAMESPACE delete authorizationpolicy require-jwt-productpage
kubectl -n $MY_NAMESPACE delete requestauthentication jwt-example
```

**Referência:** [Istio / Authentication Policy (JWT)](https://istio.io/v1.30/docs/tasks/security/authentication/authn-policy/), seção [8.6](#86-autorizando-por-jwt-eou-mtls-requestauthentication--authorizationpolicy) deste guia.

---

# 10. Referências

**Istio (documentação oficial, versão 1.30):**

- [Istio / What is Istio?](https://istio.io/v1.30/docs/overview/what-is-istio/)
- [Istio / Dataplane modes](https://istio.io/v1.30/docs/overview/dataplane-modes/)
- [Istio / Architecture (deployment)](https://istio.io/v1.30/docs/ops/deployment/architecture/)
- [Istio / Ambient Overview](https://istio.io/latest/docs/ambient/overview/)
- [Istio / Ambient Architecture — Traffic Redirection](https://istio.io/v1.30/docs/ambient/architecture/traffic-redirection/)
- [Istio / Glossary](https://istio.io/v1.30/docs/reference/glossary/)
- [Istio / Application Requirements (portas)](https://istio.io/v1.30/docs/ops/deployment/application-requirements/)
- [Istio / Debug Endpoints](https://istio.io/latest/docs/ops/integrations/integration-guide/debug-endpoints/)
- [Istio / Debugging Envoy and Istiod (istioctl proxy-status / proxy-config)](https://istio.io/v1.30/docs/ops/diagnostic-tools/proxy-cmd/)
- [Istio / istioctl analyze](https://istio.io/v1.30/docs/ops/diagnostic-tools/istioctl-analyze/)
- [Istio / istioctl command reference](https://istio.io/v1.30/docs/reference/commands/istioctl/)
- [Istio / Common Problems — Network Issues](https://istio.io/v1.30/docs/ops/common-problems/network-issues/)
- [Istio / Mutual TLS Migration](https://istio.io/v1.30/docs/tasks/security/authentication/mtls-migration/)
- [Istio / PeerAuthentication reference](https://istio.io/v1.30/docs/reference/config/security/peer_authentication/)
- [Istio / Fault Injection](https://istio.io/v1.30/docs/tasks/traffic-management/fault-injection/)
- [Istio / Troubleshoot connectivity issues with ztunnel](https://istio.io/latest/docs/ambient/usage/troubleshoot-ztunnel/)
- [Istio / Troubleshoot issues with waypoints](https://istio.io/latest/docs/ambient/usage/troubleshoot-waypoint/)
- [Istio / Supported Releases](https://istio.io/latest/docs/releases/supported-releases/)
- [Istio / Security Concepts](https://istio.io/v1.30/docs/concepts/security/)
- [Istio / Standard Metrics Reference](https://istio.io/v1.30/docs/reference/config/metrics/)
- [Istio / Envoy Statistics](https://istio.io/latest/docs/ops/configuration/telemetry/envoy-stats/)
- [Istio / Performance and Scalability](https://istio.io/v1.30/docs/ops/deployment/performance-and-scalability/)
- [Istio / Grafana](https://istio.io/v1.30/docs/ops/integrations/grafana/)
- [Istio / Visualizing Metrics with Grafana](https://istio.io/latest/docs/tasks/observability/metrics/using-istio-dashboard/)
- [Istio / DestinationRule reference](https://istio.io/v1.30/docs/reference/config/networking/destination-rule/)
- [Istio / Component logging](https://istio.io/v1.30/docs/ops/diagnostic-tools/component-logging/)
- [Istio / Envoy Access Logs](https://istio.io/v1.30/docs/tasks/observability/logs/access-log/)
- [Istio / EnvoyFilter reference](https://istio.io/v1.30/docs/reference/config/networking/envoy-filter/)
- [Istio / Telemetry API reference](https://istio.io/v1.30/docs/reference/config/telemetry/)
- [Istio / MeshConfig reference](https://istio.io/v1.30/docs/reference/config/istio.mesh.v1alpha1/)
- [Istio / Plug in CA Certificates](https://istio.io/v1.30/docs/tasks/security/cert-management/plugin-ca-cert/)
- [Istio / Sidecar reference](https://istio.io/v1.30/docs/reference/config/networking/sidecar/)
- [Istio / AuthorizationPolicy reference](https://istio.io/v1.30/docs/reference/config/security/authorization-policy/)
- [Istio / Authentication Policy (JWT)](https://istio.io/v1.30/docs/tasks/security/authentication/authn-policy/)
- [Istio / RequestAuthentication reference](https://istio.io/v1.30/docs/reference/config/security/request_authentication/)
- [Istio / Deployment Models](https://istio.io/v1.30/docs/ops/deployment/deployment-models/)
- [howardjohn's blog — Everything you need to know about Istio installation](https://blog.howardjohn.info/posts/istio-install/)
- [github.com/istio/istio — Releases](https://github.com/istio/istio/releases/)

**Envoy (documentação oficial):**

- [Envoy / Introduction](https://www.envoyproxy.io/docs/envoy/latest/intro/intro)
- [Envoy / Architecture Overview](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/arch_overview)
- [Envoy / Life of a Request](https://www.envoyproxy.io/docs/envoy/latest/intro/life_of_a_request)
- [Envoy / xDS REST and gRPC protocol](https://www.envoyproxy.io/docs/envoy/latest/api-docs/xds_protocol)
- [Envoy / Administration interface](https://www.envoyproxy.io/docs/envoy/latest/operations/admin)
- [Envoy / Threading Model](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/intro/threading_model)
- [Envoy / Hot Restart](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/hot_restart)
- [Envoy / Circuit breaking](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/circuit_breaking)
- [Envoy / Cluster statistics (GitHub)](https://github.com/envoyproxy/envoy/blob/main/docs/root/configuration/upstream/cluster_manager/cluster_stats.rst)
- [Envoy / Substitution Formatter — Response Flags](https://www.envoyproxy.io/docs/envoy/latest/configuration/advanced/substitution_formatter)
- [Envoy / Stats proto (prefixo `envoy_` e tags no formato Prometheus)](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/metrics/v3/stats.proto)
- [Envoy / Access logging](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage)
- [Envoy / Compressor filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/compressor_filter)
- [Envoy / RBAC filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/rbac_filter)
- [Envoy / Response Code Details](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/response_code_details)

**Prometheus, Grafana e Datadog:**

- [Datadog / Istio Integration](https://docs.datadoghq.com/integrations/istio/)
- [Datadog — How to Monitor Istiod](https://www.datadoghq.com/blog/how-to-monitor-istiod/)
- [Datadog / Envoy Integration](https://github.com/DataDog/integrations-core/tree/master/envoy#istio)

**Livros e outras referências complementares:**

- POSTA, Christian E.; MALOKU, Rinor. *Istio in Action*. Manning Publications, 2022, ISBN 9781617295829 (capítulos 3 — "Istio's data plane: The Envoy proxy", 11 — "Performance-tuning the control plane", e 12 — "Scaling Istio in your organization", usados nas seções [2.6](#26-topologias-de-implantação-da-mesh-multicluster-multi-network-multi-mesh), [5.3](#53-métricas-do-control-plane-istiod-e-tempo-de-convergência) e [8.4](#84-sidecar-restringindo-com-quais-endereços-um-workload-se-relaciona) deste guia)
- [Christian Posta's blog — posts sobre Istio](https://blog.christianposta.com/tags/istio/)
- [Christian Posta — Understanding Istio Telemetry v2](https://blog.christianposta.com/understanding-istio-telemetry-v2/)
- [github.com/Apress/Practical-Istio](https://github.com/Apress/Practical-Istio) — código-fonte do livro *Practical Istio* (Josyula, Singh e Mehta; Apress, 2025)
- [AWS / Transit Gateway](https://aws.amazon.com/transit-gateway/)

**Matt Klein (criador do Envoy):**

- [mattklein123.dev/appearances — lista oficial de palestras](https://mattklein123.dev/appearances/)
- [Envoy Internals Deep Dive — KubeCon EU 2018 (YouTube)](https://www.youtube.com/watch?v=gQF23Vw0keg)
- [Envoy - Embracing a Service Mesh — QCon NY 2018 (YouTube)](https://www.youtube.com/watch?v=55yi4MMVBi4)
- [The Mechanics of Deploying Envoy at Lyft — KubeCon NA 2017 (YouTube)](https://www.youtube.com/watch?v=IeJDjq-COjk)
- [Lyft's Envoy: Experiences Operating a Large Service Mesh — SREcon Americas 2017 (YouTube)](https://www.youtube.com/watch?v=fZBf6MLjNqM)
- [Service Mesh Day Keynote 2019 (YouTube)](https://www.youtube.com/watch?v=IC0eduR56DA)
- [Kubernetes Podcast from Google — Ep. 221: Creating Envoy, with Matt Klein](https://kubernetespodcast.com/episode/221-creating-envoy/)
- [Blog — Service mesh data plane vs. control plane](https://blog.envoyproxy.io/service-mesh-data-plane-vs-control-plane-2774e720f7fc)
- [Blog — Introduction to modern network load balancing and proxying](https://blog.envoyproxy.io/introduction-to-modern-network-load-balancing-and-proxying-a57f6ff80236)

**Material complementar deste repositório:**

- [README.md](README.md)
- [day1.md](day1.md) — instalação do Istio em modo `ambient` + Bookinfo
- [day2.md](day2.md) — traffic management (routing, shifting, fault injection, mirroring)
- [OLD/day1-sidecar.md](OLD/day1-sidecar.md) — instalação do Istio em modo `sidecar`
- [INSTALL_REQUIREMENTS.md](INSTALL_REQUIREMENTS.md) — Docker, Helm, kubectl, kind, MetalLB
- [UNINSTALL_ISTIO.md](UNINSTALL_ISTIO.md)
