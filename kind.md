<!-- TOC -->

- [Requisitos](#requisitos)
- [Pacotes gerais](#pacotes-gerais)
- [asdf](#asdf)
- [Docker](#docker)
- [Helm](#helm)
- [kubectl](#kubectl)
- [kind](#kind)

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

Instale o **Docker** com as instruções da página: [docker](docker.md).

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
VERSION="3.17.0"

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
VERSION_OPTION_1="1.32.1"

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
VERSION="0.26.0"
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

# Config compatible with kind v0.26.0
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  podSubnet: "10.244.0.0/16"
  serviceSubnet: "10.96.0.0/12"
nodes:
  - role: control-plane
    image: kindest/node:v1.32.0@sha256:c48c62eac5da28cdadcf560d1d8616cfa6783b58f0d94cf63ad1bf49600cb027
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
    image: kindest/node:v1.32.0@sha256:c48c62eac5da28cdadcf560d1d8616cfa6783b58f0d94cf63ad1bf49600cb027
  - role: worker
    image: kindest/node:v1.32.0@sha256:c48c62eac5da28cdadcf560d1d8616cfa6783b58f0d94cf63ad1bf49600cb027
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
