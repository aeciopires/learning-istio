<!-- TOC -->

- [Requisitos](#requisitos)
- [Pacotes gerais](#pacotes-gerais)
- [asdf](#asdf)
- [Docker](#docker)
- [Helm](#helm)
- [kubectl](#kubectl)
- [kind](#kind)
- [MetalLB](#metallb)

<!-- TOC -->

# Requisitos

# Pacotes gerais

Instale os seguintes pacotes de acordo com o sistema operacional.

Ubuntu 24.04:

```bash
sudo apt install -y vim telnet netcat-openbsd git elinks curl wget openssl net-tools jq
```

# asdf

Execute os seguintes comandos:

> Atenção!!! Para atualizar o asdf utilize APENAS o seguinte comando:

```bash
asdf update
```

> Se tentar reinstalar ou atualizar mudando a versão nos comandos seguintes, será necessário reinstalar todos os plugins/comandos instalados antes, por isso é muito importante fazer backup do diretório $HOME/.asdf.

```bash
ASDF_VERSION="v0.15.0"
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch $ASDF_VERSION

# Adicionando no $HOME/.bashrc
echo ". \"\$HOME/.asdf/asdf.sh\"" >> ~/.bashrc
echo ". \"\$HOME/.asdf/completions/asdf.bash\"" >> ~/.bashrc
source ~/.bashrc
```

Fonte: https://asdf-vm.com/guide/introduction.html

# Docker

Instale o Docker CE (Community Edition) com os seguintes comandos:

```bash
curl -fsSL https://get.docker.com | sudo bash

# References:
# https://stackoverflow.com/questions/18836853/sudo-cat-eof-file-doesnt-work-sudo-su-does
# https://earthly.dev/blog/managing-k8s-with-kubeadm/
sudo bash -c 'cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF'


# Add your user to the Docker group
sudo apt install -y acl uidmap
dockerd-rootless-setuptool.sh install
sudo usermod -aG docker $USER
sudo setfacl -m user:$USER:rw /var/run/docker.sock
sudo setfacl -m user:$USER:rw /var/run/containerd/containerd.sock

# Start the Docker service
sudo systemctl start docker
sudo systemctl start containerd

# Configure Docker to boot up with the OS
sudo systemctl enable docker
sudo systemctl enable containerd
```

# Helm

Execute os seguintes comandos para instalar o helm:

> Antes de continuar, se tiver o helm instalado via apt, remova-o com os seguintes comandos:

```bash
sudo apt remove helm
# ou
sudo rm /usr/local/bin/helm
sudo rm /etc/apt/sources.list.d/helm-stable-debian.list
```

> Antes de prosseguir, certifique-se de ter instalado o comando [asdf](#asdf).

Documentação: https://helm.sh/docs/

```bash
VERSION="3.18.6"

asdf plugin list all | grep helm
asdf plugin add helm https://github.com/Antiarchitect/asdf-helm.git
asdf latest helm

asdf install helm $VERSION
asdf list helm

# Definindo a versão padrão
asdf global helm $VERSION
asdf list helm
```

# kubectl

Execute os seguintes comandos.

Documentação: https://kubernetes.io/docs/reference/kubectl/overview/

```bash
VERSION_OPTION_1="1.33.3"

asdf plugin list all | grep kubectl
asdf plugin add kubectl https://github.com/asdf-community/asdf-kubectl.git
asdf latest kubectl

asdf install kubectl $VERSION_OPTION_1
asdf list kubectl

# Definindo a versão padrão
asdf global kubectl $VERSION_OPTION_1
asdf list kubectl

# Criando um link simbólico
sudo ln -s $HOME/.asdf/shims/kubectl /usr/local/bin/kubectl
```

# kind

O kind (Kubernetes in Docker) é outra alternativa para executar o Kubernetes num ambiente local para testes e aprendizado, mas não é recomendado para uso em produção.

Para instalar o kind execute os seguintes comandos.

> Antes de continuar, se tiver o kind instalado, remova-o com o seguinte comando:

```bash
sudo rm /usr/local/bin/kind
```

> Antes de prosseguir, certifique-se de ter instalado o comando [asdf](#asdf).

```bash
VERSION="0.29.0"
asdf plugin list all | grep kind
asdf plugin add kind https://github.com/johnlayton/asdf-kind.git
asdf latest kind
asdf install kind $VERSION
asdf list kind
# Definindo a versão padrão
asdf global kind $VERSION
```

Para criar um cluster com múltiplos nós locais com o Kind, crie um arquivo do tipo YAML para definir a quantidade e o tipo de nós no cluster que você deseja.

No exemplo a seguir, será criado o arquivo ``$HOME/kind-3nodes.yaml`` para especificar um cluster com 1 nó master (que executará o control plane do Kubernetes) e 2 workers (que executará o data plane do Kubernetes).

```bash
cat << EOF > $HOME/kind-3nodes.yaml
# References:
# Kind release image: https://github.com/kubernetes-sigs/kind/releases
# Configuration: https://kind.sigs.k8s.io/docs/user/configuration/
# Metal LB in Kind: https://kind.sigs.k8s.io/docs/user/loadbalancer
# Ingress in Kind: https://kind.sigs.k8s.io/docs/user/ingress

# Config compatible with kind v0.29.0
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  podSubnet: "10.244.0.0/16"
  serviceSubnet: "10.96.0.0/12"
nodes:
  - role: control-plane
    image: kindest/node:v1.33.1@sha256:050072256b9a903bd914c0b2866828150cb229cea0efe5892e2b644d5dd3b34f
    kubeadmConfigPatches:
    - |
      kind: InitConfiguration
      nodeRegistration:
        kubeletExtraArgs:
          node-labels: "nodeapp=loadbalancer"
    extraPortMappings:
    - containerPort: 80
      hostPort: 80
      listenAddress: "0.0.0.0" # Optional, defaults to "0.0.0.0"
      protocol: TCP
    - containerPort: 443
      hostPort: 443
      listenAddress: "0.0.0.0" # Optional, defaults to "0.0.0.0"
      protocol: TCP
  - role: worker
    image: kindest/node:v1.33.1@sha256:050072256b9a903bd914c0b2866828150cb229cea0efe5892e2b644d5dd3b34f
  - role: worker
    image: kindest/node:v1.33.1@sha256:050072256b9a903bd914c0b2866828150cb229cea0efe5892e2b644d5dd3b34f
EOF
```

Crie um cluster chamado ``kind-multinodes`` utilizando as especificações definidas no arquivo ``$HOME/kind-3nodes.yaml``.

```bash
kind create cluster --name kind-multinodes --config $HOME/kind-3nodes.yaml
```

Para visualizar os seus clusters utilizando o kind, execute o comando a seguir.

```bash
kind get clusters
```

Para destruir o cluster, execute o seguinte comando que irá selecionar e remover todos os clusters locais criados no Kind.

```bash
kind delete clusters $(kind get clusters)
```

Referências:
- https://github.com/badtuxx/DescomplicandoKubernetes/blob/master/day-1/DescomplicandoKubernetes-Day1.md#kind
- https://kind.sigs.k8s.io/docs/user/quick-start/
- https://github.com/kubernetes-sigs/kind/releases
- https://kubernetes.io/blog/2020/05/21/wsl-docker-kubernetes-on-the-windows-desktop/#kind-kubernetes-made-easy-in-a-container

Repositório alternativo para uso do kind com nginx-controller, linkerd e outras ferramentas: https://github.com/rafaelperoco/kind

# MetalLB

Execute os seguintes comandos para instalar o [MetalLB](https://metallb.universe.tf):

> Referências:
> - https://github.com/rafaelperoco/kind/blob/main/createCluster.sh
> - https://metallb.universe.tf/installation/

```bash
# actually apply the changes, returns nonzero returncode on errors only
kubectl get configmap kube-proxy -n kube-system -o yaml | \
sed -e "s/strictARP: false/strictARP: true/" | \
kubectl apply -f - -n kube-system

# Install and upgrade Helm repositories
helm repo add metallb https://metallb.github.io/metallb
helm repo update

# Install MetalLB and check if it is installed
helm upgrade --install metallb metallb/metallb \
  --create-namespace \
  --namespace metallb-system \
  --version 0.15.2 \
  --set "controller.tolerations[0].key=node-role.kubernetes.io/master" \
  --set "controller.tolerations[0].effect=NoSchedule" \
  --set speaker.tolerateMaster=true \
  --wait --debug --timeout=900s
```

> Pode demorar uns 5 minutos até os daemonsets do MetalLB ficarem health.

Configure o endereçamento IP a ser usado pelo ingress gerenciado pelo MetalLB.

```bash
# Install jq package
sudo apt update
sudo apt install -y jq

# Get default gateway interface
KIND_ADDRESS=$(docker network inspect kind | jq '.[].IPAM | .Config | .[0].Subnet' | cut -d \" -f 2 | cut -d"." -f1-3)

# Radomize Loadbalancer IP Range
#KIND_ADDRESS_END=$(shuf -i 100-150 -n1)
KIND_ADDRESS_END="100"
NETWORK_SUBMASK="27"

# Create address range
KIND_LB_RANGE="$(echo $KIND_ADDRESS.$KIND_ADDRESS_END)/$NETWORK_SUBMASK"

cat << EOF > $HOME/metallb-ingress-address.yaml
# References:
# https://metallb.universe.tf/configuration/#layer-2-configuration

apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: first-pool
spec:
  addresses:
  - $KIND_LB_RANGE
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: example
spec:
  ipAddressPools:
  - first-pool
EOF

kubectl apply -f $HOME/metallb-ingress-address.yaml --namespace metallb-system
```
