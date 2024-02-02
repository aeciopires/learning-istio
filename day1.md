# Descomplicando Istio

<!-- TOC -->

- [Descomplicando Istio](#descomplicando-istio)
- [Day 1](#day-1)

<!-- TOC -->


# Day 1

Comandos executados em sequencia durante o treinamento.

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

# Install Docker
# References:
# https://github.com/badtuxx/DescomplicandoKubernetes/blob/main/pt/day_one/descomplicando_kubernetes.md#instala%C3%A7%C3%A3o-do-docker-e-do-kubernetes
sudo su
cd /home/ubuntu
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

dockerd-rootless-setuptool.sh install
sudo systemctl daemon-reload
sudo systemctl restart docker
sudo systemctl restart containerd
sudo systemctl enable docker
sudo systemctl enable containerd
sudo usermod -aG docker $USER
sudo setfacl -m user:$USER:rw /var/run/docker.sock
sudo setfacl -m user:$USER:rw /var/run/containerd/containerd.sock

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



#------- Specifics (worker1 and worker2)
# Allow all these ports: https://github.com/badtuxx/DescomplicandoKubernetes/blob/main/pt/day_one/descomplicando_kubernetes.md#portas-que-devemos-nos-preocupar
#
# In master node print a join command for add worker node in cluster
sudo kubeadm token create --print-join-command
#
# Example of command to run in worker node:
kubeadm join 172.31.17.64:6443 --token st7hlu.6hbl4c6f2crh08ys \
	--discovery-token-ca-cert-hash sha256:331afc84c3c71a369d314ad3be27d738ccc1535fa2492876fcfaec42c5fa3135





#------- Specifics (master)
# References:
# https://istio.io/latest/docs/setup/getting-started/
# Allow all these ports: https://istio.io/latest/docs/ops/deployment/requirements/
curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.11.4 sh -
export ISTIO_DIR_BASE="/home/ubuntu/istio-1.11.4"
cd $ISTIO_DIR_BASE
export PATH="$PATH:$ISTIO_DIR_BASE/bin"

# Complementary files
git clone https://github.com/aeciopires/learning-istio
export COMPLEMENTARY_FILES=/home/ubuntu/learning-istio/files

istioctl x precheck

# Install profile of demonstration
istioctl install --set profile=demo -y

# Install Addons (grafana, prometheus, kiali, jaeger, zipkin)
kubectl apply -f $ISTIO_DIR_BASE/samples/addons
kubectl rollout status deployment/kiali -n istio-system
kubectl rollout status deployment/grafana -n istio-system
kubectl rollout status deployment/prometheus -n istio-system
kubectl rollout status deployment/jaeger -n istio-system

# List api-resources and crd
kubectl api-resources | grep istio
kubectl get crd | grep istio
kubectl get svc -n istio-system
kubectl get pods -n istio-system

# First Deployment
kubectl label namespace default istio-injection=enabled

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/platform/kube/bookinfo.yaml
kubectl get services
kubectl get pods

# Verify everything is working correctly up to this point. Run this command to see if the app is running inside the cluster and serving HTML pages by checking for the page title in the response:
kubectl exec "$(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl -sS productpage:9080/productpage | grep -o "<title>.*</title>"

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/bookinfo-gateway.yaml
kubectl get gateway
kubectl get virtualservices -o yaml
kubectl get services
kubectl get pods
kubectl get virtualservices.networking.istio.io bookinfo -o yaml
kubectl get virtualservices.networking.istio.io
kubectl describe virtualservices.networking.istio.io bookinfo
kubectl get gateways.networking.istio.io
kubectl describe gateways.networking.istio.io bookinfo-gateway

# Port forward in background for access Productpage application
# Allow port 9080/TCP
kubectl port-forward svc/productpage 9080:9080 -n default --address=0.0.0.0 > /dev/null 2>&1 &
# Access URL: http://master:9080/productpage

# Port forward in background for access Kiali
# Allow port 20001/TCP
kubectl port-forward svc/kiali 20001:20001 -n istio-system --address=0.0.0.0 > /dev/null 2>&1 &
# Access URL: http://master:20001/kiali

# Port forward in background for access Grafana
# Allow port 3000/TCP
kubectl port-forward svc/grafana 3000:3000 -n istio-system --address=0.0.0.0 > /dev/null 2>&1 &
# Access URL: http://master:3000

# Port forward in background for access Prometheus
# Allow port 9090/TCP
kubectl port-forward svc/prometheus 9090:9090 -n istio-system --address=0.0.0.0 > /dev/null 2>&1 &
# Access URL: http://master:9090

# Port forward in background for access Jaeger
# Allow port 16685/TCP
kubectl port-forward svc/tracing 16685:80 -n istio-system --address=0.0.0.0 > /dev/null 2>&1 &
# Access URL: http://master:16685

# List jobs running in background
jobs -l

# Finish jobs running in background
kill %1 %2 %3 %4 %5
```