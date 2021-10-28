# learning-istio

Esta documentação foi gerada durante o estudo do curso **[Descomplicando Istio](https://www.linuxtips.io/products/descomplicando-o-istio)** da [LinuxTips](https://www.linuxtips.io), como uma forma de aprender a usar o Istio e os conceitos de Service mesh.

É recomendado ter um bom domínio sobre orquestração de conteinêres com Docker e Kubernetes antes de começar a estudar o Istio. Os links citados nas referências podem ajudar nisso.

Para fins de aprendizado, o Kubernetes foi provisionado manualmente em instâncias EC2 para entender os requisitos de rede, hardware, software e o funcionamento dos componentes do Kubernetes e Istio. Mas em ambientes de produção é uma boa ideia utilizar serviços gerenciados com o o [EKS](https://aws.amazon.com/eks), [GKE](https://cloud.google.com/kubernetes-engine), [AKS](https://azure.microsoft.com/en-us/free/kubernetes-service), [DOKS](https://www.digitalocean.com/products/kubernetes/), entre outros.

> Para o Kubernetes e o Istio funcionarem corretamente, foi necessário liberar um conjunto de 40 portas no security group associado as instâncias EC2.

# Day 1

Veja os comandos [aqui](day1.md)

# Day 2

Veja os comandos aqui.

# Day 3

Veja os comandos aqui.

# Day 4

Veja os comandos aqui.

# Referências

* https://github.com/badtuxx/DescomplicandoKubernetes
* https://github.com/badtuxx/DescomplicandoKubernetes/blob/main/pt/day_one/descomplicando_kubernetes.md
* http://blog.aeciopires.com/primeiros-passos-com-docker/
* https://helm.sh/docs/
* https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
* https://istio.io/latest/docs/setup/getting-started/
* https://istio.io/latest/docs/ops/deployment/requirements/
* https://istio.io/latest/docs/setup/install/istioctl/
* https://www.alibabacloud.com/blog/kubernetes-configure-liveness-and-readiness-probes_594833 
* https://medium.com/sitewards/deploying-on-kubernetes-10-health-checking-a4986e807afe 
* https://githubmemory.com/repo/istio/istio/issues/32963 
* https://support.sisense.com/kb/en/article/pod-in-crashloopbackoff-state-readinessliveness-probe-failed-get-httppod-ip8082actuatorhealth-dial-tcp-pod-ip8082-connect-connection-refused 
* https://forums.rancher.com/t/liveness-probe-failed-connection-refused/20837 
* https://discuss.konghq.com/t/container-ingress-controller-failed-liveness-probe/6796  
