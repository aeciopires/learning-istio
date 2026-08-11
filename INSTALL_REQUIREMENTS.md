<!-- TOC -->

- [Requisitos](#requisitos)
- [Pacotes gerais](#pacotes-gerais)
- [mise](#mise)
  - [Opção 1 (recomendada, qualquer versão do Ubuntu): script de instalação](#opção-1-recomendada-qualquer-versão-do-ubuntu-script-de-instalação)
  - [Opção 2: pacote `apt` (Ubuntu 26.04+)](#opção-2-pacote-apt-ubuntu-2604)
  - [Opção 3: pacote `apt` via `extrepo` (Ubuntu 24.04, Debian 11+)](#opção-3-pacote-apt-via-extrepo-ubuntu-2404-debian-11)
- [Docker](#docker)
- [Helm](#helm)
- [kubectl](#kubectl)
- [istioctl](#istioctl)
- [kind](#kind)
- [MetalLB](#metallb)

<!-- TOC -->

# Requisitos

# Pacotes gerais

Instale os seguintes pacotes de acordo com o sistema operacional.

Ubuntu 22.04/24.04/26.04:

```bash
sudo apt install -y vim telnet netcat-openbsd git elinks curl wget openssl net-tools jq
```

# mise

O [mise](https://mise.jdx.dev) (antigo `rtx`) é o gerenciador de versões de
ferramentas usado neste repositório. Ele é compatível com o formato de
arquivo `.tool-versions` do [asdf](https://asdf-vm.com) (o mesmo arquivo
`.tool-versions` já usado por este repositório continua funcionando sem
mudanças), mas não depende de plugins externos para ferramentas comuns
como Helm, kubectl e kind - elas já vêm no registro embutido do mise.

> Atenção!!! Para atualizar o mise utilize o seguinte comando:

```bash
mise self-update
```

Escolha **uma** das opções de instalação abaixo, de acordo com a sua
versão do Ubuntu.

## Opção 1 (recomendada, qualquer versão do Ubuntu): script de instalação

```bash
# Opcional: fixe a versão do mise (senão instala a mais recente disponível)
# MISE_VERSION="v2026.8.0"

curl https://mise.run | sh
# ou, fixando a versão:
# curl https://mise.run | MISE_VERSION=$MISE_VERSION sh

# Adicionando no $HOME/.bashrc
echo 'eval "$(~/.local/bin/mise activate bash)"' >> ~/.bashrc

# Garante que os "shims" do mise fiquem no PATH também em shells não
# interativos (scripts, Makefile, etc.) - necessário para os scripts em
# agents/upgrade/scripts/*.sh encontrarem kind/helm/kubectl.
echo 'export PATH="$HOME/.local/share/mise/shims:$PATH"' >> ~/.bashrc

source ~/.bashrc
mise --version
```

## Opção 2: pacote `apt` (Ubuntu 26.04+)

```bash
sudo add-apt-repository -y ppa:jdxcode/mise
sudo apt update
sudo apt install -y mise

echo 'eval "$(mise activate bash)"' >> ~/.bashrc
echo 'export PATH="$HOME/.local/share/mise/shims:$PATH"' >> ~/.bashrc
source ~/.bashrc
mise --version
```

## Opção 3: pacote `apt` via `extrepo` (Ubuntu 24.04, Debian 11+)

```bash
sudo apt install -y extrepo
sudo extrepo enable mise
sudo apt update
sudo apt install -y mise

echo 'eval "$(mise activate bash)"' >> ~/.bashrc
echo 'export PATH="$HOME/.local/share/mise/shims:$PATH"' >> ~/.bashrc
source ~/.bashrc
mise --version
```

Fonte: https://mise.jdx.dev/installing-mise.html

> Atalho: como este repositório já traz um arquivo
> [`.tool-versions`](.tool-versions) com as versões de `kind`, `helm` e
> `kubectl` já definidas, depois de instalar o mise basta rodar
> `mise install helm kubectl kind` na raiz do repositório para instalar as
> três de uma vez, na versão exata que está pinada (também disponível via
> `make install-tools` em [agents/upgrade](agents/upgrade)). Os comandos
> `mise install <ferramenta>@<versão>` e `mise use --pin --path
> .tool-versions <ferramenta>@<versão>` mostrados abaixo são úteis para
> instalar uma ferramenta isolada ou trocar a versão pinada.
>
> Atenção: `mise install` **sem** argumentos também instala qualquer
> ferramenta pinada na sua configuração *global* do mise
> (`~/.config/mise/config.toml`), não só as deste repositório - se isso
> falhar por causa de alguma ferramenta não relacionada, use
> `mise install helm kubectl kind` (com os nomes explícitos) para
> restringir a instalação só ao que este repositório precisa.

# Docker

O [Docker](https://www.docker.com) é a engine de containers usada para
executar o cluster Kubernetes local (via kind) e os containers da
aplicação de exemplo Bookinfo.

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

O [Helm](https://helm.sh) é o gerenciador de pacotes do Kubernetes,
usado neste repositório para instalar o Istio e o MetalLB.

Execute os seguintes comandos para instalar o helm:

> Antes de continuar, se tiver o helm instalado via apt, remova-o com os seguintes comandos:

```bash
sudo apt remove helm
# ou
sudo rm /usr/local/bin/helm
sudo rm /etc/apt/sources.list.d/helm-stable-debian.list
```

> Antes de prosseguir, certifique-se de ter instalado o comando [mise](#mise).

Documentação: https://helm.sh/docs/

```bash
VERSION="3.21.2"

mise latest helm
mise install helm@$VERSION
mise ls helm

# Definindo a versão padrão deste repositório (grava em .tool-versions)
mise use --pin --path .tool-versions helm@$VERSION
mise ls helm
```

# kubectl

O `kubectl` é a ferramenta de linha de comando para interagir com o
cluster Kubernetes, usada em praticamente todos os passos deste
repositório.

Execute os seguintes comandos.

Documentação: https://kubernetes.io/docs/reference/kubectl/overview/

```bash
VERSION_OPTION_1="1.36.3"

mise latest kubectl
mise install kubectl@$VERSION_OPTION_1
mise ls kubectl

# Definindo a versão padrão deste repositório (grava em .tool-versions)
mise use --pin --path .tool-versions kubectl@$VERSION_OPTION_1
mise ls kubectl
```

> Se algum programa/script não interativo não encontrar o `kubectl` (por
> não passar pelos "shims" do mise), crie um link simbólico como reforço:
> `sudo ln -s $HOME/.local/share/mise/shims/kubectl /usr/local/bin/kubectl`

# istioctl

O `istioctl` é a ferramenta de linha de comando do Istio. Não é
estritamente necessária para instalar a malha (a instalação é feita via
Helm - veja [day1.md](day1.md)), mas é bastante útil para **troubleshooting**:
inspecionar o estado dos proxies (`istioctl proxy-status`), validar a
configuração da malha em busca de problemas (`istioctl analyze`),
inspecionar a configuração efetiva de um proxy Envoy
(`istioctl proxy-config`), entre outros comandos de diagnóstico.

> Antes de prosseguir, certifique-se de ter instalado o comando [mise](#mise).

Documentação: https://istio.io/latest/docs/reference/commands/istioctl/

```bash
VERSION="1.30.3"

mise latest istioctl
mise install istioctl@$VERSION
mise ls istioctl

# Definindo a versão padrão deste repositório (grava em .tool-versions)
mise use --pin --path .tool-versions istioctl@$VERSION
mise ls istioctl
```

Verifique a instalação e a conexão com o cluster (depois que o Istio
estiver instalado, veja [day1.md](day1.md)):

```bash
istioctl version
istioctl proxy-status
```

# kind

O kind (Kubernetes in Docker) é outra alternativa para executar o Kubernetes num ambiente local para testes e aprendizado, mas não é recomendado para uso em produção.

Para instalar o kind execute os seguintes comandos.

> Antes de continuar, se tiver o kind instalado, remova-o com o seguinte comando:

```bash
sudo rm /usr/local/bin/kind
```

> Antes de prosseguir, certifique-se de ter instalado o comando [mise](#mise).

```bash
VERSION="0.32.0"

mise latest kind
mise install kind@$VERSION
mise ls kind
# Definindo a versão padrão deste repositório (grava em .tool-versions)
mise use --pin --path .tool-versions kind@$VERSION
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

# Config compatible with kind v0.32.0
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  podSubnet: "10.244.0.0/16"
  serviceSubnet: "10.96.0.0/12"
nodes:
  - role: control-plane
    image: kindest/node:v1.35.5@sha256:ce977ae6d65918d0b58a5f8b5e940429c2ce42fa3a5619ec2bbc60b949c0ac95
    kubeadmConfigPatches:
    - |
      kind: InitConfiguration
      apiVersion: kubeadm.k8s.io/v1beta4
      nodeRegistration:
        kubeletExtraArgs:
          - name: nodeapp
            value: "loadbalancer"
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
    image: kindest/node:v1.35.5@sha256:ce977ae6d65918d0b58a5f8b5e940429c2ce42fa3a5619ec2bbc60b949c0ac95
  - role: worker
    image: kindest/node:v1.35.5@sha256:ce977ae6d65918d0b58a5f8b5e940429c2ce42fa3a5619ec2bbc60b949c0ac95
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

MetalLB é uma implementação de balanceador de carga para clusters Kubernetes bare metal, utilizando protocolos de roteamento padrão.

O Kubernetes não oferece uma implementação de balanceadores de carga de rede (serviços do tipo LoadBalancer) para clusters bare-metal. As implementações de balanceadores de carga de rede que o Kubernetes fornece são todas código de integração que se conecta a várias plataformas IaaS (GCP, AWS, Azure...). Se você não estiver executando em uma plataforma IaaS compatível (GCP, AWS, Azure...), os LoadBalancers permanecerão no estado "pendente" indefinidamente após serem criados.

Os operadores de clusters bare-metal ficam limitados a duas ferramentas menos eficazes para direcionar o tráfego de usuários para seus clusters: os serviços "NodePort" e "externalIPs". Ambas as opções apresentam desvantagens significativas para uso em produção, o que torna os clusters bare-metal cidadãos de segunda classe no ecossistema Kubernetes.

O MetalLB visa corrigir esse desequilíbrio, oferecendo uma implementação de balanceador de carga de rede que se integra com equipamentos de rede padrão, para que os serviços externos em clusters bare-metal também "funcionem" da melhor forma possível.

Execute os seguintes comandos para instalar o [MetalLB](https://metallb.universe.tf):

> Referências:
> - https://github.com/rafaelperoco/kind/blob/main/createCluster.sh
> - https://metallb.universe.tf/installation/

```bash
export HELM_CHART_VERSION="0.16.1"

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
  --version $HELM_CHART_VERSION \
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
# IPv4 can be in 0 or 1 array position of IPAM result of command bellow
KIND_ADDRESS=$(docker network inspect kind | jq -r '.[].IPAM.Config[] | select(.Subnet | test("^[0-9]+\\.")) | .Subnet' | cut -d"." -f1-3)

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
