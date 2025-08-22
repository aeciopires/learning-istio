<!-- TOC -->

- [uninstall-istio](#uninstall-istio)

<!-- TOC -->

# uninstall-istio

- Para desinstalar o Istio siga as instruções da página: https://istio.io/v1.27/docs/ambient/install/helm/#uninstall

- Para remover o cluster Kind, execute o seguinte comando:

```bash
kind delete clusters $(kind get clusters)
```

- Para remover variáveis de ambiente usadas na instalação, execute os seguintes comandos:

```bash
unset ISTIO_RELEASE
unset VERSION_ISTIO
unset ISTIO_BASE_URL
unset ISTIO_BOOKINFO_URL
unset ISTIO_ADDONS_URL
unset MY_NAMESPACE
unset ISTIO_HTTPBIN_URL
```
