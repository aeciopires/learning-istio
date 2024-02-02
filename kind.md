<!-- TOC -->

- [Requirements](#requirements)
- [General Packages](#general-packages)
- [Docker](#docker)
- [Kubectl](#kubectl)
- [Kind](#kind)

<!-- TOC -->

# Requirements

# General Packages

Instale os seguintes pacotes de acordo com o sistema operacional.

Ubuntu 20.04/22.04:

```bash
sudo apt install -y vim telnet git curl wget openssl netcat net-tools jq
```

# Docker

Instale o Docker CE (Community Edition) com os seguintes comandos:

```bash
sudo su
curl -fsSL https://get.docker.com | bash

cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF
exit

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

# Kubectl

Instale o kubectl com os seguintes comandos.

```bash
sudo su

VERSION=v1.29.1
KUBECTL_BIN=kubectl

function install_kubectl {
if [ -z $(which $KUBECTL_BIN) ]; then
    curl -LO https://storage.googleapis.com/kubernetes-release/release/$VERSION/bin/linux/amd64/$KUBECTL_BIN
    chmod +x ${KUBECTL_BIN}
    mv ${KUBECTL_BIN} /usr/local/bin/${KUBECTL_BIN}
    ln -sf /usr/local/bin/${KUBECTL_BIN} /usr/bin/${KUBECTL_BIN}
else
    echo "Kubectl is most likely installed"
fi
}

install_kubectl

which kubectl

kubectl version --client

exit
```

Referências:
* https://kubernetes.io/docs/reference/kubectl/overview/
* https://kubernetes.io/docs/tasks/tools/install-kubectl/

# Kind

Instale o kind com os seguintes comandos.

```bash
VERSION=v0.20.0
curl -Lo ./kind https://kind.sigs.k8s.io/dl/$VERSION/kind-$(uname)-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

Para criar um cluster com múltiplos nós locais com o Kind, crie um arquivo do tipo YAML para definir a quantidade e o tipo de nós no cluster que você deseja. 

No exemplo a seguir, será criado o arquivo $HOME/kind-3nodes.yaml para especificar um cluster com 1 nó master (que executará o control plane do Kubernetes) e 2 workers (que executará o data plane do Kubernetes).

```yaml
cat << EOF > $HOME/kind-3nodes.yaml
# References:
# Kind release image: https://github.com/kubernetes-sigs/kind/releases
# Configuration: https://kind.sigs.k8s.io/docs/user/configuration/
# Metal LB in Kind: https://kind.sigs.k8s.io/docs/user/loadbalancer
# Ingress in Kind: https://kind.sigs.k8s.io/docs/user/ingress

# Config compatible with kind v0.20.0
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  podSubnet: "10.244.0.0/16"
  serviceSubnet: "10.96.0.0/12"
nodes:
  - role: control-plane
    image: kindest/node:v1.29.0@sha256:eaa1450915475849a73a9227b8f201df25e55e268e5d619312131292e324d570
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
    image: kindest/node:v1.29.0@sha256:eaa1450915475849a73a9227b8f201df25e55e268e5d619312131292e324d570
  - role: worker
    image: kindest/node:v1.29.0@sha256:eaa1450915475849a73a9227b8f201df25e55e268e5d619312131292e324d570
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

* https://github.com/badtuxx/DescomplicandoKubernetes/blob/master/day-1/DescomplicandoKubernetes-Day1.md#kind 
* https://kind.sigs.k8s.io
* https://kind.sigs.k8s.io/docs/user/quick-start/
* https://github.com/kubernetes-sigs/kind/releases
* https://kubernetes.io/blog/2020/05/21/wsl-docker-kubernetes-on-the-windows-desktop/#kind-kubernetes-made-easy-in-a-container
