<!-- TOC -->

- [learning-istio](#learning-istio)
- [Configurações do Ambiente](#configurações-do-ambiente)
- [Requisitos para executar o Istio](#requisitos-para-executar-o-istio)
- [Day 1](#day-1)
- [Day 2](#day-2)
- [Referências](#referências)
  - [Documentação do Istio](#documentação-do-istio)
- [Mantenedores](#mantenedores)
- [Licença](#licença)

<!-- TOC -->

# learning-istio

> Foi combinado com a [LinuxTips](https://www.linuxtips.io), via email, a publicação deste repositório como uma forma de contribuir com a comunidade e ajudar na atualização do material do curso. O próprio [Jefferson Fernando](https://twitter.com/badtux_), representante da LinuxTips, começou a divulgação deste repositório junto à comunidade neste [grupo no Telegram](https://t.me/joinchat/GmIMiRVkN55gwDhTlDWCqA).

Esta documentação foi gerada durante o estudo do curso **[Descomplicando Istio](https://www.linuxtips.io/products/descomplicando-o-istio)** da [LinuxTips](https://www.linuxtips.io), como uma forma de aprender a usar o Istio e os conceitos de Service mesh.

**É recomendado ter um bom domínio sobre orquestração de conteinêres com Docker e Kubernetes** antes de começar a estudar o Istio. Os links citados nas referências podem ajudar nisso. Você também pode aprender através dos cursos [Descomplicando Docker](https://www.linuxtips.io/products/descomplicando-o-docker) e [Descomplicando Kubernetes](https://www.linuxtips.io/products/descomplicando-o-kubernetes).

Para fins de aprendizado, o Kubernetes foi provisionado utilizando o [kind](kind.md).

**Em ambientes de teste, homologação e produção**, que ficam na cloud, é uma boa ideia utilizar serviços gerenciados como: [EKS](https://aws.amazon.com/eks), [GKE](https://cloud.google.com/kubernetes-engine), [AKS](https://azure.microsoft.com/en-us/free/kubernetes-service), [DOKS](https://www.digitalocean.com/products/kubernetes/), entre outros. Em ambientes on-premisses, o [k0s](https://k0sproject.io) e o [kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/) podem ser utilizados.

**Em ambientes de desenvolvimento** é uma boa ideia utilizar o [kind](https://kind.sigs.k8s.io/), [microk8s](https://microk8s.io), [k3d](https://k3d.io) ou [minikube](https://minikube.sigs.k8s.io/).

# Configurações do Ambiente

As versões dos softwares utilizados foram:

- Docker: 28.3.2
- Kubernetes (kind): 1.33.1
- kubectl: 1.33.1
- Helm: 3.18.4
- Istio: 1.27

# Requisitos para executar o Istio

Veja os requisitos de CPU e memória para executar o Istio: https://istio.io/v1.27/docs/setup/platform-setup/docker/

Veja as provedores de Kubernetes suportados: https://istio.io/v1.27/docs/ambient/install/platform-prerequisites/

Veja as versões de Kubernetes suportadas por cada versão do Istio: https://istio.io/v1.27/docs/releases/supported-releases/#support-status-of-istio-releases

Veja mais detalhes sobre os requisitos da aplicação em:

- https://istio.io/v1.27/docs/ops/deployment/application-requirements/
- https://istio.io/v1.27/docs/ops/deployment/platform-requirements/

# Day 1

Veja os comandos [aqui](day1.md)

# Day 2

Veja os comandos [aqui](day2.md)

# Referências

- https://github.com/badtuxx/DescomplicandoKubernetes
- http://blog.aeciopires.com/primeiros-passos-com-docker/
- https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
- https://www.alibabacloud.com/blog/kubernetes-configure-liveness-and-readiness-probes_594833 
- https://medium.com/sitewards/deploying-on-kubernetes-10-health-checking-a4986e807afe 
- https://support.sisense.com/kb/en/article/pod-in-crashloopbackoff-state-readinessliveness-probe-failed-get-httppod-ip8082actuatorhealth-dial-tcp-pod-ip8082-connect-connection-refused 
- https://discuss.konghq.com/t/container-ingress-controller-failed-liveness-probe/6796
- https://www.udemy.com/course/istio-hands-on-for-kubernetes/

## Documentação do Istio

- https://istio.io/v1.27/docs/overview/what-is-istio/
- https://istio.io/v1.27/docs/overview/why-choose-istio/
- https://istio.io/v1.27/docs/overview/dataplane-modes/
- https://istio.io/v1.27/docs/concepts/
- https://istio.io/v1.27/docs/setup/
- https://istio.io/v1.27/docs/ambient/
- https://istio.io/v1.27/docs/tasks/
- https://istio.io/v1.27/docs/examples/
- https://istio.io/v1.27/docs/ops/

# Mantenedores

- Aécio dos Santos Pires ([linkedin.com/in/aeciopires](https://www.linkedin.com/in/aeciopires))

# Licença

GPL-3.0 Aécio dos Santos Pires
