# Descomplicando Istio

<!-- TOC -->

- [Descomplicando Istio](#descomplicando-istio)
- [Day 3](#day-3)

<!-- TOC -->


# Day 3

Comandos executados em sequencia durante o treinamento.

```bash
#------- Specifics (master)
export ISTIO_DIR_BASE="/home/ubuntu/istio-1.11.4"
cd $ISTIO_DIR_BASE
export PATH="$PATH:$ISTIO_DIR_BASE/bin"

# Complementary files
export COMPLEMENTARY_FILES=/home/ubuntu/learning-istio/files

#----------------- Authentication Policy
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
kubectl get pods --all-namespaces

kubectl exec -ti -n strigus sleep-7d457d69b5-h8frp -- curl http://httpbin.giropops:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl exec -ti -n giropopssleep-7d457d69b5-tw5pt -- curl http://httpbin.strigus:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl exec -ti -n giropops sleep-7d457d69b5-tw5pt -- curl http://httpbin.girus:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl get policies.authentication.istio.io --all-namespaces
kubectl get meshpolicies.authentication.istio.io --all-namespaces
kubectl get destinationrules.networking.istio.io --all-namespaces -o yaml | grep "host:"

vim $COMPLEMENTARY_FILES/authentication-policy/meshpolicy.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/meshpolicy.yaml

kubectl exec -ti -n strigus sleep-7d457d69b5-h8frp -- curl http://httpbin.giropops:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl exec -ti -n giropopssleep-7d457d69b5-tw5pt -- curl http://httpbin.strigus:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl exec -ti -n giropopssleep-7d457d69b5-tw5pt -- curl http://httpbin.girus:8000/ip -s -o /dev/null -w "%{http_code}\n"

vim $COMPLEMENTARY_FILES/authentication-policy/destination-rule-meshpolicy.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/destination-rule-meshpolicy.yaml

kubectl exec -ti -n strigus sleep-7d457d69b5-h8frp -- curl http://httpbin.giropops:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl exec -ti -n giropopssleep-7d457d69b5-tw5pt -- curl http://httpbin.strigus:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl exec -ti -n giropopssleep-7d457d69b5-tw5pt -- curl http://httpbin.girus:8000/ip -s -o /dev/null -w "%{http_code}\n"

vim $COMPLEMENTARY_FILES/authentication-policy/destination-rule-fix-girus-conn.yaml


kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/destination-rule-fix-girus-conn.yaml

kubectl delete meshpolicy default
kubectl delete destinationrules httpbin-girus -n girus
kubectl delete destinationrules api-server -n istio-system
kubectl delete destinationrules default -n istio-system

vim $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-namespace.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-namespace.yaml

for from in "giropops" "strigus"; do for to in "giropops" "strigus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

vim $COMPLEMENTARY_FILES/authentication-policy/destination-rule-enable-mtls-namespace.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/destination-rule-enable-mtls-namespace.yaml

for from in "giropops" "strigus"; do for to in "giropops" "strigus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

vim $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-service.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-service.yaml

for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

vim $COMPLEMENTARY_FILES/authentication-policy/destination-rule-enable-mtls-service.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/destination-rule-enable-mtls-service.yaml

for from in "giropops" "strigus" "girus"; do for to in "giropops" "strigus" "girus"; do kubectl exec $(kubectl get pod -l app=sleep -n ${from} -o jsonpath={.items..metadata.name}) -c sleep -n ${from} -- curl "http://httpbin.${to}:8000/ip" -s -o /dev/null -w "sleep.${from} to httpbin.${to}: %{http_code}\n"; done; done

vim $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-specific-service.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/policy-enable-mtls-specific-service.yaml

vim $COMPLEMENTARY_FILES/authentication-policy/destination-rule-specific-service.yaml

kubectl apply -f $COMPLEMENTARY_FILES/authentication-policy/destination-rule-specific-service.yaml

kubectl exec -ti -n girus sleep-7d457d69b5-sr5xw -- curl http://httpbin.strigus:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl exec -ti -n girus sleep-7d457d69b5-sr5xw -- curl http://httpbin.giropops:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl exec $(kubectl get pod -l app=sleep -n girus -o jsonpath={.items..metadata.name}) -c sleep -n girus -- curl http://httpbin.strigus:8000/ip -s -o /dev/null -w "%{http_code}\n"

kubectl delete policy default -n giropops
kubectl delete policy httpbin -n strigus
kubectl delete destinationrules default -n giropops
kubectl delete destinationrules httpbin -n strigus

#----------------- Authorization
```