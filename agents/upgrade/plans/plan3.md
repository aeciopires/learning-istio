# Plano executado 3 — day1-sidecar.md, suporte a macOS, addons do Istio e Bookinfo no Makefile

> Continuação de [plan1.md](plan1.md) e [plan2.md](plan2.md).

## Pedido

1. Restaurar/atualizar um `day1-sidecar.md` na raiz do repo (hoje só existe uma versão
   desatualizada em `OLD/day1-sidecar.md`) e fazer o CLI (`docs update`), o Makefile, os
   testes, o `CLAUDE.md`, o `CHANGELOG.md` e os diagramas mermaid passarem a considerá-lo.
2. Suporte a build e execução do CLI (todos os comandos) em **macOS amd64 e arm64**.
3. Opção no Makefile para instalar todos os addons do Istio citados no fim do
   `day1.md`/`day1-sidecar.md` (Grafana, Prometheus, Kiali, Jaeger), com o CLI dando suporte
   a rastrear/atualizar a versão desses addons também.
4. Opção no Makefile para instalar o Bookinfo citado em `day1.md`/`day1-sidecar.md`.
5. As novas opções (addons e Bookinfo) devem entrar em `bootstrap-sidecar`/`bootstrap-ambient`,
   permitindo instalar tudo de uma vez ou como passos separados - igual já acontece hoje para
   os componentes do Istio (`create-cluster`, `install-metallb`, `install-istio-{sidecar,ambient}`
   são passos individuais; `bootstrap-*` roda tudo em sequência).

## Decisões de design

### 1. `day1-sidecar.md`
Novo arquivo na raiz do repo, espelhando a estrutura e as versões atuais de `day1.md`
(`ISTIO_RELEASE=1.30`, `VERSION_ISTIO=${ISTIO_RELEASE}.3`, `GATEWAY_API_VERSION=v1.6.1`), mas
com os passos de instalação do modo **sidecar** (sem `--set profile=ambient`, sem `ztunnel`,
label `istio-injection=enabled` em vez de `istio.io/dataplane-mode=ambient`, e o passo extra de
`kubectl rollout restart deployment` para os pods do Bookinfo ganharem o sidecar). Baseado em
`OLD/day1-sidecar.md`, mas sem a seção avançada de customização de logs (fora do escopo do
walkthrough padrão; quem quiser isso ainda a encontra em `OLD/`). `agents/upgrade`'s
`docsupdater.UpdateDayFile` já é genérico o bastante (regex sobre `ISTIO_RELEASE`/`VERSION_ISTIO`/
`GATEWAY_API_VERSION`/links `istio.io/vX.Y/`) para ser reaproveitado sem mudanças - só precisa
entrar na lista de arquivos que `Updater.Plan` processa.

### 2. macOS (amd64/arm64)
- O CLI em Go já compila nativamente para `darwin/amd64`/`darwin/arm64` (sem cgo, sem chamadas
  de syscall específicas de SO). Adiciona-se ao `Makefile` os alvos `build-linux-amd64`,
  `build-linux-arm64`, `build-darwin-amd64`, `build-darwin-arm64` e `build-all` (via
  `GOOS`/`GOARCH`), publicando em `bin/istio-upgrade-<os>-<arch>`.
- A imagem Docker continua **só Linux** (Docker Desktop no Mac roda containers Linux via VM;
  não existe "imagem Docker para macOS") - documentado, não é uma lacuna.
- Os `scripts/*.sh` (automação de cluster) têm dois pontos não portáveis para macOS/BSD,
  corrigidos:
  - `mktemp -t kind-config-XXXXXX.yaml` (sintaxe GNU) trocado por
    `mktemp "${TMPDIR:-/tmp}/kind-config.XXXXXX"` (portável entre GNU e BSD mktemp).
  - A checagem de `fs.inotify` em `preflight.sh` é específica do kernel Linux (não existe no
    macOS - o kernel relevante ali é o da VM Linux do Docker Desktop, não o host). Adiciona-se
    detecção de SO (`uname -s`) para pular essa checagem no Darwin com uma nota informativa,
    em vez de falhar tentando rodar `sysctl -n fs.inotify...` (que não existe no macOS).
  - O restante dos scripts (`sed` sem `-i`, `docker network inspect`, `jq`, `kubectl`, `helm`)
    já é portável (POSIX/GNU e BSD se comportam igual nesses comandos específicos).
- `INSTALL_REQUIREMENTS.md` continua com foco em Ubuntu/Debian (não é o pedido - o pedido é
  sobre o CLI/build/scripts funcionarem em macOS, não reescrever o guia de pacotes de SO); a
  nota já existente sobre macOS + `docker-mac-net-connect` em `day1.md`/`day1-sidecar.md` é
  mantida como está.

### 3. Addons do Istio (Grafana, Prometheus, Kiali, Jaeger)
- Novo `scripts/install-addons.sh`, espelhando o trecho final de `day1.md`/`day1-sidecar.md`
  (`kubectl apply` dos 4 manifests em `$ISTIO_ADDONS_URL` + `kubectl rollout status` de cada
  deployment). Novo alvo `make install-addons` no `Makefile` (não depende de modo
  sidecar/ambient - os manifests são os mesmos).
- **Rastreio de versão pelo CLI**: os manifests de addons são baixados prontos do repositório
  `istio/istio` (pinados por `ISTIO_RELEASE`, já mantido automaticamente pelo `docs update`) -
  não há um número de versão independente para "escolher" nesse mecanismo de instalação.
  Ainda assim, Grafana/Prometheus/Kiali/Jaeger entram no catálogo de
  `internal/toolversions` (repositórios `grafana/grafana`, `prometheus/prometheus`,
  `kiali/kiali`, `jaegertracing/jaeger`), aparecendo em `istio-upgrade versions latest` e em
  `configs/versions.env` (`GRAFANA_VERSION`, `PROMETHEUS_VERSION`, `KIALI_VERSION`,
  `JAEGER_VERSION`) como referência informativa - útil para saber se a versão empacotada pelo
  Istio está defasada frente ao upstream. Isso é documentado explicitamente (não é um pin que
  controla o que é de fato instalado, já que a instalação usa os manifests do Istio).

### 4. Bookinfo no Makefile
Já existe (`deploy-bookinfo-sidecar`/`deploy-bookinfo-ambient` → `scripts/deploy-bookinfo.sh`),
implementado numa sessão anterior. Mantido como está; só passa a ser citado/linkado também a
partir de `day1-sidecar.md`.

### 5. Integração em `bootstrap-sidecar`/`bootstrap-ambient`
`scripts/bootstrap.sh` passa a incluir addons como parte do "instalar tudo", controlável por
variáveis de ambiente (paralelo ao que já existe para escolher `sidecar`/`ambient` via
argumento posicional):

```bash
make bootstrap-sidecar                    # tudo, incluindo addons (padrão)
INSTALL_ADDONS=false make bootstrap-sidecar   # pula os addons
```

As "opções separadas" continuam sendo os alvos individuais já existentes
(`create-cluster`, `install-metallb`, `install-istio-{sidecar,ambient}`,
`deploy-bookinfo-{sidecar,ambient}`) mais o novo `install-addons` - cada um pode ser rodado
isoladamente fora do `bootstrap-*`.

## Arquivos afetados

- **Novo**: `day1-sidecar.md`, `agents/upgrade/scripts/install-addons.sh`,
  `agents/upgrade/plans/plan3.md` (este arquivo).
- **Go (`agents/upgrade/app/`)**: `internal/docsupdater/updater.go` (adiciona
  `day1-sidecar.md` à lista de day-files), `internal/toolversions/toolversions.go`
  (catálogo com Grafana/Prometheus/Kiali/Jaeger), `internal/docsupdater/files.go`
  (`RenderVersionsEnv` com as 4 novas variáveis), testes correspondentes.
- **Scripts/Makefile**: `scripts/create-cluster.sh` (mktemp portável),
  `scripts/preflight.sh` (detecção de SO para o check de inotify), `scripts/bootstrap.sh`
  (`INSTALL_ADDONS`), `Makefile` (alvos `build-<os>-<arch>`, `build-all`, `install-addons`).
- **Documentação**: `agents/upgrade/README.md` (diagramas mermaid, seção de build
  multiplataforma, addons), `CLAUDE.md`, `agents/upgrade/app/CHANGELOG.md`.

## Verificação

- `make test`/`make lint` (unitários, offline).
- `make build-all` gerando os 4 binários (linux/darwin × amd64/arm64) sem erro de compilação
  cruzada.
- `make docs-update` (dry-run) e `--apply` contra o repo real, incluindo o novo
  `day1-sidecar.md`.
- Execução real de `make install-addons` (e, se o tempo permitir, um `make bootstrap-sidecar`
  completo com addons) contra um cluster kind de verdade, como já foi feito para os demais
  componentes nesta sessão.
