# Descomplicando Istio

<!-- TOC -->

- [Descomplicando Istio](#descomplicando-istio)
- [Day 4](#day-4)

<!-- TOC -->


# Day 4

Os comandos a seguir foram executados apenas no **master**.

```bash
# Loop for test the access in productpage application
# Run in your compute
for aux in $(seq 1 1000); do curl -I -X GET http://master:9080/productpage; done

# Loop for test the access in httpbin application
# Run in your compute
for aux in $(seq 1 1000); do curl -I -X GET http://master:8000; done

#------- Specifics (master)
sudo su
export ISTIO_DIR_BASE="/home/ubuntu/istio-1.11.4"
cd $ISTIO_DIR_BASE
export PATH="$PATH:$ISTIO_DIR_BASE/bin"

# Complementary files
export COMPLEMENTARY_FILES=/home/ubuntu/learning-istio/files

#----------------- Uninstall Istio
# Uninstall App and Istio
kubectl delete -f $ISTIO_DIR_BASE/samples/addons
istioctl manifest generate --set profile=demo | kubectl delete --ignore-not-found=true -f -
kubectl delete namespace istio-system
kubectl get crd | grep --color=never 'istio.io' | awk '{print $1}' | xargs -n1 kubectl delete crd
kubectl label namespace default istio-injection-

#----------------- Install Helm
# Reference: https://helm.sh/docs/intro/install/
cd $HOME
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 $HOME/get_helm.sh
$HOME/get_helm.sh
helm version

#----------------- Install Istio with Helm
# Reference: https://istio.io/latest/docs/setup/install/helm/
kubectl create namespace istio-system
helm install istio-base $ISTIO_DIR_BASE/manifests/charts/base -n istio-system

helm install istiod $ISTIO_DIR_BASE/manifests/charts/istio-control/istio-discovery \
    -n istio-system

helm install istio-ingress $ISTIO_DIR_BASE/manifests/charts/gateways/istio-ingress \
    -n istio-system

helm install istio-egress $ISTIO_DIR_BASE/manifests/charts/gateways/istio-egress \
    -n istio-system

kubectl get pods -n istio-system
kubectl label namespace default istio-injection=enabled

#----------------- Egress
kubectl apply -f $ISTIO_DIR_BASE/samples/sleep/sleep.yaml
kubectl get deployments
kubectl get pods
kubectl get configmaps istio -n istio-system

SLEEP_POD=$(kubectl get pod | grep sleep | awk '{ print $1 }' | head -n1)
kubectl exec -ti $SLEEP_POD -c sleep -- curl -I https://www.linuxtips.io

kubectl edit configmaps istio -n istio-system

kubectl apply -f $COMPLEMENTARY_FILES/egress/libera_httpbin_https.yaml
kubectl apply -f $COMPLEMENTARY_FILES/egress/libera_httpbin.yaml
kubectl apply -f $COMPLEMENTARY_FILES/egress/vs_httpbin_org.yaml

kubectl exec -ti $SLEEP_POD -c sleep -- curl -I http://httpbin.org/delay/2

kubectl exec -ti $SLEEP_POD -c sleep -- curl -I http://httpbin.org/delay/4

kubectl exec -ti $SLEEP_POD -c sleep -- curl -I http://httpbin.org/delay/3

# Remove configurations
kubectl delete -f $COMPLEMENTARY_FILES/egress

#----------------- Retry policy
vim $COMPLEMENTARY_FILES/retry-policy/destination_rule.yaml

kubectl apply -f $COMPLEMENTARY_FILES/retry-policy/destination_rule.yaml

vim $COMPLEMENTARY_FILES/retry-policy/virtualservice.yaml

kubectl apply -f $COMPLEMENTARY_FILES/retry-policy/virtualservice.yaml

# Remove configurations
kubectl delete -f $COMPLEMENTARY_FILES/retry-policy/

# Uninstall Istio with helm
helm delete istio-egress -n istio-system
helm delete istio-ingress -n istio-system
helm delete istiod -n istio-system
helm delete istio-base -n istio-system
kubectl delete namespace istio-system
kubectl get crd | grep --color=never 'istio.io' | awk '{print $1}' | xargs -n1 kubectl delete crd

# Remove configurations
kubectl delete ns giropops girus strigus

kubectl delete deployment details-v1 httpbin-v1 httpbin-v2 productpage-v1 ratings-v1 reviews-v1 reviews-v2 reviews-v3 sleep

kubectl delete svc details fortio httpbin productpage ratings reviews sleep

kubectl get all -A
```