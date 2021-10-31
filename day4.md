# Descomplicando Istio

<!-- TOC -->

- [Descomplicando Istio](#descomplicando-istio)
- [Day 4](#day-4)

<!-- TOC -->


# Day 4

Comandos executados em sequencia durante o treinamento.

```bash
#------- Specifics (master)
export ISTIO_DIR_BASE="/home/ubuntu/istio-1.11.4"
cd $ISTIO_DIR_BASE
export PATH="$PATH:$ISTIO_DIR_BASE/bin"

# Complementary files
export COMPLEMENTARY_FILES=/home/ubuntu/learning-istio/files

#----------------- Uninstall Istio
kubectl delete -f $ISTIO_DIR_BASE/install/kubernetes/istio-demo-auth.yaml
kubectl delete -f $ISTIO_DIR_BASE/install/kubernetes/istio-demo.yaml
for i in $ISTIO_DIR_BASE/install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl delete -f $i; done
kubectl get virtualservices.networking.istio.io --all-namespaces

#----------------- Install Helm
# Reference: https://helm.sh/docs/intro/install/
cd $HOME
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 $HOME/get_helm.sh
$HOME/get_helm.sh

#----------------- Install Istio with Helm
# Reference: https://istio.io/latest/docs/setup/install/helm/
for i in $ISTIO_DIR_BASE/install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl apply -f $i; done

# helm template install/kubernetes/helm/istio --name istio --namespace istio-system --values install/kubernetes/helm/istio/values.yaml --set gateways.istio-ingresssgateway.type=NodePort --set grafana.enabled=true --set kiali.enabled=true --set tracing.enabled=true --set kiali.dashboard.username=admin --set kiali.dashboard.passphrase=admin --set servicegraph.enabled=true > $HOME/meu_istio.yaml

vim $HOME/meu_istio.yaml
kubectl create namespace istio-system
kubectl apply -f $HOME/meu_istio.yaml
kubectl get pods -n istio-system
kubectl get pods -n istio-system --watch

 # KIALI_USERNAME=$(read -p 'Kiali Username: ' uval && echo -n $uval | base64)
 # KIALI_PASSPHRASE=$(read -sp 'Kiali Passphrase: ' pval && echo -n $pval | base64)
 # echo $KIALI_USERNAME
 # echo $KIALI_PASSPHRASE
 # NAMESPACE=istio-system
 # cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: kiali
  namespace: $NAMESPACE
  labels:
    app: kiali
type: Opaque
data:
  username: $KIALI_USERNAME
  passphrase: $KIALI_PASSPHRASE
EOF

kubectl get services -n istio-system
kubectl port-forward svc/kiali 20001:20001 -n istio-system --address=0.0.0.0 &


#----------------- Egress
kubectl apply -f $ISTIO_DIR_BASE/samples/sleep/sleep.yaml
kubectl get deployments.
kubectl get pods
kubectl get configmaps istio -n istio-system
kubectl get configmaps istio -n istio-system -o yaml | grep -o "mode: ALLOW_ANY"

kubectl exec -ti sleep-7d457d69b5-wbcbw -c sleep -- curl -I https://www.linuxtips.io

kubectl edit configmaps istio -n istio-system

kubectl get configmaps istio -n istio-system -o yaml | sed 's/mode: ALLOW_ANY/mode: REGISTRY_ONLY/g' | kubectl replace -n istio-system -f -

kubectl get configmaps istio -n istio-system -o yaml | sed 's/mode: REGISTRY_ONLY/mode: ALLOW_ANY/g' | kubectl replace -n istio-system -f -

kubectl apply -f $COMPLEMENTARY_FILES/egress/libera_httpbin_https.yaml
kubectl apply -f $COMPLEMENTARY_FILES/egress/libera_httpbin.yaml
kubectl apply -f $COMPLEMENTARY_FILES/egress/vs_httpbin_org.yaml
kubectl exec -ti sleep-7d457d69b5-wbcbw -c sleep -- curl -I http://httpbin.org/delay/2
kubectl exec -ti sleep-7d457d69b5-wbcbw -c sleep -- curl -I http://httpbin.org/delay/4
kubectl exec -ti sleep-7d457d69b5-wbcbw -c sleep -- curl -I http://httpbin.org/delay/3

kubectl delete -f $COMPLEMENTARY_FILES/egress

#----------------- Retry policy
vim $COMPLEMENTARY_FILES/retry-policy/destination_rule.yaml

kubectl apply -f $COMPLEMENTARY_FILES/retry-policy/destination_rule.yaml

vim $COMPLEMENTARY_FILES/retry-policy/virtualservice.yaml

kubectl apply -f $COMPLEMENTARY_FILES/retry-policy/virtualservice.yaml

kubectl delete -f $COMPLEMENTARY_FILES/retry-policy/
```