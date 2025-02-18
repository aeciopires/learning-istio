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

Instale o **Docker** em cada instância com as instruções da página: [docker](docker.md).

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

# Reference:
# https://comtechies.com/setup-kubernetes-cluster-kubeadm.html
# https://devopscube.com/setup-kubernetes-cluster-kubeadm/
# Config allow system network comunication
# Execute the following commands on all the nodes to enable IPtables bridged traffic and to enable SWAP.
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter

# sysctl params required by setup, params persist across reboots
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

# Apply sysctl params without reboot
sudo sysctl --system

sudo swapoff -a
(crontab -l 2>/dev/null; echo "@reboot /sbin/swapoff -a") | crontab - || true

# Install Kubernetes with kubeadm
# References:
# https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
# https://earthly.dev/blog/managing-k8s-with-kubeadm/
# https://www.densify.com/kubernetes-tools/kubeadm/
# https://stackoverflow.com/questions/43201643/kubeadm-and-weave-not-working-together#43204264
# https://devopscube.com/setup-kubernetes-cluster-kubeadm/
K8S_VERSION='1.29'
sudo swapoff -a
sudo curl -fsSL "https://pkgs.k8s.io/core:/stable:/v${K8S_VERSION}/deb/Release.key" | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v${K8S_VERSION}/deb/ /" | sudo tee /etc/apt/sources.list.d/kubernetes.list
cat /etc/apt/sources.list.d/kubernetes.list
sudo apt update
sudo apt install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
kubectl version --client


# Add the node IP to KUBELET_EXTRA_ARGS.
INTEFACE_NAME='eth0'
sudo apt-get install -y jq
export local_ip="$(ip --json a s | jq -r --arg interface "$INTEFACE_NAME" '.[] | if .ifname == $interface then .addr_info[] | if .family == "inet" then .local else empty end else empty end')"
sudo bash -c "cat > /etc/default/kubelet << EOF
KUBELET_EXTRA_ARGS=--node-ip=$local_ip
EOF"

#------- Specifics (master)
IPADDR="$local_ip"
NODENAME=$(hostname -s)
POD_CIDR="192.168.0.0/16"

sudo kubeadm config images pull
sudo kubeadm init --apiserver-advertise-address=$IPADDR  --apiserver-cert-extra-sans=$IPADDR  --pod-network-cidr=$POD_CIDR --node-name $NODENAME --ignore-preflight-errors Swap
# Reset
# sudo kubeadm reset
# sudo rm -rf /etc/cni/net.d

# Config kubectl
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
kubectl get nodes

# You verify all the cluster component health statuses using the following command.
kubectl get --raw='/readyz?verbose'
kubectl cluster-info

# Reference:
# https://docs.tigera.io/calico/latest/getting-started/kubernetes/self-managed-onprem/onpremises
# Kubeadm does not configure any network plugin. You need to install a network plugin of your choice for kubernetes pod networking and enable network policy.
# I am using the Calico network plugin for this setup.
kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/tigera-operator.yaml

curl https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/custom-resources.yaml -O

kubectl create -f custom-resources.yaml

# Instale o **MetalLB** com as instruções da página: [metallb](kind.md#metallb).

#------- Specifics (worker1 and worker2)
# Allow all these ports: https://github.com/badtuxx/DescomplicandoKubernetes/blob/main/pt/day_one/descomplicando_kubernetes.md#portas-que-devemos-nos-preocupar
#
# In master node print a join command for add worker node in cluster
sudo kubeadm token create --print-join-command
#
# Example of command to run in worker node:
sudo kubeadm join 172.31.23.105:6443 --token iax219.i23cr4o1hcnln4t7 \
	--discovery-token-ca-cert-hash sha256:daf0c3356f3ec0c108d492c71ff3d719ef4ac4e5af4bdfcca786a1e0c123e469
