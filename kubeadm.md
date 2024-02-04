# Using k8s with kubeadm

Crie três máquinas virtuais ou instâncias no cloud provider com as seguites configurações.

AWS-EC2 Instances:

* type: on-demand (t2.large 2 CPU e 8 GB de memória)
* HD: 30 GB
* SO: Ubuntu 22.04 64 bits
* login: ubuntu
* SSH: 22/TCP
* key: teste-aecio-treinamento-istio

```bash
ssh -o ServerAliveInterval=30 -i ~/teste-aecio-treinamento-istio.pem ubuntu@master
ssh -o ServerAliveInterval=30 -i ~/teste-aecio-treinamento-istio.pem ubuntu@worker1
ssh -o ServerAliveInterval=30 -i ~/teste-aecio-treinamento-istio.pem ubuntu@worker2
```

Execute os seguintes comandos de acordo com a instância.

```bash
#------- Specifics for each instance
sudo hostnamectl set-hostname master
sudo hostnamectl set-hostname worker1
sudo hostnamectl set-hostname worker2


#------- Generic (all instances)
sudo systemctl stop ufw
sudo systemctl disable ufw
sudo apt update
sudo apt upgrade -y
sudo apt install -y vim traceroute telnet git tcpdump elinks curl wget openssl netcat net-tools jq gpg acl uidmap

# Add aliases in files: ~/.bashrc and /root/.bashrc
alias egrep='egrep --color=auto'
alias fgrep='fgrep --color=auto'
alias grep='grep --color=auto'
alias k='kubectl'
alias l='ls -CF'
alias la='ls -A'
alias live='curl parrot.live'
alias ll='ls -alF'
alias ls='ls --color=auto'
alias nettools='kubectl run --rm -it nettools-$(< /dev/urandom tr -dc a-z-0-9 | head -c${1:-4}) --image=aeciopires/nettools:1.0.0 -n default -- bash'
alias randompass='< /dev/urandom tr -dc _A-Z-a-z-0-9 | head -c${1:-16}'
alias sc='source ~/.bashrc'
alias show-hidden-files='du -sch .[!.]* * |sort -h'
alias gitlog='git log -p'

# Apply new aliases
source ~/.bashrc
source /root/.bashrc

# Reboot system
sudo reboot

# Check whether the Cgroup driver has been set correctly
# If the output was Cgroup Driver: systemd, all right!
docker info | grep -i cgroup

# Reference: https://github.com/containerd/containerd/discussions/8033
# Edit /etc/containerd/config.toml file and uncomment the line
# Before:
#  disabled_plugins = ["cri"]
# After:
#  #disabled_plugins = ["cri"]
# Now, run the command:
sudo systemctl restart containerd.service


# Install Kubernetes with kubeadm
# References:
# https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
K8S_VERSION='1.29'
sudo swapoff -a
sudo curl -fsSL "https://pkgs.k8s.io/core:/stable:/v${K8S_VERSION}/deb/Release.key" | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v${K8S_VERSION}/deb/ /" | sudo tee /etc/apt/sources.list.d/kubernetes.list
cat /etc/apt/sources.list.d/kubernetes.list
sudo apt update
sudo apt install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
kubectl version --client




#------- Specifics (master)
# Config allow system network comunication
sudo modprobe br_netfilter ip_vs_rr ip_vs_wrr ip_vs_sh nf_conntrack_ipv4 ip_vs
sudo sysctl net.ipv4.conf.all.forwarding=1
echo "net.ipv4.conf.all.forwarding=1" | sudo tee -a /etc/sysctl.conf


sudo kubeadm config images pull
sudo kubeadm init
# Reset
# sudo kubeadm reset
# sudo rm -rf /etc/cni/net.d

# Config kubectl
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
kubectl get nodes

# Install weavenet or calico
# WeaveNet:
#   https://www.weave.works/docs/net/latest/kubernetes/kube-addon/#install
#   https://www.pluralsight.com/cloud-guru/labs/aws/setting-up-kubernetes-networking-with-weave-net
# Calico:
#   https://docs.tigera.io/calico/latest/getting-started/kubernetes/self-managed-onprem/onpremises


# Install MetalLB following the instructions of the page: https://metallb.universe.tf/installation/


#------- Specifics (worker1 and worker2)
# Allow all these ports: https://github.com/badtuxx/DescomplicandoKubernetes/blob/main/pt/day_one/descomplicando_kubernetes.md#portas-que-devemos-nos-preocupar
#
# In master node print a join command for add worker node in cluster
sudo kubeadm token create --print-join-command
#
# Example of command to run in worker node:
kubeadm join 172.31.17.64:6443 --token st7hlu.6hbl4c6f2crh08ys \
	--discovery-token-ca-cert-hash sha256:331afc84c3c71a369d314ad3be27d738ccc1535fa2492876fcfaec42c5fa3135
