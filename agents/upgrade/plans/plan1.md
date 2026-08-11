# Plano de implantação executado — CLI `istio-upgrade` (`agents/upgrade`)

> Este documento registra, de forma retrospectiva, o plano que foi de fato
> seguido para criar o conteúdo de `agents/upgrade/` nesta sessão -
> incluindo as correções de escopo pedidas pelo usuário no meio do
> trabalho e os bugs reais encontrados/corrigidos durante a validação com
> infraestrutura real (kind + Docker). Serve como referência para quem
> for entender "por que as coisas foram feitas assim" ou planejar a
> próxima extensão do CLI (`plan2.md`, etc.).

## Contexto e pedido original

O repositório é material de estudo (pt-BR) de Istio/Envoy, com
walkthroughs manuais (`day1.md`/`day2.md`) fixados em versões específicas
de ferramentas (`README.md`, `INSTALL_REQUIREMENTS.md`, `.tool-versions`).
Foi pedido um CLI em Go, na pasta `agents/upgrade`, capaz de:

1. Ler as release notes do Istio e do Envoy e listar mudanças de código,
   configuração, *breaking changes* e instalação entre uma versão inicial
   e final, considerando todas as versões intermediárias (ex.: Istio 1.10
   → 1.30).
2. Ter opções como `--help`, `--version`, uma opção para analisar release
   notes, e uma opção para atualizar o conteúdo de instalação deste
   repositório.
3. Automatizar a criação do cluster e a instalação do Istio, MetalLB e
   Bookinfo, nos modos ambient e sidecar.
4. Ter um arquivo `VERSION` para controle centralizado de versão, um
   `Makefile` para compilar/executar/rodar em Docker, um `README.md`
   explicando o funcionamento, e testes automatizados para o CLI e para
   validar o funcionamento da instalação do cluster e do Istio.
5. Usar as versões de software mais novas dos sites oficiais.
6. Criar um `CLAUDE.md` na raiz do repositório com informações úteis para IA.

## Correções de escopo feitas pelo usuário durante a sessão

O plano inicial (definido em modo de planejamento, antes de qualquer
código) foi ajustado três vezes por instrução explícita do usuário,
**depois** da aprovação do plano original:

1. **Comparação de versões minor e patch/hotfix**: o `analyze` precisa
   suportar tanto `1.10` → `1.30` (minor a minor) quanto `1.29.0` →
   `1.29.3` (hotfix dentro do mesmo minor), não só ranges de minor.
2. **O CLI não cria o cluster**: a automação de cluster (kind, MetalLB,
   Istio, Bookinfo) é responsabilidade do **Makefile + scripts shell**,
   não do binário Go. O CLI só analisa release notes e atualiza arquivos
   de instalação/configuração (incluindo os que os scripts consomem).
3. **`CHANGELOG.md`** para o CLI, no formato Keep a Changelog.

Depois da primeira execução real do ambiente (ver seção de bugs abaixo),
o usuário também pediu explicitamente que o `istio-cni` fosse corrigido e
habilitado **também no modo sidecar** (inicialmente só era instalado no
modo ambient).

## Decisões de arquitetura

- **Separação CLI × automação de cluster**: `cmd/` + `internal/` (Go) só
  fazem análise de release notes e sincronização de versões; `scripts/*.sh`
  + `Makefile` fazem a orquestração real de infraestrutura (Docker/kind),
  algo que não faz sentido "containerizar" dentro da imagem do CLI. A
  ponte entre os dois mundos é um único arquivo gerado,
  `configs/versions.env`.
- **`internal/semverx`**: parser de versão próprio (sem dependência
  externa) que expande um intervalo `from`→`to` em uma lista de séries
  `Major.Minor`, com uma flag por série indicando se há recorte por patch
  (`IncludePatchDetail`, `FromPatch`, `ToPatch`). Isso permite tratar
  `1.10`→`1.30` (21 séries, sem detalhe de patch) e `1.29.0`→`1.29.3`
  (1 série, com detalhe de patch 0→3) com a mesma estrutura de dados.
- **`internal/releasenotes`**: um `Fetcher` HTTP com cache em disco
  (`internal/releasenotes/cache.go`) compartilhado por dois clientes:
  - `istio.go`: busca `announcing-X.Y/_index.md` e `.../upgrade-notes/index.md`
    (markdown cru do repositório `istio/istio.io` no GitHub) para
    features/breaking changes por minor, e usa a API de *compare* do
    GitHub (`istio/istio`) para notas de patch/hotfix.
  - `envoy.go`: busca as páginas RST de `version_history` publicadas em
    `envoyproxy.io` (fonte por patch, cumulativa por minor).
  - A versão exata do Envoy embarcada em cada release do Istio **não é
    adivinhada** (é um commit SHA de fork em `istio.deps`, não uma versão
    limpa) - o usuário precisa informar `--envoy-from`/`--envoy-to`
    explicitamente.
- **`internal/toolversions`**: última versão estável de cada ferramenta
  via API de Releases do GitHub, com fallback para a API de tags quando o
  repositório não publica GitHub Releases (ex.: MetalLB, cujo "latest
  release" às vezes aponta para uma tag `metallb-chart-X.Y.Z`, não
  `vX.Y.Z` - tratado por `normalizeVersionTag`, que extrai o sufixo
  semver da tag independente do prefixo).
- **`internal/docsupdater`**: regras de reescrita por regex, cada uma
  ancorada em texto único o suficiente para não confundir seções
  parecidas (ex.: `VERSION="..."` aparece tanto na seção do Helm quanto
  na do kind em `INSTALL_REQUIREMENTS.md` - a âncora é o texto seguinte).
  Inclui uma proteção deliberada: **bump de versão MAJOR nunca é aplicado
  automaticamente** (ex.: Helm 3 → 4), só com `--allow-major-bump`,
  porque isso costuma quebrar os comandos exatos documentados no
  repositório.
- **Testes**: unitários (`go test`, sem rede, usando `httptest` e
  fixtures) para toda a lógica em `internal/` e para validação de flags
  em `cmd/`; a validação de que "o cluster e o Istio realmente funcionam"
  é feita por `scripts/validate-cluster.sh` (shell), não por `go test`,
  porque depende de infraestrutura real (Docker/kind).

## Ordem de execução do plano

1. Scaffold do módulo Go (`go.mod`, cobra, yaml.v3) e esqueleto dos
   comandos (`root`, `version`, `analyze`, `docs update`, `versions latest`).
2. `internal/semverx` com testes (expansão de range minor e patch).
3. `internal/releasenotes` (cache HTTP, parser de markdown do Istio,
   parser de RST do Envoy, API de tags/compare do GitHub, classificação
   por categoria, renderização em markdown/JSON) com testes via
   `httptest` e fixtures.
4. `internal/toolversions` (últimas versões via GitHub) com testes.
5. `internal/docsupdater` (regras de reescrita + proteção contra bump
   major) com testes.
6. Execução **real** do `analyze` contra as APIs verdadeiras: Istio
   `1.10`→`1.30`, Istio `1.29.0`→`1.29.3` (patch), Envoy `1.31`→`1.32` -
   saídas salvas em `examples/`.
7. Execução **real** do `docs update --apply` contra o repositório: viu
   drift real (kubectl 1.36.2→1.36.3, Gateway API v1.6.0→v1.6.1), manteve
   Istio em 1.30.3 (não 1.31, que está em alpha) e bloqueou o bump de
   Helm 3→4 com um aviso.
8. `scripts/*.sh` + `Makefile` para o ciclo de vida do cluster
   (`preflight`, `create-cluster`, `install-metallb`, `install-istio`,
   `deploy-bookinfo`, `validate-cluster`, `destroy-cluster`, `bootstrap`).
9. Execução **real** de `make bootstrap-sidecar` e `make bootstrap-ambient`
   (kind + Docker de verdade), com validação de tráfego passando pela
   mesh, seguida de `make destroy-cluster` nos dois casos.
10. `Dockerfile` + build/run real da imagem, incluindo teste de HTTPS
    (`ca-certificates`) dentro do container.
11. `README.md` (com os diagramas mermaid, adicionados depois a pedido do
    usuário), `CHANGELOG.md`, `VERSION`.
12. `CLAUDE.md` na raiz do repositório.
13. Testes de `cmd/` (cobertura de `cmd/` estava em 0%, foi adicionada
    cobertura para validação de flags e detecção de repo-root).
14. Correção do `istio-cni` também no modo sidecar (`pilot.cni.enabled=true`)
    a pedido do usuário, com nova execução real de `make bootstrap-sidecar`
    para confirmar (inspecionando os init containers do pod: `istio-validation`
    em vez do antigo `istio-init` privilegiado).

## Bugs reais encontrados e corrigidos durante a execução

Nenhum destes apareceria só lendo a documentação - só foram descobertos
rodando a infraestrutura de verdade:

1. **`--set profile=sidecar` não existe** nos charts Helm do Istio. O
   chart só tem `profile=ambient`; o modo sidecar é o profile padrão do
   chart (nenhuma flag `--set profile=...`). Corrigido em
   `scripts/install-istio.sh`.
2. **`istio-cni` em `CrashLoopBackOff`** com
   `couldn't initialize inotify: too many open files` - os limites de
   `fs.inotify.max_user_watches`/`max_user_instances` do host estavam
   baixos demais (esse problema já era documentado manualmente em
   `day1.md`, mas não havia automação para ele). Corrigido em
   `scripts/preflight.sh`: detecta o limite baixo e tenta subir via
   `sudo -n sysctl -w ...`; se não houver sudo sem senha, imprime os
   comandos exatos a rodar manualmente.
3. **`istio-cni` só era instalado no modo ambient.** A pedido do usuário,
   passou a ser instalado nos **dois modos**: no modo sidecar, o
   `istiod` é instalado com `--set pilot.cni.enabled=true`, o que faz o
   injetor de sidecar parar de injetar o container de init privilegiado
   `istio-init` e delegar o redirecionamento de tráfego ao `istio-cni`
   (mesmo mecanismo usado no modo ambient, aplicado ao caso sidecar).
   Confirmado inspecionando `spec.initContainers` do pod `productpage`
   após a instalação (`istio-validation`, não `istio-init`).

## Verificação final

- `make test` / `make lint` (gofmt + go vet + golangci-lint) limpos.
- Cobertura: `semverx` ~92%, `toolversions` ~87%, `releasenotes` ~80%,
  `docsupdater` ~73%, `cmd` ~31% (partia de 0%).
- `make bootstrap-sidecar` e `make bootstrap-ambient` executados de
  verdade, duas vezes cada (uma antes e uma depois da correção do
  `istio-cni`/inotify), com `make validate` confirmando: cluster kind
  acessível, control plane do Istio pronto, `istio-cni-node` pronto (nos
  dois modos), `ztunnel` pronto (modo ambient), Bookinfo disponível e
  tráfego real chegando ao `productpage` de dentro da mesh.
- `make docker-build` + `docker run istio-upgrade-cli --version`/`--help`/
  `versions latest` (testando HTTPS dentro do container).
- `docs update --apply` aplicado de fato no repositório (diff pequeno e
  seguro, revisado antes: `.tool-versions`, `INSTALL_REQUIREMENTS.md`,
  `README.md`, `day1.md`, `day2.md`).

## O que ficou fora do escopo (deliberadamente)

- Reescrita automática da tag+digest sha256 da imagem `kindest/node` em
  `INSTALL_REQUIREMENTS.md` (precisa de verificação manual do digest).
- `ISTIO_ENVOY_GUIA_DE_ESTUDOS.md`, `UNINSTALL_ISTIO.md` e `OLD/` não são
  tocados pelo `docs update` (prosa curada / conteúdo arquivado).
- Bump automático de versão MAJOR (ex.: Helm 4) sem `--allow-major-bump`.
- Instalação de addons de observabilidade (Grafana/Prometheus/Kiali/Jaeger)
  na automação de cluster - não foi pedido; `ISTIO_ADDONS_URL` já está
  disponível em `scripts/lib/common.sh` para uma extensão futura.
