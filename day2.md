# Descomplicando Istio

<!-- TOC -->

- [Descomplicando Istio](#descomplicando-istio)
- [Day 2](#day-2)

<!-- TOC -->


# Day 2

Estratégias de gerenciamento de tráfego (*traffic management*):

> * **traffic request** => roteamento baseado nas configurações das requisições (requests), como por exemplo campos do header.
> 
> * **traffic shifting** => roteamento baseado no peso (valor ou percentual) das requisições. É útil para **canary deploy**.
> 
> * **fault injection** => roteamento baseado em cenários de falhas das requisições para testar a resiliência das aplicações e microsserviços.
> 
> * **request timeout** => roteamento baseado em cenários em que ocorre timeout das requisições para testar a resiliência das aplicações e microsserviços.
> 
> * **circuit breaking** => roteamento baseado em cenários em que ocorre problemas de conectividade ou falha de um ou mais microsserviços que formam uma ou mais aplicações.
> 
> * **mirroring** => espelhamento de tráfego, bem útil para testes de uma versão nova de um aplicação sem impactar a versão em uso no ambiente de produção. Necessário ter o ambiente de produção replicado (duplica o custo). Necessário ter o ambiente de produção replicado (duplica o custo).O Istio descarta as requisições espelhadas após chegarem nas aplicações para evitar duplicidade no processamento para o usuário final.


Os comandos a seguir foram executados apenas no **master**.

```bash
#------- Specifics (master)
sudo su
export ISTIO_DIR_BASE="/home/ubuntu/istio-1.11.4"
cd $ISTIO_DIR_BASE
export PATH="$PATH:$ISTIO_DIR_BASE/bin"

# Complementary files
export COMPLEMENTARY_FILES=/home/ubuntu/learning-istio/files

#----------------- Traffic request
vim $ISTIO_DIR_BASE/samples/bookinfo/networking/destination-rule-all.yaml

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/destination-rule-all.yaml
kubectl get destinationrules -o yaml

vim $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-all-v1.yaml

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-all-v1.yaml
kubectl get virtualservices

vim $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-reviews-test-v2.yaml

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-reviews-test-v2.yaml
kubectl get virtualservice reviews -o yaml
kubectl delete -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-all-v1.yaml

#----------------- Traffic shifting
kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-all-v1.yaml
kubectl get virtualservice reviews -o yaml

vim $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-reviews-50-v3.yaml
kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-reviews-50-v3.yaml

kubectl delete -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-all-v1.yaml

#----------------- Fault injection
kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-all-v1.yaml

vim $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-reviews-test-v2.yaml

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-reviews-test-v2.yaml

vim $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-ratings-test-delay.yaml

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-ratings-test-delay.yaml 
kubectl get virtualservice ratings -o yaml

vim $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-ratings-test-abort.yaml

kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-ratings-test-abort.yaml

kubectl delete -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-all-v1.yaml

#----------------- Request timeout
kubectl apply -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-all-v1.yaml

vim $COMPLEMENTARY_FILES/request-timeout/virtual-service-reviews-normal.yaml

kubectl apply -f $COMPLEMENTARY_FILES/request-timeout/virtual-service-normal.yaml

vim $COMPLEMENTARY_FILES/request-timeout/virtual-service-ratings-delay.yaml

kubectl apply -f $COMPLEMENTARY_FILES/request-timeout/virtual-service-ratings-delay.yaml

vim $COMPLEMENTARY_FILES/request-timeout/virtual-service-reviews-timeout.yaml

kubectl apply -f $COMPLEMENTARY_FILES/request-timeout/virtual-service-reviews-timeout.yaml
kubectl delete -f $ISTIO_DIR_BASE/samples/bookinfo/networking/virtual-service-all-v1.yaml

#----------------- Circuit breaking
vim $ISTIO_DIR_BASE/samples/httpbin/httpbin.yaml

kubectl apply -f $ISTIO_DIR_BASE/samples/httpbin/httpbin.yaml

vim $COMPLEMENTARY_FILES/circuit-braker/circuit-breaker.yaml

kubectl apply -f $COMPLEMENTARY_FILES/circuit-braker/circuit-breaker.yaml
kubectl get destinationrule httpbin -o yaml
kubectl apply -f $ISTIO_DIR_BASE/samples/httpbin/sample-client/fortio-deploy.yaml
FORTIO_POD=$(kubectl get pod | grep fortio | awk '{ print $1 }')

kubectl exec -it $FORTIO_POD  -c fortio /usr/bin/fortio -- load -curl  http://httpbin:8000/get

kubectl exec -it $FORTIO_POD  -c fortio /usr/bin/fortio -- load -c 2 -qps 0 -n 20 -loglevel Warning http://httpbin:8000/get

kubectl exec -it $FORTIO_POD  -c fortio /usr/bin/fortio -- load -c 3 -qps 0 -n 30 -loglevel Warning http://httpbin:8000/get

kubectl exec -it $FORTIO_POD  -c istio-proxy  -- sh -c 'curl localhost:15000/stats' | grep httpbin | grep pending

kubectl delete destinationrule httpbin
kubectl delete deploy httpbin fortio-deploy
kubectl delete svc httpbin

#----------------- Mirroring
vim $COMPLEMENTARY_FILES/mirroring/deployment-httpbin-v1.yaml
vim $COMPLEMENTARY_FILES/mirroring/deployment-httpbin-v2.yaml

kubectl apply -f $COMPLEMENTARY_FILES/mirroring/deployment-httpbin-v1.yaml
kubectl apply -f $COMPLEMENTARY_FILES/mirroring/deployment-httpbin-v2.yaml

vim $COMPLEMENTARY_FILES/mirroring/service-httpbin.yaml

kubectl apply -f $COMPLEMENTARY_FILES/mirroring/service-httpbin.yaml

vim $COMPLEMENTARY_FILES/mirroring/deployment-sleep.yaml

kubectl apply -f $COMPLEMENTARY_FILES/mirroring/deployment-sleep.yaml

vim $COMPLEMENTARY_FILES/mirroring/virtual-service-httpbin.yaml

kubectl apply -f $COMPLEMENTARY_FILES/mirroring/virtual-service-httpbin.yaml

vim $COMPLEMENTARY_FILES/mirroring/destination-rule-httpbin.yaml

kubectl apply -f $COMPLEMENTARY_FILES/mirroring/destination-rule-httpbin.yaml
kubectl get pods
kubectl exec -ti sleep-674f75ff4d-5bsnf -c sleep -- sh -c 'curl http://httpbin:8000/headers'
kubectl logs httpbin-v1-b9985cc7d-t5f4z  httpbin
kubectl logs httpbin-v2-5cdb74d4c7-2lc27 httpbin

vim $COMPLEMENTARY_FILES/mirroring/mirroring-httpbin.yaml

kubectl apply -f $COMPLEMENTARY_FILES/mirroring/mirroring-httpbin.yaml
kubectl get virtualservices.networking.istio.io

vim $COMPLEMENTARY_FILES/mirroring/mirroring-httpbin.yaml

kubectl get virtualservices.networking.istio.io httpbin -o yaml
kubectl exec -ti sleep-674f75ff4d-5bsnf -c sleep -- sh -c 'curl http://httpbin:8000/headers'
kubectl logs httpbin-v1-b9985cc7d-t5f4z  httpbin
kubectl logs httpbin-v2-5cdb74d4c7-2lc27 httpbin
kubectl exec -ti sleep-674f75ff4d-5bsnf -c sleep -- sh -c 'curl http://httpbin:8000/headers'
kubectl logs httpbin-v2-5cdb74d4c7-2lc27 httpbin
kubectl logs httpbin-v1-b9985cc7d-t5f4z  httpbin
```