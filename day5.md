# Descomplicando Istio

<!-- TOC -->

- [Descomplicando Istio](#descomplicando-istio)
- [Day 5](#day-5)
- [Antes de usar um EnvoyFilter](#antes-de-usar-um-envoyfilter)
  - [Sidecar x ambient: onde o EnvoyFilter funciona](#sidecar-x-ambient-onde-o-envoyfilter-funciona)
- [Campos do EnvoyFilter além do básico](#campos-do-envoyfilter-além-do-básico)
- [Quais filtros do Envoy existem no proxy do Istio](#quais-filtros-do-envoy-existem-no-proxy-do-istio)
- [Laboratório 1 (iniciante): compressão gzip no gateway de entrada](#laboratório-1-iniciante-compressão-gzip-no-gateway-de-entrada)
  - [Quando o Envoy comprime (e quando não comprime)](#quando-o-envoy-comprime-e-quando-não-comprime)
- [Laboratório 2 (intermediário): brotli, zstd e mais de um algoritmo](#laboratório-2-intermediário-brotli-zstd-e-mais-de-um-algoritmo)
- [Laboratório 3 (intermediário): compressão no sidecar de um workload](#laboratório-3-intermediário-compressão-no-sidecar-de-um-workload)
- [Laboratório 4 (intermediário): filtro Lua adicionando um header de resposta](#laboratório-4-intermediário-filtro-lua-adicionando-um-header-de-resposta)
- [Laboratório 5 (avançado): rate limit local](#laboratório-5-avançado-rate-limit-local)
- [Laboratório 6 (avançado): ajustando o HTTP Connection Manager com MERGE](#laboratório-6-avançado-ajustando-o-http-connection-manager-com-merge)
- [Modo ambient: estendendo o waypoint com TrafficExtension](#modo-ambient-estendendo-o-waypoint-com-trafficextension)
- [Cuidados em upgrades (Istio 1.32)](#cuidados-em-upgrades-istio-132)
- [Troubleshooting](#troubleshooting)
- [Limpeza](#limpeza)
- [Referências](#referências)

<!-- TOC -->

# Day 5

Neste dia vamos customizar o Envoy por baixo do Istio com o recurso **`EnvoyFilter`**: conhecer os **filtros do Envoy** disponíveis no proxy do Istio e praticar com **compressão** (gzip, brotli e zstd), **Lua**, **rate limit local** e ajustes no **HTTP Connection Manager**.

> **Nível:** intermediário/avançado. Pré-requisito: ambiente do [Day 1](day1.md) ou [Day 1 - sidecar](day1-sidecar.md) com o Bookinfo exposto pelo gateway (`bookinfo-gateway`) e as variáveis `MY_NAMESPACE` e `GATEWAY_IP` definidas.
>
> A teoria está no Guia de Estudos: [3.7 - a cadeia de filtros HTTP](ISTIO_ENVOY_GUIA_DE_ESTUDOS.md#37-a-cadeia-de-filtros-http-em-detalhe-tipos-ordem-e-controle-de-fluxo) e [7 - EnvoyFilter](ISTIO_ENVOY_GUIA_DE_ESTUDOS.md#7-envoyfilter) (`workloadSelector`, `applyTo`, `match`, `patch.operation`). Este arquivo é a parte **prática** e complementa o guia com campos e exemplos que não estão lá.

# Antes de usar um EnvoyFilter

A [referência oficial](https://istio.io/v1.32/docs/reference/config/networking/envoy-filter/) descreve o `EnvoyFilter` como um mecanismo para customizar a configuração do Envoy gerada pelo istiod, que **"deve ser usado com cuidado, pois configurações incorretas podem desestabilizar a mesh inteira"**. Três avisos da própria referência:

1. Partes dessa API são profundamente ligadas à implementação interna do Istio e à API xDS do Envoy. A API `EnvoyFilter` mantém compatibilidade, mas a configuração de Envoy que você fornece por ela **precisa ser acompanhada a cada upgrade** do proxy do Istio (campos obsoletos precisam ser removidos/substituídos).
2. Quando vários `EnvoyFilter` se aplicam ao mesmo workload em um namespace, os patches são processados **em ordem de criação**, e o comportamento é **indefinido** se eles conflitarem.
3. Um `EnvoyFilter` no **namespace raiz** (*config root namespace*, ex.: `istio-system`) **sem** `workloadSelector` vale para **todos** os workloads (sidecars e gateways) da mesh.

Por isso, suba esta escada só até onde for necessário:

| Necessidade | Use primeiro |
|---|---|
| Roteamento, retries, timeouts, fault injection, mirroring, CORS | `VirtualService` / `HTTPRoute` ([Day 2](day2.md)) |
| Load balancing, circuit breaker, outlier detection, TLS para o destino | `DestinationRule` |
| Métricas, access logs e tracing | `Telemetry` ([Day 4](day4.md)) |
| Autenticação e autorização | `PeerAuthentication`, `RequestAuthentication`, `AuthorizationPolicy` |
| Lógica customizada (Wasm ou Lua) | `WasmPlugin` ou `TrafficExtension` |
| **Nada acima atende** (ex.: habilitar um filtro nativo do Envoy que o Istio não expõe, como o de compressão) | `EnvoyFilter` |

## Sidecar x ambient: onde o EnvoyFilter funciona

| Proxy | É Envoy? | `EnvoyFilter` |
|---|---|---|
| Sidecar (`istio-proxy`) | Sim | Suportado (`context: SIDECAR_INBOUND` / `SIDECAR_OUTBOUND`) |
| Gateways de entrada/saída | Sim | Suportado (`context: GATEWAY`) |
| **Waypoint** (modo ambient) | Sim | **Não suportado** |
| **ztunnel** (modo ambient) | Não - é um proxy próprio, escrito em Rust, restrito a L3/L4 | Não se aplica |

Sobre o waypoint, a documentação oficial ([Ambient / L7 features](https://istio.io/v1.32/docs/ambient/usage/l7-features/)) diz que o `EnvoyFilter` **não é suportado com waypoints em nenhuma versão existente do Istio**: pode até funcionar em cenários limitados, mas o uso não é suportado e é ativamente desencorajado pelos mantenedores. O guia de [migração para o ambient](https://istio.io/v1.32/docs/ambient/migrate/) repete: configurações de `EnvoyFilter` dos sidecars **não podem ser levadas** para os waypoints.

Na prática:

- Os laboratórios 1 e 2 usam o **gateway de entrada** do Bookinfo, que existe nos dois modos - funcionam tanto no ambiente do Day 1 (ambient) quanto no do Day 1 - sidecar.
- Os laboratórios 3 a 5 usam `SIDECAR_INBOUND` - **exigem o modo sidecar**.
- Para estender o L7 no modo ambient, veja [TrafficExtension](#modo-ambient-estendendo-o-waypoint-com-trafficextension).

# Campos do EnvoyFilter além do básico

Além de `workloadSelector` e `configPatches` (explicados na seção 7 do Guia de Estudos), a referência documenta:

- **`targetRefs`**: alternativa ao `workloadSelector` (só um dos dois pode ser usado) para anexar o `EnvoyFilter` a um recurso. Tipos suportados: `Gateway` (grupo `gateway.networking.k8s.io`, no mesmo namespace), `GatewayClass` (no namespace raiz), `Service` (somente para waypoints) e `ServiceEntry` (grupo `networking.istio.io`). É a forma natural de mirar um gateway criado pelo Kubernetes Gateway API, como o `bookinfo-gateway` deste repositório.
- **`priority`**: ordem de aplicação entre conjuntos de patches. O padrão é `0`; valores negativos são processados antes do padrão e positivos depois. A ordenação é por prioridade, depois data de criação, depois nome totalmente qualificado. A referência recomenda usar múltiplos de 10 para deixar espaço para inserções futuras.
- **`patch.filterClass`** (usado com `operation: ADD`): define o ponto de inserção em relação aos filtros que o próprio Istio insere - `AUTHN` (depois dos filtros de autenticação do Istio), `AUTHZ` (depois dos de autorização) e `STATS` (antes dos de estatísticas). A referência diz que este é o mecanismo **preferido** para adicionar filtros, em vez das operações `INSERT_*`, porque estas dependem de nomes de filtros potencialmente instáveis.
- **`match.listener.filterChain.transportSocket`** (**novo no Istio 1.32**): casa pelo nome do *transport socket* da filter chain. Como uma filter chain que casa `tls` não necessariamente **termina** TLS, este campo permite mirar com precisão as chains que terminam TLS.
- **`match.context`**: `ANY`, `SIDECAR_INBOUND`, `SIDECAR_OUTBOUND` ou `GATEWAY`.

# Quais filtros do Envoy existem no proxy do Istio

O `istio-proxy` é um build próprio do Envoy: **nem toda extensão do Envoy upstream está compilada nele**. A lista do que entra no build fica no arquivo [`bazel/extension_config/extensions_build_config.bzl`](https://github.com/istio/proxy/blob/master/bazel/extension_config/extensions_build_config.bzl) do repositório `istio/proxy`. Para saber o que existe **no seu proxy**, pergunte a ele:

```bash
export GW_POD=$(kubectl -n $MY_NAMESPACE get pod -l gateway.networking.k8s.io/gateway-name=bookinfo-gateway -o jsonpath='{.items[0].metadata.name}')

# Versão do Envoy empacotada
kubectl -n $MY_NAMESPACE exec $GW_POD -- pilot-agent request GET server_info | grep '"version"'

# Filtros HTTP REALMENTE configurados neste proxy, na ordem de execução
istioctl proxy-config listener $GW_POD.$MY_NAMESPACE -o json | grep '"name": "envoy.filters.http' | sort | uniq -c
```

Filtros HTTP presentes no build do proxy usado pelo Istio 1.32 que mais aparecem em `EnvoyFilter` (conferidos no arquivo acima, no commit do `istio/proxy` referenciado pelo `istio.deps` da tag `1.32.0`):

| Filtro (`name`) | Para que serve | Documentação do Envoy |
|---|---|---|
| `envoy.filters.http.compressor` | Comprime respostas (e opcionalmente requisições). Bibliotecas: `gzip`, `brotli`, `zstd` | [Compressor](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/compressor_filter) |
| `envoy.filters.http.decompressor` | Descomprime requisições/respostas | [Decompressor](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/decompressor_filter) |
| `envoy.filters.http.lua` | Executa scripts Lua na requisição e na resposta | [Lua](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/lua_filter) |
| `envoy.filters.http.local_ratelimit` | Rate limit local (token bucket no próprio proxy) | [Local rate limit](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/local_rate_limit_filter) |
| `envoy.filters.http.ratelimit` | Rate limit global (consulta um serviço externo) | [Rate limit](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/rate_limit_filter) |
| `envoy.filters.http.ext_authz` | Autorização delegada a um serviço externo | [External Authorization](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter) |
| `envoy.filters.http.ext_proc` | Processamento externo da requisição/resposta | [External Processing](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_proc_filter) |
| `envoy.filters.http.wasm` | Executa extensões WebAssembly | [Wasm](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/wasm_filter) |
| `envoy.filters.http.grpc_json_transcoder` | Expõe um serviço gRPC como API JSON/REST | [gRPC-JSON transcoder](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/grpc_json_transcoder_filter) |
| `envoy.filters.http.grpc_web` | Ponte gRPC-Web | [gRPC-Web](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/grpc_web_filter) |
| `envoy.filters.http.header_mutation` | Altera headers de requisição/resposta | [Header mutation](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/header_mutation_filter) |
| `envoy.filters.http.buffer` | Faz buffer da requisição inteira antes de encaminhar | [Buffer](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/buffer_filter) |
| `envoy.filters.http.oauth2` | Fluxo OAuth2 no proxy | [OAuth2](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/oauth2_filter) |
| `envoy.filters.http.tap` | Captura requisições/respostas para depuração | [Tap](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/tap_filter) |
| `envoy.filters.http.cors`, `fault`, `jwt_authn`, `rbac` | CORS, injeção de falhas, validação de JWT e autorização | Já são configurados pelo próprio Istio (`VirtualService`, `RequestAuthentication`, `AuthorizationPolicy`) - **prefira essas APIs** |
| `envoy.filters.http.router` | Encaminha a requisição ao cluster de destino; é sempre o **último** da cadeia | [Router](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/router_filter) |

> A lista completa de filtros HTTP do Envoy upstream está em [HTTP filters](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/http_filters). A documentação `latest` do Envoy descreve a versão em desenvolvimento: um campo documentado lá pode ainda não existir na versão do Envoy empacotada pelo seu Istio. Na dúvida, o `istiod`/Envoy rejeita a configuração - veja [Troubleshooting](#troubleshooting).

# Laboratório 1 (iniciante): compressão gzip no gateway de entrada

Objetivo: fazer o gateway de entrada do Bookinfo comprimir as respostas com gzip quando o cliente pedir.

Antes, confirme que **hoje** a resposta não vem comprimida:

```bash
curl -s -o /dev/null -D - -H "Accept-Encoding: gzip" http://${GATEWAY_IP}/productpage | grep -i -E 'content-encoding|content-length|transfer-encoding|vary'
```

Aplique o `EnvoyFilter`. Ele mira o gateway com `targetRefs` e insere o filtro `compressor` imediatamente antes do `router`:

```bash
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: gateway-gzip
  namespace: $MY_NAMESPACE
spec:
  targetRefs:
    - kind: Gateway
      group: gateway.networking.k8s.io
      name: bookinfo-gateway
  configPatches:
    - applyTo: HTTP_FILTER
      match:
        context: GATEWAY
        listener:
          filterChain:
            filter:
              name: "envoy.filters.network.http_connection_manager"
              subFilter:
                name: "envoy.filters.http.router"
      patch:
        operation: INSERT_BEFORE
        value:
          name: envoy.filters.http.compressor
          typed_config:
            "@type": type.googleapis.com/envoy.extensions.filters.http.compressor.v3.Compressor
            response_direction_config:
              common_config:
                min_content_length: 100
                content_type:
                  - text/html
                  - application/json
            compressor_library:
              name: text_optimized
              typed_config:
                "@type": type.googleapis.com/envoy.extensions.compression.gzip.compressor.v3.Gzip
                memory_level: 3
                window_bits: 10
                compression_level: BEST_COMPRESSION
                compression_strategy: DEFAULT_STRATEGY
EOF
```

> O bloco `value` é o exemplo de configuração da [documentação oficial do filtro Compressor](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/compressor_filter) do Envoy.

Teste de novo:

```bash
curl -s -o /dev/null -D - -H "Accept-Encoding: gzip" http://${GATEWAY_IP}/productpage | grep -i -E 'content-encoding|content-length|transfer-encoding|vary'

# Compare o tamanho transferido com e sem compressão
curl -s -o /dev/null -w 'sem gzip: %{size_download} bytes\n' http://${GATEWAY_IP}/productpage
curl -s -o /dev/null -w 'com gzip: %{size_download} bytes\n' -H "Accept-Encoding: gzip" http://${GATEWAY_IP}/productpage
```

Resultado esperado, conforme a documentação do filtro - quando a compressão é aplicada:

- O header `content-length` é **removido** da resposta.
- A resposta passa a ter `transfer-encoding: chunked` e `content-encoding: gzip`.
- O header `vary: accept-encoding` é inserido em toda resposta.

## Quando o Envoy comprime (e quando não comprime)

A compressão de **resposta** vem habilitada por padrão no filtro, mas é **pulada** quando (lista da documentação oficial):

- A requisição **não** tem o header `accept-encoding`, ou ele não contém `gzip` nem `*`, ou contém com peso `q=0`.
- O `accept-encoding` inclui outro tipo de codificação com peso **maior** que o do `gzip`, e o filtro de compressão correspondente está presente na cadeia.
- A resposta já tem um header `content-encoding`.
- A resposta tem `cache-control` com `no-transform`.
- A resposta tem `transfer-encoding` com um nome de compressão conhecido.
- O `content-type` da resposta não está na lista de mime-types selecionados. O padrão é: `application/javascript`, `application/json`, `application/xhtml+xml`, `image/svg+xml`, `text/css`, `text/html`, `text/plain`, `text/xml`.
- A resposta não tem `content-length` nem `transfer-encoding`.
- A resposta é menor que 30 bytes (aplicável quando o `transfer-encoding` não é chunked).
- O código de resposta está na lista de códigos não comprimíveis (vazia por padrão).

Outros pontos documentados:

- **ETag**: por padrão, ETags fracas (`W/...`) são preservadas e ETags fortes são **removidas** da resposta comprimida (o corpo mudou, a ETag não corresponde mais). Há opções para pular a compressão quando houver ETag (`disable_on_etag_header`) - consulte a seção *ETag handling* da documentação para a versão do Envoy em uso.
- A compressão de **requisição** vem **desabilitada** por padrão.
- **Custo**: a própria documentação resume o trade-off - compressão é útil quando a banda é escassa e o payload comprime bem, **ao custo de mais CPU** no proxy.
- **Estatísticas**: cada filtro Compressor gera estatísticas em `<stat_prefix>.compressor.<compressor_library.name>.<compressor_library_stat_prefix>.<direction_prefix>.*` (`compressed`, `not_compressed`, `total_uncompressed_bytes`, `total_compressed_bytes`, `no_accept_header`, `header_compressor_used`, entre outras).

> Por padrão o Istio configura o Envoy para registrar um **conjunto mínimo** de estatísticas ([Envoy Statistics](https://istio.io/v1.32/docs/ops/configuration/telemetry/envoy-stats/)). Para ver as estatísticas do compressor é preciso incluí-las com `proxyStatsMatcher` (annotation `proxy.istio.io/config` no Pod), do mesmo jeito mostrado no [Laboratório 5](#laboratório-5-avançado-rate-limit-local), usando a regex `.*compressor.*`.

# Laboratório 2 (intermediário): brotli, zstd e mais de um algoritmo

O filtro suporta três bibliotecas de compressão. Para trocar o algoritmo, mude **somente** o `compressor_library`:

| Algoritmo | `Accept-Encoding` | `@type` do `compressor_library.typed_config` |
|---|---|---|
| gzip | `gzip` | `type.googleapis.com/envoy.extensions.compression.gzip.compressor.v3.Gzip` |
| Brotli | `br` | `type.googleapis.com/envoy.extensions.compression.brotli.compressor.v3.Brotli` |
| Zstandard | `zstd` | `type.googleapis.com/envoy.extensions.compression.zstd.compressor.v3.Zstd` |

Para oferecer **mais de um algoritmo** e deixar o cliente escolher pelo `Accept-Encoding`, instale **um filtro Compressor por algoritmo** na cadeia. É a própria documentação do filtro que descreve esse arranjo: um filtro deixa de comprimir quando o `accept-encoding` dá peso maior a outra codificação cujo filtro está presente na cadeia (estatística `header_compressor_overshadowed`).

Adicione o Brotli ao gateway (o gzip do Laboratório 1 continua aplicado):

```bash
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: gateway-brotli
  namespace: $MY_NAMESPACE
spec:
  targetRefs:
    - kind: Gateway
      group: gateway.networking.k8s.io
      name: bookinfo-gateway
  configPatches:
    - applyTo: HTTP_FILTER
      match:
        context: GATEWAY
        listener:
          filterChain:
            filter:
              name: "envoy.filters.network.http_connection_manager"
              subFilter:
                name: "envoy.filters.http.router"
      patch:
        operation: INSERT_BEFORE
        value:
          name: envoy.filters.http.compressor.brotli
          typed_config:
            "@type": type.googleapis.com/envoy.extensions.filters.http.compressor.v3.Compressor
            response_direction_config:
              common_config:
                min_content_length: 100
                content_type:
                  - text/html
                  - application/json
            compressor_library:
              name: brotli
              typed_config:
                "@type": type.googleapis.com/envoy.extensions.compression.brotli.compressor.v3.Brotli
EOF
```

> O campo `name` do filtro HTTP é um identificador da configuração; aqui ele recebe o sufixo `.brotli` apenas para distinguir os dois filtros na saída dos comandos. O que define o filtro é o `@type` do `typed_config`. No exemplo oficial do Envoy com dois compressores na mesma cadeia ([compressor-filter-request-response.yaml](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/compressor_filter#using-different-compressors-for-requests-and-responses)), os dois usam o mesmo nome `envoy.filters.http.compressor`.

Teste a negociação com os pesos (`q`) do `Accept-Encoding`:

```bash
# Prefere brotli
curl -s -o /dev/null -D - -H "Accept-Encoding: br;q=1.0, gzip;q=0.8" http://${GATEWAY_IP}/productpage | grep -i content-encoding

# Prefere gzip
curl -s -o /dev/null -D - -H "Accept-Encoding: br;q=0.8, gzip;q=1.0" http://${GATEWAY_IP}/productpage | grep -i content-encoding

# Confirme os dois filtros na cadeia do gateway
istioctl proxy-config listener $GW_POD.$MY_NAMESPACE -o json | grep '"name": "envoy.filters.http.compressor'
```

> Se o seu gateway é o **Envoy Gateway** (e não o do Istio), nada disso é necessário: a compressão é um campo de alto nível (`compressor`) do recurso `BackendTrafficPolicy`. Veja o [Day 6](day6.md#compressão-de-respostas).

# Laboratório 3 (intermediário): compressão no sidecar de um workload

> **Exige o modo sidecar** ([Day 1 - sidecar](day1-sidecar.md)).

O mesmo filtro pode ser aplicado só a um workload, na cadeia **inbound** do sidecar. O exemplo completo (`enable-gzip-compression`, com `workloadSelector: app=productpage` e `context: SIDECAR_INBOUND`) está na seção [7.6 do Guia de Estudos](ISTIO_ENVOY_GUIA_DE_ESTUDOS.md#76-exemplo-prático-habilitando-compressão-gzip-via-envoyfilter), e o [Lab 12](ISTIO_ENVOY_GUIA_DE_ESTUDOS.md#lab-12---compressão-gzip-via-envoyfilter-não-aplicada) do guia traz um exercício de troubleshooting desse cenário.

Onde comprimir?

| Onde | Efeito |
|---|---|
| **Gateway de entrada** (Laboratório 1) | Comprime o tráfego para os clientes externos (normalmente onde a banda é mais cara). Um único ponto de configuração |
| **Sidecar inbound** de um workload | Comprime as respostas daquele workload para qualquer chamador que envie `Accept-Encoding`, inclusive chamadas internas da mesh. O custo de CPU é pago em cada Pod |

Teste a partir de outro Pod da mesh:

```bash
kubectl -n $MY_NAMESPACE exec deploy/ratings-v1 -c ratings -- \
  curl -s -o /dev/null -D - -H "Accept-Encoding: gzip" http://productpage:9080/productpage | grep -i -E 'content-encoding|vary'
```

# Laboratório 4 (intermediário): filtro Lua adicionando um header de resposta

> **Exige o modo sidecar.** Estrutura baseada no exemplo `reviews-lua` da referência oficial do `EnvoyFilter`.

```bash
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: productpage-lua
  namespace: $MY_NAMESPACE
spec:
  workloadSelector:
    labels:
      app: productpage
  configPatches:
    - applyTo: HTTP_FILTER
      match:
        context: SIDECAR_INBOUND
        listener:
          portNumber: 9080
          filterChain:
            filter:
              name: "envoy.filters.network.http_connection_manager"
              subFilter:
                name: "envoy.filters.http.router"
      patch:
        operation: INSERT_BEFORE
        value:
          name: envoy.filters.http.lua
          typed_config:
            "@type": "type.googleapis.com/envoy.extensions.filters.http.lua.v3.Lua"
            defaultSourceCode:
              inlineString: |
                function envoy_on_response(response_handle)
                  response_handle:headers():add("x-learning-istio", "lua")
                end
EOF

kubectl -n $MY_NAMESPACE exec deploy/ratings-v1 -c ratings -- \
  curl -s -o /dev/null -D - http://productpage:9080/productpage | grep -i x-learning-istio
```

O filtro Lua expõe as funções `envoy_on_request(request_handle)` e `envoy_on_response(response_handle)`. O exemplo da referência oficial vai além: faz uma chamada HTTP a um serviço externo de dentro do script (`request_handle:httpCall(...)`) e, para isso, usa um segundo patch (`applyTo: CLUSTER`, `operation: ADD`) que **adiciona o cluster** chamado pelo script.

# Laboratório 5 (avançado): rate limit local

> **Exige o modo sidecar.** Fonte: [Enabling Rate Limits using Envoy](https://istio.io/v1.32/docs/tasks/policy-enforcement/rate-limit/) - a tarefa oficial também cobre o rate limit **global** (com um serviço de rate limit externo).

O rate limit local é aplicado por instância, no próprio proxy, sem chamar nenhum outro serviço. O `EnvoyFilter` abaixo (da tarefa oficial) limita o `productpage` a 4 requisições por minuto e adiciona o header `x-local-rate-limit` nas respostas bloqueadas:

```bash
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: filter-local-ratelimit-svc
  namespace: istio-system
spec:
  workloadSelector:
    labels:
      app: productpage
  configPatches:
    - applyTo: HTTP_FILTER
      match:
        context: SIDECAR_INBOUND
        listener:
          filterChain:
            filter:
              name: "envoy.filters.network.http_connection_manager"
      patch:
        operation: INSERT_BEFORE
        value:
          name: envoy.filters.http.local_ratelimit
          typed_config:
            "@type": type.googleapis.com/udpa.type.v1.TypedStruct
            type_url: type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit
            value:
              stat_prefix: http_local_rate_limiter
              token_bucket:
                max_tokens: 4
                tokens_per_fill: 4
                fill_interval: 60s
              filter_enabled:
                runtime_key: local_rate_limit_enabled
                default_value:
                  numerator: 100
                  denominator: HUNDRED
              filter_enforced:
                runtime_key: local_rate_limit_enforced
                default_value:
                  numerator: 100
                  denominator: HUNDRED
              response_headers_to_add:
                - append: false
                  header:
                    key: x-local-rate-limit
                    value: 'true'
EOF
```

> Repare que a tarefa oficial cria o recurso no namespace raiz (`istio-system`) com `workloadSelector: app=productpage`: ele vale para os workloads com esse label em **qualquer** namespace.

Teste: a partir da 5ª requisição dentro de um minuto, o código passa a ser `429`:

```bash
for i in $(seq 1 6); do
  kubectl -n $MY_NAMESPACE exec deploy/ratings-v1 -c ratings -- \
    curl -s -o /dev/null -w '%{http_code}\n' http://productpage:9080/productpage
done
```

As estatísticas do filtro vêm **desabilitadas** por padrão no Istio. Para habilitá-las, a tarefa oficial indica a annotation abaixo no template do Deployment:

```yaml
template:
  metadata:
    annotations:
      proxy.istio.io/config: |-
        proxyStatsMatcher:
          inclusionRegexps:
          - ".*http_local_rate_limit.*"
```

```bash
kubectl -n $MY_NAMESPACE exec deploy/productpage-v1 -c istio-proxy -- pilot-agent request GET stats | grep http_local_rate_limit
```

A tarefa oficial mostra ainda como limitar **apenas uma rota** (patch `HTTP_ROUTE` com `typed_per_filter_config`).

# Laboratório 6 (avançado): ajustando o HTTP Connection Manager com MERGE

Nem todo `EnvoyFilter` insere um filtro. Com `operation: MERGE` é possível sobrescrever campos de um objeto que o istiod já gera. O exemplo abaixo, **da referência oficial**, ajusta o *idle timeout* HTTP e o número de saltos confiáveis do `X-Forwarded-For` em um ingress gateway instalado em `istio-system`, apenas para o host SNI `app.example.com`:

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: hcm-tweaks
  namespace: istio-system
spec:
  workloadSelector:
    labels:
      istio: ingressgateway
  configPatches:
  - applyTo: NETWORK_FILTER # http connection manager is a filter in Envoy
    match:
      context: GATEWAY
      listener:
        filterChain:
          sni: app.example.com
          filter:
            name: "envoy.filters.network.http_connection_manager"
    patch:
      operation: MERGE
      value:
        typed_config:
          "@type": "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager"
          xff_num_trusted_hops: 5
          common_http_protocol_options:
            idle_timeout: 30s
```

Dois detalhes da referência sobre `MERGE` x `REPLACE`:

- O `MERGE` usa a semântica de *proto merge*, que **não consegue remover campos**.
- Quando é preciso trocar a configuração inteira de um filtro (sem herdar campos), a referência recomenda `REPLACE`.

> O `workloadSelector` precisa casar com os labels reais do seu gateway. Confira com `kubectl -n istio-system get pods --show-labels`. No chart `istio/gateway`, o label `istio` é o nome do release sem o prefixo `istio-` (um release chamado `istio-ingress`, como o deste repositório, recebe `istio: ingress`).

# Modo ambient: estendendo o waypoint com TrafficExtension

Como o `EnvoyFilter` não é suportado em waypoints, o caminho documentado para estender o L7 no modo ambient é a API [`TrafficExtension`](https://istio.io/v1.32/docs/reference/config/proxy_extensions/traffic_extension/), que provê um mecanismo unificado para extensões **WebAssembly e Lua** e, segundo a referência, substitui (*supersedes*) o `WasmPlugin`. No modo ambient, o recurso precisa ser anexado ao waypoint (ou ao gateway) com `targetRefs`.

> ATENÇÃO!!! Na documentação do Istio 1.32 essas tarefas estão marcadas com o status **Alpha**.

Exemplo da tarefa oficial [Extend waypoints with Lua](https://istio.io/v1.32/docs/ambient/usage/extend-waypoint-lua/), aplicado ao gateway do Bookinfo: lê o header de requisição `x-number` e devolve o header de resposta `x-parity` (`odd`/`even`):

```bash
kubectl -n $MY_NAMESPACE apply -f - <<EOF
apiVersion: extensions.istio.io/v1alpha1
kind: TrafficExtension
metadata:
  name: parity-at-gateway
spec:
  targetRefs:
    - kind: Gateway
      group: gateway.networking.k8s.io
      name: bookinfo-gateway
  phase: STATS
  lua:
    inlineCode: |
      function envoy_on_request(request_handle)
        local number = tonumber(request_handle:headers():get("x-number"))
        if number == nil then return end
        local parity = number % 2 == 0 and "even" or "odd"
        request_handle:streamInfo():dynamicMetadata():set(
          "envoy.filters.http.lua", "parity", parity)
      end
      function envoy_on_response(response_handle)
        local meta = response_handle:streamInfo():dynamicMetadata():get(
          "envoy.filters.http.lua")
        if meta == nil then return end
        response_handle:headers():add("x-parity", meta["parity"])
      end
EOF

curl -s -o /dev/null -D - -H "x-number: 4" http://${GATEWAY_IP}/productpage | grep -i x-parity
```

Para aplicar a um **waypoint**, crie-o (`istioctl waypoint apply -n $MY_NAMESPACE --enroll-namespace --wait`) e troque o `name` do `targetRefs` pelo nome do waypoint (por padrão, `waypoint`). A versão com WebAssembly está em [Extend waypoints with WebAssembly plugins](https://istio.io/v1.32/docs/ambient/usage/extend-waypoint-wasm/).

# Cuidados em upgrades (Istio 1.32)

Das [notas de upgrade do Istio 1.32](https://istio.io/latest/news/releases/1.32.x/announcing-1.32/upgrade-notes/):

- **Patches que casam o filtro RBAC do gateway pelo nome**: com `ENABLE_GATEWAY_API_HTTPROUTE_AUTH=true`, o filtro HTTP da ação `ALLOW` nos ingress gateways passa a se chamar `istio.authorization.allow` em vez de `envoy.filters.http.rbac` (os filtros de `DENY` e `AUDIT` continuam como `envoy.filters.http.rbac`). Um `EnvoyFilter` com `subFilter.name: envoy.filters.http.rbac` deixa de casar o filtro `ALLOW`: `MERGE`, `REPLACE` e `REMOVE` param de ser aplicados **silenciosamente**, e `INSERT_BEFORE`/`INSERT_AFTER` **descartam** o filtro fornecido em vez de inseri-lo em outra posição.
- **Nenhum aviso é emitido** quando um patch deixa de casar, a menos que `PILOT_ENVOY_FILTER_STATS=true` esteja definido no `istiod`, o que expõe a métrica `pilot_envoy_filter_status`.

E do [anúncio da release](https://istio.io/latest/news/releases/1.32.x/announcing-1.32/):

- Novo campo `transportSocket` no match de filter chain (visto acima).
- Novo aviso `IST0178` do `istioctl analyze`, emitido quando um `EnvoyFilter` insere um filtro HTTP junto de uma `AuthorizationPolicy` que mira um `HTTPRoute`.

Regra geral: a cada upgrade do Istio, reaplique o roteiro de [Troubleshooting](#troubleshooting) em **todos** os `EnvoyFilter` existentes. Uma ferramenta que ajuda a prever mudanças entre versões do Istio e do Envoy é a CLI deste repositório: [agents/upgrade](agents/upgrade/README.md) (`istio-upgrade analyze`).

# Troubleshooting

```bash
# 1. O recurso existe e está no namespace certo?
kubectl get envoyfilter -A

# 2. O istioctl aponta algum problema?
istioctl analyze -n $MY_NAMESPACE

# 3. O filtro chegou na configuração REAL do proxy? (a única prova de que o patch casou)
istioctl proxy-config listener $GW_POD.$MY_NAMESPACE -o json | grep -B2 -A12 'envoy.filters.http.compressor'

# 4. O proxy aceitou a configuração? (SYNCED) Uma configuração inválida é rejeitada pelo Envoy
istioctl proxy-status
kubectl -n istio-system logs deploy/istiod | grep -i -E 'envoyfilter|rejected|NACK' | tail -20

# 5. Qual é a ordem real dos filtros HTTP?
istioctl proxy-config listener $GW_POD.$MY_NAMESPACE -o json | grep '"name": "envoy.filters.http'
```

| Sintoma | Causa provável |
|---|---|
| O patch "não faz nada" | `workloadSelector`/`targetRefs` não casa com o proxy; `context` errado (ex.: `SIDECAR_INBOUND` em um gateway); `subFilter` com um nome que não existe naquela cadeia |
| Resposta sem `content-encoding` | Uma das condições de [Quando o Envoy comprime](#quando-o-envoy-comprime-e-quando-não-comprime) - geralmente `Accept-Encoding` ausente, `content-type` fora da lista ou resposta menor que o mínimo |
| Proxy fica `STALE`/configuração rejeitada | Campo inexistente na versão do Envoy empacotada, ou `@type` incorreto |
| Funcionava e parou após um upgrade | Nome de filtro/estrutura interna mudou - veja [Cuidados em upgrades](#cuidados-em-upgrades-istio-132) |
| `EnvoyFilter` aplicado a um waypoint não tem efeito | Não é suportado - use [TrafficExtension](#modo-ambient-estendendo-o-waypoint-com-trafficextension) |

Mais comandos na seção [7.7 do Guia de Estudos](ISTIO_ENVOY_GUIA_DE_ESTUDOS.md#77-comandos-de-troubleshooting-de-envoyfilter).

# Limpeza

```bash
kubectl -n $MY_NAMESPACE delete envoyfilter gateway-gzip gateway-brotli productpage-lua --ignore-not-found
kubectl -n istio-system delete envoyfilter filter-local-ratelimit-svc --ignore-not-found
kubectl -n $MY_NAMESPACE delete trafficextension parity-at-gateway --ignore-not-found
```

# Referências

- https://istio.io/v1.32/docs/reference/config/networking/envoy-filter/
- https://istio.io/v1.32/docs/tasks/policy-enforcement/rate-limit/
- https://istio.io/v1.32/docs/ops/configuration/telemetry/envoy-stats/
- https://istio.io/v1.32/docs/ops/diagnostic-tools/proxy-cmd/
- https://istio.io/v1.32/docs/ambient/usage/l7-features/
- https://istio.io/v1.32/docs/ambient/migrate/
- https://istio.io/v1.32/docs/ambient/usage/extend-waypoint-lua/
- https://istio.io/v1.32/docs/ambient/usage/extend-waypoint-wasm/
- https://istio.io/v1.32/docs/reference/config/proxy_extensions/traffic_extension/
- https://istio.io/v1.32/docs/reference/config/proxy_extensions/wasm-plugin/
- https://istio.io/latest/news/releases/1.32.x/announcing-1.32/
- https://istio.io/latest/news/releases/1.32.x/announcing-1.32/upgrade-notes/
- https://github.com/istio/proxy/blob/master/bazel/extension_config/extensions_build_config.bzl
- https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/http_filters
- https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/compressor_filter
- https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/lua_filter
- https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/local_rate_limit_filter
- https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/http/http_filters
