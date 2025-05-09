# Descomplicando Istio

<!-- TOC -->

- [Descomplicando Istio](#descomplicando-istio)
- [Day 3](#day-3)

<!-- TOC -->

# Day 3

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

#----------------- Authentication Policy
kubectl get all --all-namespaces
kubectl create ns giropops
kubectl create ns strigus
kubectl create ns girus
kubectl label namespace strigus istio-injection=enabled
kubectl label namespace giropops istio-injection=enabled
cd $ISTIO_DIR_BASE

kubectl apply -f $ISTIO_DIR_BASE/samples/httpbin/httpbin.yaml -n giropops
kubectl apply -f $ISTIO_DIR_BASE/samples/httpbin/httpbin.yaml -n strigus
kubectl apply -f $ISTIO_DIR_BASE/samples/httpbin/httpbin.yaml -n girus
kubectl apply -f $ISTIO_DIR_BASE/samples/sleep/sleep.yaml -n giropops
kubectl apply -f $ISTIO_DIR_BASE/samples/sleep/sleep.yaml -n strigus
kubectl apply -f $ISTIO_DIR_BASE/samples/sleep/sleep.yaml -n girus
kubectl get deployments --all-namespaces
kubectl get pods --all-namespaces

# Testing the comunication between pods of namespaces
SLEEP_POD_STRIGUS=$(kubectl get pod -n strigus | grep sleep | awk '{ print $1 }' | head -n1)

kubectl exec -ti -n strigus $SLEEP_POD_STRIGUS -- curl http://httpbin.giropops:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl exec -ti -n strigus $SLEEP_POD_STRIGUS -- curl http://httpbin.girus:8000/ip -s -o /dev/null -w "%{http_code}\n"

SLEEP_POD_GIROPOPS=$(kubectl get pod -n giropops | grep sleep | awk '{ print $1 }' | head -n1)

kubectl exec -ti -n giropops $SLEEP_POD_GIROPOPS -- curl http://httpbin.strigus:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl exec -ti -n giropops $SLEEP_POD_GIROPOPS -- curl http://httpbin.girus:8000/ip -s -o /dev/null -w "%{http_code}\n"

SLEEP_POD_GIRUS=$(kubectl get pod -n girus | grep sleep | awk '{ print $1 }' | head -n1)

kubectl exec -ti -n girus $SLEEP_POD_GIRUS -- curl http://httpbin.strigus:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl exec -ti -n girus $SLEEP_POD_GIRUS -- curl http://httpbin.giropops:8000/ip -s -o /dev/null -w "%{http_code}\n"

# Enable MTLS in cluster
vim $COMPLEMENTARY_FILES/authentication-policy/meshpolicy.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/meshpolicy.yaml

kubectl get authorizationpolicies.security.istio.io --all-namespaces
kubectl get peerauthentications.security.istio.io --all-namespaces
kubectl get requestauthentications.security.istio.io --all-namespaces
kubectl get destinationrules.networking.istio.io --all-namespaces -o yaml | grep "host:"

# Testing the comunication between pods of namespaces
for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

# Enable MTLS local services of cluster
vim $COMPLEMENTARY_FILES/authentication-policy/destination-rule-meshpolicy.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/destination-rule-meshpolicy.yaml

# Testing the comunication between pods of namespaces
for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

# Disable MTLS for services of girus namespace
vim $COMPLEMENTARY_FILES/authentication-policy/destination-rule-fix-girus-conn.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/destination-rule-fix-girus-conn.yaml

# Testing the comunication between pods of namespaces
for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

# Remove configurations
kubectl delete PeerAuthentication default
kubectl delete destinationrules httpbin-girus -n girus
kubectl delete destinationrules api-server -n istio-system
kubectl delete destinationrules default -n istio-system

# Enable MTLS for giropops namespace
vim $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-namespace.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-namespace.yaml

# Testing the comunication between pods of namespaces
for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

# Enable MTLS for services of giropops namespace
vim $COMPLEMENTARY_FILES/authentication-policy/destination-rule-enable-mtls-namespace.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/destination-rule-enable-mtls-namespace.yaml

# Testing the comunication between pods of namespaces
for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

# Enable MTLS for strigus namespace
vim $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-service.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-service.yaml

# Testing the comunication between pods of namespaces
for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

# Enable MTLS for services of strigus namespace
vim $COMPLEMENTARY_FILES/authentication-policy/destination-rule-enable-mtls-service.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/destination-rule-enable-mtls-service.yaml

# Testing the comunication between pods of namespaces
for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

# Enable MTLS only for port of specif service of strigus namespace
vim $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-specific-service.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-specific-service.yaml

# Testing the comunication between pods of namespaces
for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

# Disable MTLS only for port of specif service of strigus namespace
vim $COMPLEMENTARY_FILES/authentication-policy/destination-rule-specific-service.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/destination-rule-specific-service.yaml

# Testing the comunication between pods of namespaces
for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

kubectl delete peerauthentication default -n giropops
kubectl delete peerauthentication httpbin -n strigus
kubectl delete destinationrules default -n giropops
kubectl delete destinationrules httpbin -n strigus

#----------------- Authorization
kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/meshpolicy.yaml
kubectl get peerauthentication.security.istio.io --all-namespaces
kubectl get peerauthentication

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/platform/kube/bookinfo.yaml
kubectl get deployments --all-namespaces

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/destination-rule-all-mtls.yaml
kubectl get destinationrules -o yaml

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/bookinfo-gateway.yaml
kubectl get deployments --all-namespaces
kubectl get svc

# Port forward in background for access Productpage application
# Allow port 9080/TCP
kubectl port-forward svc/productpage 9080:9080 -n default --address=0.0.0.0 > /dev/null 2>&1 &
# Access URL: http://master:9080/productpage

# Enable authorization in namespace
vim $COMPLEMENTARY_FILES/authorization/rbac/rbac-config-ON.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authorization/rbac/rbac-config-ON.yaml

vim $COMPLEMENTARY_FILES/authorization/rbac/namespace-policy.yaml

# Authorize get and head HTTP methods in default namespace
kubectl apply -f $COMPLEMENTARY_FILES/authorization/rbac/namespace-policy.yaml
kubectl get authorizationpolicy
kubectl get authorizationpolicy.security.istio.io

vim $COMPLEMENTARY_FILES/authorization/rbac/productpage-policy.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authorization/rbac/productpage-policy.yaml

vim $COMPLEMENTARY_FILES/authorization/rbac/ratings-policy.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authorization/rbac/ratings-policy.yaml

kubectl get authorizationpolicy

# Remove configurations
kubectl delete -f $COMPLEMENTARY_FILES/authorization/rbac/
```
