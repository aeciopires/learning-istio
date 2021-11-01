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
sudo apt install -y vim traceroute telnet git tcpdump elinks curl wget openssl netcat net-tools jq

# Add aliases in files: ~/.bashrc and /root/.bashrc
alias egrep='egrep --color=auto'
alias fgrep='fgrep --color=auto'
alias grep='grep --color=auto'
alias k='kubectl'
alias kmongo='kubectl run --rm -it mongoshell --image=mongo:latest -n default -- bash'
alias kmysql='kubectl run --rm -it mysql --image=mysql:5.7 -n default -- bash'
alias kredis='kubectl run --rm -it redis-cli --image=redis:latest -n default -- bash'
alias kssh='kubectl run ssh-client -it --rm --image=kroniak/ssh-client -n default -- bash'
alias nettools='kubectl run --rm -it nettools --image=travelping/nettools:latest -n default -- bash'
alias live='curl parrot.live'
alias l='ls -CF'
alias la='ls -A'
alias ll='ls -alF'
alias ls='ls --color=auto'
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

sudo mkdir -p /etc/systemd/system/docker.service.d
sudo systemctl daemon-reload
sudo systemctl restart docker

# Check whether the Cgroup driver has been set correctly
# If the output was Cgroup Driver: systemd, all right!
docker info | grep -i cgroup

# Install Kubernetes with kubeadm
# References:
# https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
sudo swapoff -a
sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/k8s.list
cat /etc/apt/sources.list.d/k8s.list
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
kubectl version --client


#------- Specifics (master)
kubeadm config images pull
kubeadm init
# Reset
# kubeadm reset
# rm -rf /etc/cni/net.d

# Config kubectl
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
kubectl get nodes

# Config pod network (weave-net)
sudo modprobe br_netfilter ip_vs_rr ip_vs_wrr ip_vs_sh nf_conntrack_ipv4 ip_vs
kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
kubectl get pods -A
kubectl get nodes


#------- Specifics (worker1 and worker2)
# Allow all these ports: https://github.com/badtuxx/DescomplicandoKubernetes/blob/main/pt/day_one/descomplicando_kubernetes.md#portas-que-devemos-nos-preocupar
#
# In master node print a join command for add worker node in cluster
kubeadm token create --print-join-command
#
# Example of command to run in worker node:
# kubeadm join 10.0.35.25:6443 --token xzlzjw.bksen5z6somtgh22 --discovery-token-ca-cert-hash sha256:86e507a7af3de3b47aceff4c9a2466e965e72ff7236a37031ea76258425b5c72

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
kubectl apply -f samples/addons
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

kubectl apply -f $ISTIO_DIR_BASE/samples/addons
kubectl rollout status deployment/kiali -n istio-system

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