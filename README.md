<!-- TOC -->

- [learning-istio](#learning-istio)
- [Configurações do Ambiente](#configurações-do-ambiente)
- [Requisitos para executar o Istio](#requisitos-para-executar-o-istio)
- [Trilha de estudos (do iniciante ao avançado)](#trilha-de-estudos-do-iniciante-ao-avançado)
- [Day 1](#day-1)
- [Day 2](#day-2)
- [Day 3](#day-3)
- [Day 4](#day-4)
- [Day 5](#day-5)
- [Day 6](#day-6)
- [Guia de Estudos: Istio e Envoy](#guia-de-estudos-istio-e-envoy)
- [Referências](#referências)
  - [Documentação do Istio](#documentação-do-istio)
  - [Documentação do Envoy e do Envoy Gateway](#documentação-do-envoy-e-do-envoy-gateway)
- [Mantenedores](#mantenedores)
- [Licença](#licença)

<!-- TOC -->

# learning-istio

> Foi combinado com a [LinuxTips](https://www.linuxtips.io), via email, a publicação deste repositório como uma forma de contribuir com a comunidade e ajudar na atualização do material do curso. O próprio [Jefferson Fernando](https://twitter.com/badtux_), representante da LinuxTips, começou a divulgação deste repositório junto à comunidade neste [grupo no Telegram](https://t.me/joinchat/GmIMiRVkN55gwDhTlDWCqA).

Esta documentação foi gerada durante o estudo do curso **[Descomplicando Istio](https://www.linuxtips.io/products/descomplicando-o-istio)** da [LinuxTips](https://www.linuxtips.io), como uma forma de aprender a usar o Istio e os conceitos de Service mesh.

**É recomendado ter um bom domínio sobre orquestração de conteinêres com Docker e Kubernetes** antes de começar a estudar o Istio. Os links citados nas referências podem ajudar nisso. Você também pode aprender através dos cursos [Descomplicando Docker](https://www.linuxtips.io/products/descomplicando-o-docker) e [Descomplicando Kubernetes](https://www.linuxtips.io/products/descomplicando-o-kubernetes).

Para fins de aprendizado, o Kubernetes foi provisionado utilizando o [kind](https://kind.sigs.k8s.io).

**Em ambientes de teste, homologação e produção**, que ficam na cloud, é uma boa ideia utilizar serviços gerenciados como: [EKS](https://aws.amazon.com/eks), [GKE](https://cloud.google.com/kubernetes-engine), [AKS](https://azure.microsoft.com/en-us/free/kubernetes-service), [DOKS](https://www.digitalocean.com/products/kubernetes/), entre outros. Em ambientes on-premisses, o [k0s](https://k0sproject.io) e o [kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/) podem ser utilizados.

**Em ambientes de desenvolvimento** é uma boa ideia utilizar o [kind](https://kind.sigs.k8s.io), [microk8s](https://microk8s.io), [k3d](https://k3d.io) ou [minikube](https://minikube.sigs.k8s.io/).

# Configurações do Ambiente

As versões dos softwares utilizados foram:

- Docker: 29.2.1
- Kubernetes (kind): 1.37.0
- kubectl: 1.37.2
- Helm: 3.22.1
- Istio: 1.32
- istioctl: 1.32.0

# Requisitos para executar o Istio

Veja os requisitos de CPU e memória para executar o Istio: https://istio.io/v1.32/docs/setup/platform-setup/docker/

Veja as provedores de Kubernetes suportados: https://istio.io/v1.32/docs/ambient/install/platform-prerequisites/

Veja as versões de Kubernetes suportadas por cada versão do Istio: https://istio.io/v1.32/docs/releases/supported-releases/#support-status-of-istio-releases

Veja mais detalhes sobre os requisitos da aplicação em:

- https://istio.io/v1.32/docs/ops/deployment/application-requirements/
- https://istio.io/v1.32/docs/ops/deployment/platform-requirements/

# Trilha de estudos (do iniciante ao avançado)

| Nível | Conteúdo | Modo sidecar | Modo ambient |
|---|---|---|---|
| Iniciante | [Day 1](day1.md) / [Day 1 - sidecar](day1-sidecar.md): cluster kind, instalação do Istio via Helm, Bookinfo e addons (Grafana, Prometheus, Kiali, Jaeger) | Sim | Sim |
| Iniciante | [Day 2](day2.md): gerenciamento de tráfego (request routing, traffic shifting, fault injection, mirroring) | Sim | Sim |
| Iniciante a avançado | [Guia de Estudos](ISTIO_ENVOY_GUIA_DE_ESTUDOS.md): glossário, arquitetura do Istio e do Envoy, xDS, observabilidade, segurança e 14 laboratórios de troubleshooting | Sim | Sim |
| Intermediário a avançado | [Day 3](day3.md): multicluster e mesh única (multi-primary em redes diferentes, east-west gateway, confiança comum) | Sim | Sim |
| Iniciante a avançado | [Day 4](day4.md): exportação de métricas para Prometheus/Grafana (addons, Prometheus existente, produção) e Datadog | Sim | Sim |
| Intermediário a avançado | [Day 5](day5.md): EnvoyFilter, filtros do Envoy, compressão (gzip, brotli, zstd), Lua e rate limit local | Sim | Gateways (em waypoints, use `TrafficExtension`) |
| Iniciante a avançado | [Day 6](day6.md): Envoy Gateway (Gateway API, políticas, compressão, observabilidade, EnvoyPatchPolicy) | Não usa mesh | Não usa mesh |

# Day 1

Veja os comandos [aqui](day1.md) (modo ambient) ou [aqui](day1-sidecar.md) (modo sidecar).

# Day 2

Veja os comandos [aqui](day2.md)

# Day 3

Multicluster e mesh única, nos modos sidecar e ambient. Veja os comandos [aqui](day3.md).

# Day 4

Exportação de métricas para o Prometheus, Grafana e Datadog. Veja os comandos [aqui](day4.md).

# Day 5

EnvoyFilter, filtros do Envoy e compressão. Veja os comandos [aqui](day5.md).

# Day 6

Envoy Gateway. Veja os comandos [aqui](day6.md).

# Guia de Estudos: Istio e Envoy

Glossário, diagramas (mermaid), comandos e exercícios de troubleshooting sobre Istio (modos sidecar e ambient) e Envoy (arquitetura interna, xDS e o ciclo de vida de uma requisição), com todas as referências oficiais citadas. Veja [aqui](ISTIO_ENVOY_GUIA_DE_ESTUDOS.md).

# Referências

- https://speakerdeck.com/aeciopires/introducao-a-service-mesh-usando-o-istio
- https://github.com/badtuxx/DescomplicandoKubernetes
- http://blog.aeciopires.com/primeiros-passos-com-docker/
- https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
- https://www.alibabacloud.com/blog/kubernetes-configure-liveness-and-readiness-probes_594833 
- https://medium.com/sitewards/deploying-on-kubernetes-10-health-checking-a4986e807afe 
- https://support.sisense.com/kb/en/article/pod-in-crashloopbackoff-state-readinessliveness-probe-failed-get-httppod-ip8082actuatorhealth-dial-tcp-pod-ip8082-connect-connection-refused 
- https://discuss.konghq.com/t/container-ingress-controller-failed-liveness-probe/6796
- https://www.udemy.com/course/istio-hands-on-for-kubernetes/
- https://blog.christianposta.com/tags/istio/
- https://www.envoyproxy.io/docs/envoy/latest/intro/life_of_a_request
- https://github.com/Apress/Practical-Istio

## Documentação do Istio

- https://istio.io/v1.32/docs/overview/what-is-istio/
- https://istio.io/v1.32/docs/overview/why-choose-istio/
- https://istio.io/v1.32/docs/overview/dataplane-modes/
- https://istio.io/v1.32/docs/concepts/
- https://istio.io/v1.32/docs/setup/
- https://istio.io/v1.32/docs/ambient/
- https://istio.io/v1.32/docs/tasks/
- https://istio.io/v1.32/docs/examples/
- https://istio.io/v1.32/docs/ops/

## Documentação do Envoy e do Envoy Gateway

- https://www.envoyproxy.io/docs/envoy/latest/
- https://gateway.envoyproxy.io
- https://gateway-api.sigs.k8s.io

# Mantenedores

- Aécio dos Santos Pires ([linkedin.com/in/aeciopires](https://www.linkedin.com/in/aeciopires))

# Licença

GPL-3.0 Aécio dos Santos Pires
